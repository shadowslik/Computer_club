package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if xrip := r.Header.Get("X-Real-Ip"); xrip != "" {
		return xrip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
