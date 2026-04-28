// Copyright 2025 The Score Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package uriget provides a mechanism for "loading" the contents of a file from a uri with flexible support for
// different uri schemes. This is primarily used for loading Score provisioner files for various Score implementations.
// There are similar packages such as hashicorp/go-getter, however this version is maintained to have zero dependencies
// and use the binaries that exist on the local system where possible.
package uriget

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

// FileContent holds the URI and content of a file retrieved by GetFiles.
type FileContent struct {
	// URI is the URI or path of the file.
	URI string
	// Content is the raw bytes of the file.
	Content []byte
}

// options is a struct holding fields that may need to have overrides in certain environments or during unit testing.
// The options struct can be modified by using Option functions. See defaultOptions.
type options struct {
	// limit is the limit of bytes to read from the target file. This is a safety mechanism to ensure we aren't loading
	// from a file that is too large. See WithLimit.
	limit int
	// logger is used to log messages from within the getter functions. These messages are informational only and can
	// be muted as necessary. See WithLogger.
	logger *log.Logger

	// httpClient is the http client implementation to use. See WithHttpClient.
	httpClient HttpDoer

	// tempDir is a temporary directory which may be used for storing buffers or temporary files.
	tempDir string
}

// HttpDoer is an http.Client interface used for overrides during testing or other http fetching implementations.
type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Option is an option function that modifies the options structure in place.
type Option func(*options)

// WithLimit sets the io reader limit in bytes.
func WithLimit(b int) Option {
	return func(o *options) {
		o.limit = b
	}
}

// WithHttpClient sets the http client that may be used.
func WithHttpClient(c HttpDoer) Option {
	return func(o *options) {
		o.httpClient = c
	}
}

// WithTempDir sets the temporary data directory.
func WithTempDir(p string) Option {
	return func(o *options) {
		o.tempDir = p
	}
}

// WithLogger sets the logging output
func WithLogger(l *log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

var defaultOptions = []Option{
	WithLimit(1024 * 1024 * 1024),
	WithLogger(log.Default()),
	WithHttpClient(&http.Client{
		Timeout: time.Second * 30,
	}),
	WithTempDir(os.TempDir()),
}

const ()

// GetFile attempts to parse and retrieve file contents from the given url/uri. The scheme is used to inform what sources
// are supported and how the remainder of the url should be handled.
// Supported schemes:
// - http/https: reads the file using http.
// - file or no scheme: attempts to read the file from local file system.
// - git-ssh / git-https: attempts to perform a sparse checkout of just the target file.
// - oci: retrieves a file from a remote OCI registry based on the reference and optional fragment.
//
// Deprecated: Use GetFiles instead, which supports both single files and directories.
func GetFile(ctx context.Context, rawUri string, optionFuncs ...Option) ([]byte, error) {
	u, err := url.Parse(rawUri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	opts := &options{}
	for _, optionFunc := range append(defaultOptions, optionFuncs...) {
		optionFunc(opts)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		return opts.getHttp(ctx, u)
	case "file", "":
		return opts.getFile(ctx, u)
	case "git-ssh", "git-https":
		return opts.getGit(ctx, u)
	case "oci":
		return opts.getOci(ctx, u)
	default:
		return nil, fmt.Errorf("unsupported scheme '%s'", u.Scheme)
	}
}

// GetFiles is like GetFile but with support for importing multiple files from a directory. Currently, directory
// support is implemented for the file and git schemes. For other schemes (http, oci), the target is treated as a
// single file and returned as a single-element slice.
//
// TODO: Add directory support for oci scheme.
func GetFiles(ctx context.Context, rawUri string, optionFuncs ...Option) ([]FileContent, error) {
	u, err := url.Parse(rawUri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	opts := &options{}
	for _, optionFunc := range append(defaultOptions, optionFuncs...) {
		optionFunc(opts)
	}
	var content []byte
	switch strings.ToLower(u.Scheme) {
	case "file", "":
		return opts.getFileOrDir(ctx, u)
	case "http", "https":
		content, err = opts.getHttp(ctx, u)
	case "git-ssh", "git-https":
		return opts.getGitFileOrDir(ctx, u)
	case "oci":
		content, err = opts.getOci(ctx, u)
	default:
		return nil, fmt.Errorf("unsupported scheme '%s'", u.Scheme)
	}
	if err != nil {
		return nil, err
	}
	return []FileContent{{URI: rawUri, Content: content}}, nil
}

func getStdinFile(ctx context.Context) ([]byte, error) {
	// Check if stdin is being piped
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	// Check if stdin is a pipe
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("stdin is empty")
		}
		return data, nil
	}

	return nil, fmt.Errorf("no stdin data provided")
}

func readLimited(r io.Reader, limit int) ([]byte, error) {
	if buff, err := io.ReadAll(io.LimitReader(r, int64(limit+1))); err == nil && len(buff) > limit {
		return nil, fmt.Errorf("%d byte limit exceeded", limit)
	} else {
		return buff, err
	}
}

func (o *options) getHttp(ctx context.Context, u *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("bad url: %w", err)
	}
	req = req.WithContext(ctx)
	res, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make get request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s non-200 status code: %d", req.Method, req.URL, res.StatusCode)
	}
	buff, err := readLimited(res.Body, o.limit)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	o.logger.Printf("Read %d bytes from %s %s", len(buff), req.Method, req.URL)
	return buff, nil
}

