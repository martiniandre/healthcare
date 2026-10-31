package cache

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectReturnsNilForEmptyUrl(testingInstance *testing.T) {
	assert.Nil(testingInstance, Connect("  "))
}

func TestConnectReturnsNilWhenRedisIsUnreachable(testingInstance *testing.T) {
	assert.Nil(testingInstance, Connect("127.0.0.1:1"))
}

func TestConnectAcceptsHostPortAndUrlForms(testingInstance *testing.T) {
	fakeRedis := miniredis.RunT(testingInstance)

	hostPortClient := Connect(fakeRedis.Addr())
	require.NotNil(testingInstance, hostPortClient)
	require.NoError(testingInstance, hostPortClient.Close())

	urlClient := Connect("redis://" + fakeRedis.Addr() + "/0")
	require.NotNil(testingInstance, urlClient)
	require.NoError(testingInstance, urlClient.Close())
}
