package uriget

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	testUrl := "oci://ghcr.io/score-spec/score-compose-community-provisioners:v0.1.0#00-service.provisioners.yaml"
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
	ociTestUrl := "oci://ghcr.io/score-spec/score-compose-community-provisioners:v0.1.0#00-service.provisioners.yaml"
	ociBuff, err := GetFile(context.Background(), ociTestUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	gitTestUrl := "git-https://github.com/score-spec/community-provisioners.git/score-compose/00-service.provisioners.yaml"
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
	ociTestUrl := "oci://ghcr.io/score-spec/score-compose-community-provisioners:v0.1.0#00-service.provisioners.yaml"
	ociBuff, err := GetFile(context.Background(), ociTestUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	httpsTestUrl := "https://github.com/score-spec/community-provisioners/raw/v0.1.0/score-compose/00-service.provisioners.yaml"
	httpsbuff, err := GetFile(context.Background(), httpsTestUrl)
	if err != nil {
		fmt.Println("failed to pull file by HTTPS:", err)
		return
	}
	fmt.Println(len(ociBuff) == len(httpsbuff))
	// Output:
	// true
}
