package uriget

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

func ExampleGetFile_http() {
	buff, err := GetFile(context.Background(), "http://example.com")
	fmt.Println(len(buff) > 0, err)
	_, err = GetFile(context.Background(), "https://example.com/does/not/exist")
	fmt.Println(err)
	// Output:
	// true <nil>
	// GET https://example.com/does/not/exist non-200 status code: 404
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
	testUrl := "oci://ghcr.io/score-spec/score-compose:0.18.0"
	buff, err := GetFile(context.Background(), testUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	fmt.Println(len(buff) > 0)
	// Output:
	// true
}

func ExampleGetFile_ociNoTag() {
	testUrl := "oci://ghcr.io/score-spec/score-compose"
	buff, err := GetFile(context.Background(), testUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	fmt.Println(len(buff) > 0)
	// Output:
	// true
}

func ExampleGetFile_ociWithDigest() {
	testUrl := "oci://ghcr.io/score-spec/score-compose@sha256:f3d8d5485a751cbdc91e073df1b6fbcde83f85a86ee3bc7d53e05b00452baedd"
	buff, err := GetFile(context.Background(), testUrl)
	if err != nil {
		fmt.Println("failed to pull OCI image:", err)
		return
	}
	fmt.Println(len(buff) > 0)
	// Output:
	// true
}
