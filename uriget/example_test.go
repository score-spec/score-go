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

package uriget

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func ExampleGetFile_local() {
	buff, err := GetFile(context.Background(), "../README.md")
	fmt.Println(len(buff) > 0, err)
	_, err = GetFile(context.Background(), "./does/not/exist.txt")
	fmt.Println(err)
	// Output:
	// true <nil>
	// open ./does/not/exist.txt: no such file or directory
}

func TestGetStdinFile(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test_stdin_file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testData := "This is a test."
	if _, err := tempFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek in temporary file: %v", err)
	}

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin = tempFile

	ctx := context.Background()
	output, err := getStdinFile(ctx)
	if err != nil {
		t.Fatalf("Failed to read from stdin: %v", err)
	}
	if string(output) != testData {
		t.Errorf("Expected %s, but got %s", testData, string(output))
	}

}

func TestGetStdinFile_NoData(t *testing.T) {
	// Step 1: Backup the original os.Stdin
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }() // Restore original os.Stdin after the test

	// Step 2: Assign an empty os.Stdin
	tempFile, err := os.CreateTemp("", "test_stdin_empty")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	os.Stdin = tempFile // Changed line

	// Step 3: Call the function and verify the output
	ctx := context.Background()
	_, err = getStdinFile(ctx)
	if err == nil {
		t.Fatalf("Expected an error when no stdin data is provided, but got none")
	}
}

func ExampleGetFile_http() {
	buff, err := GetFile(context.Background(), "http://example.com")
	fmt.Println(len(buff) > 0, err)
	// Output:
	// true <nil>
}

func ExampleGetFile_git() {
	buff, err := GetFile(context.Background(), "git-https://github.com/score-spec/score.dev.git/README.md")
	fmt.Println(len(buff) > 0, err)
	// Output: true <nil>
}

func ExampleWithLimit() {
	_, err := GetFile(context.Background(), "../README.md", WithLimit(1))
	fmt.Println(err)
	// Output: failed to read file: 1 byte limit exceeded
}

func ExampleWithHttpClient() {
	customClient := &http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return nil, fmt.Errorf("no proxy")
			},
		},
	}
	_, err := GetFile(context.Background(), "https://example.com", WithHttpClient(customClient))
	fmt.Println(err)
	// Output: failed to make get request: Get "https://example.com": no proxy
}
func ExampleGetFile_oci() {
	testUrl := "oci://ghcr.io/score-spec/score-compose-community-provisioners:latest#10-service.provisioners.yaml"
	buff, err := GetFile(context.Background(), testUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	fmt.Println(len(buff) > 0)
	// Output:
	// true
}

func ExampleGetFile_oci_git() {
	ociTestUrl := "oci://ghcr.io/score-spec/score-compose-community-provisioners:latest#10-service.provisioners.yaml"
	ociBuff, err := GetFile(context.Background(), ociTestUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	gitTestUrl := "git-https://github.com/score-spec/community-provisioners.git/service/score-compose/10-service.provisioners.yaml"
	gitBuff, err := GetFile(context.Background(), gitTestUrl)
	if err != nil {
		fmt.Println("failed to pull file in git:", err)
		return
	}
	fmt.Println(len(ociBuff) == len(gitBuff))
	// Output:
	// true
}

func ExampleGetFile_oci_https() {
	ociTestUrl := "oci://ghcr.io/score-spec/score-compose-community-provisioners:latest#10-service.provisioners.yaml"
	ociBuff, err := GetFile(context.Background(), ociTestUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	httpsTestUrl := "https://github.com/score-spec/community-provisioners/raw/main/service/score-compose/10-service.provisioners.yaml"
	httpsbuff, err := GetFile(context.Background(), httpsTestUrl)
	if err != nil {
		fmt.Println("failed to pull file by HTTPS:", err)
		return
	}
	fmt.Println(len(ociBuff) == len(httpsbuff))
	// Output:
	// true
}

func ExampleGetFiles_local() {
	results, err := GetFiles(context.Background(), "../README.md")
	fmt.Println(len(results) == 1, len(results[0].Content) > 0, err)
	// Output:
	// true true <nil>
}

func TestGetFiles_SingleFile(t *testing.T) {
	td := t.TempDir()
	filePath := filepath.Join(td, "test.yaml")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	results, err := GetFiles(context.Background(), filePath, WithLogger(log.New(os.Stderr, "", 0)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if string(results[0].Content) != "hello" {
		t.Errorf("expected content 'hello', got '%s'", results[0].Content)
	}
}

func TestGetFiles_Directory(t *testing.T) {
	td := t.TempDir()
	// Create files in non-alphabetical order to verify sorting
	files := map[string]string{
		"c-provisioner.yaml": "content-c",
		"a-provisioner.yaml": "content-a",
		"b-provisioner.yaml": "content-b",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(td, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	results, err := GetFiles(context.Background(), td, WithLogger(log.New(os.Stderr, "", 0)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// Verify sorted order and full path URIs
	expected := []struct {
		uri     string
		content string
	}{
		{filepath.Join(td, "a-provisioner.yaml"), "content-a"},
		{filepath.Join(td, "b-provisioner.yaml"), "content-b"},
		{filepath.Join(td, "c-provisioner.yaml"), "content-c"},
	}
	for i, e := range expected {
		if results[i].URI != e.uri {
			t.Errorf("result[%d]: expected URI '%s', got '%s'", i, e.uri, results[i].URI)
		}
		if string(results[i].Content) != e.content {
			t.Errorf("result[%d]: expected content '%s', got '%s'", i, e.content, results[i].Content)
		}
	}
}

func TestGetFiles_DirectorySkipsSubdirs(t *testing.T) {
	td := t.TempDir()
	if err := os.WriteFile(filepath.Join(td, "file.yaml"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(td, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	results, err := GetFiles(context.Background(), td, WithLogger(log.New(os.Stderr, "", 0)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result (subdir skipped), got %d", len(results))
	}
	if results[0].URI != filepath.Join(td, "file.yaml") {
		t.Errorf("expected URI '%s', got '%s'", filepath.Join(td, "file.yaml"), results[0].URI)
	}
}

func TestGetFiles_EmptyDirectory(t *testing.T) {
	td := t.TempDir()
	_, err := GetFiles(context.Background(), td, WithLogger(log.New(os.Stderr, "", 0)))
	if err == nil {
		t.Fatal("expected error for empty directory, got nil")
	}
}

func TestGetFiles_NonExistentPath(t *testing.T) {
	_, err := GetFiles(context.Background(), "/does/not/exist", WithLogger(log.New(os.Stderr, "", 0)))
	if err == nil {
		t.Fatal("expected error for non-existent path, got nil")
	}
}

func TestGetFiles_DirectoryWithFileScheme(t *testing.T) {
	td := t.TempDir()
	if err := os.WriteFile(filepath.Join(td, "a.yaml"), []byte("aaa"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(td, "b.yaml"), []byte("bbb"), 0644); err != nil {
		t.Fatal(err)
	}
	results, err := GetFiles(context.Background(), "file://"+td, WithLogger(log.New(os.Stderr, "", 0)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].URI != filepath.Join(td, "a.yaml") || results[1].URI != filepath.Join(td, "b.yaml") {
		t.Errorf("unexpected URIs: %s, %s", results[0].URI, results[1].URI)
	}
}
