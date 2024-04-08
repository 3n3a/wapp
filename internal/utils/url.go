package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ada-url/goada"
)

// Normalize a Url by rewriting special characters to url standard
func UrlNormalize(inputUrl string) (string, error) {
	url, err := goada.New(inputUrl)
	if err != nil {
		return "", errors.New("url parsing failed")
	}

	return url.Href(), nil
}

// Return the Hostname part of a url
//
// See explanation here: https://url-parts.glitch.me/?url=https://cats.example.org.au:1234/stripes/fur.html?pattern=tabby#claws
// 
// Source for Explanation: https://web.dev/articles/url-parts
func UrlGetHostname(inputUrl string) (string, error) {
	// prepend http if only domain
	if !(strings.Contains(inputUrl, "://")) {
		inputUrl = fmt.Sprintf("http://%s", inputUrl)
	}

	url, err := goada.New(inputUrl)
	if err != nil {
		return "", errors.New("url parsing failed")
	}

	return url.Hostname(), nil
}


func IsURLInterface(input interface{}) bool {
	if IsString(input) {
		return IsUrl(input.(string))
	}
	return false
}

func IsUrl(input string) (bool) {
	_, err := goada.New(input)
	return err == nil
}