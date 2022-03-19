package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestCache_Success(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	cch := NewRedisCache(mr.Addr(), "")
	require.NoError(t, err)

	err = cch.Set(context.Background(), "aNjR4bEV", "http://google.com")
	require.NoError(t, err)

	url, err := cch.Get(context.Background(), "aNjR4bEV")
	require.NoError(t, err)
	require.Equal(t, url, "http://google.com")
}

func TestCache_FailIfNoKeyPresent(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	cch := NewRedisCache(mr.Addr(), "")
	require.NoError(t, err)

	_, err = cch.Get(context.Background(), "aNjR4bEV")
	require.Error(t, err)
	require.Equal(t, err.Error(), redis.Nil.Error())
}
