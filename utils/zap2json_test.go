package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestZap2JsonReturnedBytesRemainValid(t *testing.T) {
	first, err := Zap2Json([]zap.Field{zap.String("event", "first")})
	require.NoError(t, err)
	want := append([]byte(nil), first...)

	for i := 0; i < 128; i++ {
		_, err := Zap2Json([]zap.Field{zap.String("event", "later")})
		require.NoError(t, err)
	}

	require.True(t, bytes.Equal(want, first), "public bytes must remain valid after later encodes")
}
