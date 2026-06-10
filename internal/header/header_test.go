package header

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: valid spacing header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069            \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 35, n)
	assert.False(t, done)

	//Test: Check for ending ~ done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Test: Invalid name header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Different Header line to check validity
	headers = NewHeaders()
	data = []byte("Content-Length: 55 \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "55", headers["content-length"])
	assert.False(t, done)

	//Headers with same data eg: Name: Tack , Name:Pranaya
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go \r\n")
	data1 := []byte("Set-Person: prime-loves-zig \r\n")
	data2 := []byte("Set-Person: tj-loves-ocaml \r\n")
	n, done, err = headers.Parse(data)
	n, done, err = headers.Parse(data1)
	n, done, err = headers.Parse(data2)
	require.NoError(t, err)
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.False(t, done)

}
