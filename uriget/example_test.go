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
