package ratelimit

import (
	"fmt"
	"net"
	"strings"
)

func ParseTrustedProxies(trustedProxyEntries []string) ([]*net.IPNet, error) {
	trustedNetworks := make([]*net.IPNet, 0, len(trustedProxyEntries))
	for _, rawEntry := range trustedProxyEntries {
		trimmedEntry := strings.TrimSpace(rawEntry)
		if trimmedEntry == "" {
			continue
		}
		if !strings.Contains(trimmedEntry, "/") {
			parsedIP := net.ParseIP(trimmedEntry)
			if parsedIP == nil {
				return nil, fmt.Errorf("invalid trusted proxy %q", trimmedEntry)
			}
			if parsedIP.To4() != nil {
				trimmedEntry += "/32"
			} else {
				trimmedEntry += "/128"
			}
		}
		_, parsedNetwork, parseError := net.ParseCIDR(trimmedEntry)
		if parseError != nil {
			return nil, fmt.Errorf("invalid trusted proxy %q: %w", trimmedEntry, parseError)
		}
		trustedNetworks = append(trustedNetworks, parsedNetwork)
	}
	return trustedNetworks, nil
}

func IsTrustedProxy(trustedNetworks []*net.IPNet, candidateIP net.IP) bool {
	for _, trustedNetwork := range trustedNetworks {
		if trustedNetwork.Contains(candidateIP) {
			return true
		}
	}
	return false
}

func ResolveClientIP(remoteAddress string, forwardedFor string, realIP string, trustedNetworks []*net.IPNet) string {
	remoteHost, _, splitError := net.SplitHostPort(remoteAddress)
	if splitError != nil {
		remoteHost = remoteAddress
	}

	remoteIP := net.ParseIP(remoteHost)
	if remoteIP == nil || !IsTrustedProxy(trustedNetworks, remoteIP) {
		return remoteHost
	}

	if forwardedFor != "" {
		forwardedHops := strings.Split(forwardedFor, ",")
		for hopIndex := len(forwardedHops) - 1; hopIndex >= 0; hopIndex-- {
			hop := strings.TrimSpace(forwardedHops[hopIndex])
			parsedHop := net.ParseIP(hop)
			if parsedHop == nil {
				break
			}
			if !IsTrustedProxy(trustedNetworks, parsedHop) {
				return hop
			}
		}
	}

	trimmedRealIP := strings.TrimSpace(realIP)
	if trimmedRealIP != "" && net.ParseIP(trimmedRealIP) != nil {
		return trimmedRealIP
	}

	return remoteHost
}
