package interceptor

import (
	"context"
	"net"
	"testing"

	"github.com/healthcare/backend/internal/shared/ratelimit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func buildClientIPContext(peerAddress string, headers map[string]string) context.Context {
	ctx := context.Background()
	if peerAddress != "" {
		tcpAddress, _ := net.ResolveTCPAddr("tcp", peerAddress)
		ctx = peer.NewContext(ctx, &peer.Peer{Addr: tcpAddress})
	}
	if len(headers) > 0 {
		ctx = metadata.NewIncomingContext(ctx, metadata.New(headers))
	}
	return ctx
}

func mustParseTrustedNetworks(testingInstance *testing.T, entries ...string) []*net.IPNet {
	trustedNetworks, parseError := ratelimit.ParseTrustedProxies(entries)
	require.NoError(testingInstance, parseError)
	return trustedNetworks
}

func TestExtractClientIPIgnoresSpoofedHeaderFromUntrustedPeer(testingInstance *testing.T) {
	ctx := buildClientIPContext("198.51.100.20:4000", map[string]string{"x-forwarded-for": "1.2.3.4"})

	assert.Equal(testingInstance, "198.51.100.20", extractClientIP(ctx, nil))
	assert.Equal(testingInstance, "198.51.100.20", extractClientIP(ctx, mustParseTrustedNetworks(testingInstance, "10.0.0.0/8")))
}

func TestExtractClientIPDefaultTrustsNoProxy(testingInstance *testing.T) {
	ctx := buildClientIPContext("10.0.0.5:4000", map[string]string{"x-forwarded-for": "203.0.113.7"})

	assert.Equal(testingInstance, "10.0.0.5", extractClientIP(ctx, nil))
}

func TestExtractClientIPHonorsForwardedChainFromTrustedProxy(testingInstance *testing.T) {
	trustedNetworks := mustParseTrustedNetworks(testingInstance, "10.0.0.0/8")
	ctx := buildClientIPContext("10.0.0.5:4000", map[string]string{"x-forwarded-for": "1.2.3.4, 203.0.113.7, 10.0.0.9"})

	assert.Equal(testingInstance, "203.0.113.7", extractClientIP(ctx, trustedNetworks))
}

func TestExtractClientIPFallsBackToPeerWithoutHeader(testingInstance *testing.T) {
	trustedNetworks := mustParseTrustedNetworks(testingInstance, "10.0.0.0/8")
	ctx := buildClientIPContext("10.0.0.5:4000", nil)

	assert.Equal(testingInstance, "10.0.0.5", extractClientIP(ctx, trustedNetworks))
}

func TestExtractClientIPUnknownWithoutPeer(testingInstance *testing.T) {
	assert.Equal(testingInstance, "unknown", extractClientIP(context.Background(), nil))
}
