package serde_test

import (
	"testing"

	"github.com/jfk9w-go/flu/serde"
	"github.com/stretchr/testify/assert"
)

func TestSize_FromString(t *testing.T) {
	str := "100"
	size := new(serde.Size)
	assert.Nil(t, size.FromString(str))
	assert.Equal(t, int64(100), size.Bytes)

	str = "100b"
	assert.Nil(t, size.FromString(str))
	assert.Equal(t, int64(100), size.Bytes)

	str = "100Kb"
	assert.Nil(t, size.FromString(str))
	assert.Equal(t, int64(100<<10), size.Bytes)

	str = "100 Mb"
	assert.Nil(t, size.FromString(str))
	assert.Equal(t, int64(100<<20), size.Bytes)
}

func TestSize_ToString(t *testing.T) {
	size := serde.Size{Bytes: 100}
	assert.Equal(t, "100b", size.String())

	size.Bytes = 100 << 10
	assert.Equal(t, "100Kb", size.String())

	size.Bytes = 100 << 20
	assert.Equal(t, "100Mb", size.String())

	size.Bytes = 100<<30 + 100<<20
	assert.Equal(t, "100Gb", size.String())
}
