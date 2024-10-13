// Package uriget provides a mechanism for "loading" the contents of a file from a uri with flexible support for
// different uri schemes. This is primarily used for loading Score provisioner files for various Score implementations.
// There are similar packages such as hashicorp/go-getter, however this version is maintained to have zero dependencies
// and use the binaries that exist on the local system where possible.
package uriget

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry/remote"
)

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
	case "http":
		fallthrough
	case "https":
		return opts.getHttp(ctx, u)
	case "file":
		fallthrough
	case "":
		return opts.getFile(ctx, u)
	case "git-ssh":
		fallthrough
	case "git-https":
		return opts.getGit(ctx, u)
	case "oci":
		return opts.getOci(ctx, u)
	default:
		return nil, fmt.Errorf("unsupported scheme '%s'", u.Scheme)
	}
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

	td, err := os.MkdirTemp(os.TempDir(), "score-go")
	if err != nil {
		return nil, fmt.Errorf("failed to make temp dir")
	} else if err := os.Chmod(td, 0700); err != nil {
		return nil, fmt.Errorf("failed to chown temp dir")
	}
	defer func() {
		_ = os.RemoveAll(td)
	}()

	gitBinary, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("failed to find git binary on the local system: %w", err)
	}
	gitRemote := "origin"
	getRef := "HEAD"

	c := exec.CommandContext(ctx, gitBinary, "init")
	c.Dir = td
	if output, err := c.CombinedOutput(); err != nil {
		o.logger.Printf("command output: %s", output)
		return nil, fmt.Errorf("failed to init git repo in %s: %w", td, err)
	}
	c = exec.CommandContext(ctx, gitBinary, "remote", "add", gitRemote, u.String())
	c.Dir = td
	if output, err := c.CombinedOutput(); err != nil {
		o.logger.Printf("command output: %s", output)
		return nil, fmt.Errorf("failed to set git remote to %s: %w", u.String(), err)
	}
	o.logger.Printf("Initialized git remote in %s for %s", td, u.String())
	// https://stackoverflow.com/questions/61587133/cloning-single-file-from-git-repository
	c = exec.CommandContext(ctx, gitBinary, "sparse-checkout", "set", "--no-cone", "--sparse-index", subPath)
	c.Dir = td
	if output, err := c.CombinedOutput(); err != nil {
		o.logger.Printf("command output: %s", output)
		return nil, fmt.Errorf("failed to set sparse checkout: %w", err)
	}
	c = exec.CommandContext(ctx, gitBinary, "pull", gitRemote, getRef, "--depth=1")
	c.Dir = td
	if output, err := c.CombinedOutput(); err != nil {
		o.logger.Printf("command output: %s", output)
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
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

func (o *options) getOci(ctx context.Context, u *url.URL) ([]byte, error) {
	parts := strings.Split(u.Host+u.Path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid OCI URL format")
	}
	registry := parts[0]
	repo := strings.Join(parts[1:], "/")
	tag := "latest"
	if u.Fragment != "" {
		tag = u.Fragment
	}
	store, err := oci.New(o.tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCI layout store: %w", err)
	}
	repoUrl := fmt.Sprintf("%s/%s", registry, repo)
	remoteRepo, err := remote.NewRepository(repoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote repository: %w", err)
	}
	manifestDescriptor, err := oras.Copy(ctx, remoteRepo, tag, store, tag, oras.DefaultCopyOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to pull OCI image: %w", err)
	}
	o.logger.Printf("Pulled OCI image: %s with manifest descriptor : %v", u.String(), manifestDescriptor.Digest)
	return []byte(manifestDescriptor.Digest), nil
}
