package main

import (
	"net/url"
	"strings"
)

func normalizeURL(input string) (string, error) {
	url, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	fullpath := url.Host + url.Path
	fullpath = strings.ToLower(strings.TrimSuffix(fullpath, "/"))

	return fullpath, nil
}