func (o *options) getFile(ctx context.Context, u *url.URL) ([]byte, error) {
	targetPath := u.Host + u.Path
	rawUri := u.String()
	if rawUri == "-" {
		return getStdinFile(ctx)
	}

	if strings.HasPrefix(targetPath, "~/") {
		hd, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to find user home dir: %w", err)
		}
		targetPath = filepath.Join(hd, targetPath[2:])
	}
	f, err := os.Open(targetPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	buff, err := readLimited(f, o.limit)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	o.logger.Printf("Read %d bytes from %s", len(buff), targetPath)
	return buff, nil
}

// getFileOrDir resolves a file:// or bare path URI. If the path is a directory, it reads all files (non-recursively)
// sorted by name. If it's a single file, it returns a single-element slice.
func (o *options) getFileOrDir(ctx context.Context, u *url.URL) ([]FileContent, error) {
	targetPath := u.Host + u.Path
	rawUri := u.String()
	if rawUri == "-" {
		content, err := getStdinFile(ctx)
		if err != nil {
			return nil, err
		}
		return []FileContent{{URI: "-", Content: content}}, nil
	}

	if strings.HasPrefix(targetPath, "~/") {
		hd, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to find user home dir: %w", err)
		}
		targetPath = filepath.Join(hd, targetPath[2:])
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		f, err := os.Open(targetPath)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close() }()
		buff, err := readLimited(f, o.limit)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		o.logger.Printf("Read %d bytes from %s", len(buff), targetPath)
		return []FileContent{{URI: targetPath, Content: buff}}, nil
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Filter to skip sub-directories
	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	sort.Strings(fileNames)

	if len(fileNames) == 0 {
		return nil, fmt.Errorf("directory %s contains no files", targetPath)
	}

	var out []FileContent
	for _, name := range fileNames {
		filePath := filepath.Join(targetPath, name)
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
		}
		buff, err := readLimited(f, o.limit)
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
		}
		o.logger.Printf("Read %d bytes from %s", len(buff), filePath)
		out = append(out, FileContent{URI: filePath, Content: buff})
	}
	return out, nil
}

func (o *options) getGit(ctx context.Context, u *url.URL) ([]byte, error) {
	u.Scheme = strings.TrimPrefix(u.Scheme, "git-")
	u.RawQuery = ""
	u.Fragment = ""
	parts := strings.SplitN(u.Path, ".git/", 2)
	if len(parts) == 1 || parts[0] == "" || strings.HasSuffix(parts[1], "/") {
		return nil, fmt.Errorf("invalid git url, expected a path with ../<REPO>.git/<FILEPATH>")
	}
	u.Path = parts[0] + ".git"
	subPath := parts[1]

	td, err := o.gitSparseCheckout(ctx, u.String(), subPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(td) }()

	f, err := os.Open(filepath.Join(td, subPath))
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	buff, err := readLimited(f, o.limit)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	o.logger.Printf("Read %d bytes from %s", len(buff), filepath.Join(td, subPath))
	return buff, nil
}

// parseGitUrl parses a git URI and returns the remote URL and the sub-path within the repo.
// It accepts paths with or without a trailing slash.
func parseGitUrl(u *url.URL) (remoteUrl string, subPath string, err error) {
	u.Scheme = strings.TrimPrefix(u.Scheme, "git-")
	u.RawQuery = ""
	u.Fragment = ""
	parts := strings.SplitN(u.Path, ".git/", 2)
	if len(parts) == 1 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid git url, expected a path with ../<REPO>.git/<PATH>")
	}
	u.Path = parts[0] + ".git"
	subPath = strings.TrimSuffix(parts[1], "/")
	return u.String(), subPath, nil
}

