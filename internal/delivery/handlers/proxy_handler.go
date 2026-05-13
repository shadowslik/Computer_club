package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxyHandler(targetURL string) (http.Handler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(target), nil
}