// gitSparseCheckout performs a sparse checkout of the given subPath from the git remote into a temp directory.
// The caller is responsible for cleaning up the returned temp directory.
func (o *options) gitSparseCheckout(ctx context.Context, remoteUrl string, subPath string) (string, error) {
	td, err := os.MkdirTemp(os.TempDir(), "score-go")
	if err != nil {
		return "", fmt.Errorf("failed to make temp dir")
	} else if err := os.Chmod(td, 0700); err != nil {
		_ = os.RemoveAll(td)
		return "", fmt.Errorf("failed to chown temp dir")
	}

	gitBinary, err := exec.LookPath("git")
	if err != nil {
		_ = os.RemoveAll(td)
		return "", fmt.Errorf("failed to find git binary on the local system: %w", err)
	}

	for _, step := range []struct {
		args   []string
		errMsg string
	}{
		{[]string{"init"}, "failed to init git repo in " + td},
		{[]string{"remote", "add", "origin", remoteUrl}, "failed to set git remote to " + remoteUrl},
		{[]string{"sparse-checkout", "set", "--no-cone", "--sparse-index", subPath}, "failed to set sparse checkout"},
		{[]string{"pull", "origin", "HEAD", "--depth=1"}, "failed to fetch"},
	} {
		c := exec.CommandContext(ctx, gitBinary, step.args...)
		c.Dir = td
		if output, err := c.CombinedOutput(); err != nil {
			o.logger.Printf("command output: %s", output)
			_ = os.RemoveAll(td)
			return "", fmt.Errorf("%s: %w", step.errMsg, err)
		}
	}
	o.logger.Printf("Initialized git remote in %s for %s", td, remoteUrl)
	return td, nil
}

// getGitFileOrDir is like getGit but returns multiple files if the subPath is a directory.
func (o *options) getGitFileOrDir(ctx context.Context, u *url.URL) ([]FileContent, error) {
	originalUri := u.String()
	remoteUrl, subPath, err := parseGitUrl(u)
	if err != nil {
		return nil, err
	}

	td, err := o.gitSparseCheckout(ctx, remoteUrl, subPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(td) }()

	fullPath := filepath.Join(td, subPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		f, err := os.Open(fullPath)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close() }()
		buff, err := readLimited(f, o.limit)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		o.logger.Printf("Read %d bytes from %s", len(buff), fullPath)
		return []FileContent{{URI: originalUri, Content: buff}}, nil
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	sort.Strings(fileNames)

	if len(fileNames) == 0 {
		return nil, fmt.Errorf("directory %s contains no files", subPath)
	}

	var out []FileContent
	for _, name := range fileNames {
		filePath := filepath.Join(fullPath, name)
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
		}
		buff, err := readLimited(f, o.limit)
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
		}
		o.logger.Printf("Read %d bytes from %s", len(buff), filePath)
		out = append(out, FileContent{URI: subPath + "/" + name, Content: buff})
	}
	return out, nil
}

func (o *options) getOci(ctx context.Context, u *url.URL) ([]byte, error) {
	ref, err := registry.ParseReference(u.Host + u.Path)
	if err != nil {
		return nil, fmt.Errorf("invalid artifact URL: %w", err)
	}
	if ref.Reference == "" {
		ref.Reference = "latest"
	}
	specifiedFile := strings.TrimPrefix(u.Fragment, "#")
	storeOpts := credentials.StoreOptions{}
	credStore, err := credentials.NewStoreFromDocker(storeOpts)
	if err != nil {
		o.logger.Printf("Warning: Unable to load Docker credentials, continuing without auth. Error: %v", err)
	}
	remoteRepo, err := remote.NewRepository(ref.String())
	if err != nil {
		return nil, fmt.Errorf("connection to remote repository failed: %w", err)
	}
	remoteRepo.PlainHTTP = strings.HasPrefix(ref.Registry, "localhost") || strings.HasPrefix(ref.Registry, "127.0.0.1")
	remoteRepo.Client = &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(credStore),
	}
	_, rc, err := remoteRepo.Manifests().FetchReference(ctx, ref.Reference)
	if err != nil {
		return nil, fmt.Errorf("manifest fetch failed: %w", err)
	}
	defer rc.Close()
	var manifest v1.Manifest
	if err := json.NewDecoder(rc).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("manifest decode failed: %w", err)
	}
	var selectedLayer *v1.Descriptor
	yamlFileCount := 0
	for _, layer := range manifest.Layers {
		title := layer.Annotations[v1.AnnotationTitle]
		if strings.HasSuffix(title, ".yaml") {
			yamlFileCount++
			if specifiedFile == "" && yamlFileCount > 1 {
				return nil, fmt.Errorf("manifest contains %d .yaml files; specify a specific file in the URL fragment", yamlFileCount)
			}
			if specifiedFile == "" || title == specifiedFile {
				selectedLayer = &layer
				break
			}
		}
	}
	if selectedLayer == nil {
		return nil, fmt.Errorf("no matching .yaml file found in layers")
	}
	_, rc, err = remoteRepo.Blobs().FetchReference(ctx, selectedLayer.Digest.String())
	if err != nil {
		return nil, fmt.Errorf("blob fetch failed: %w", err)
	}
	defer rc.Close()
	buff, err := readLimited(rc, o.limit)
	if err != nil {
		return nil, fmt.Errorf("blob read failed: %w", err)
	}
	o.logger.Printf("Read %d bytes from %s", len(buff), specifiedFile)
	return buff, nil
}
