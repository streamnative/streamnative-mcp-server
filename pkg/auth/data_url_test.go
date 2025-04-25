// Copyright (c) 2025 StreamNative, Inc.. All Rights Reserved.

package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataURL(t *testing.T) {
	rawURL := "data:,test"
	url, err := newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", url.Mimetype)
	assert.Equal(t, "test", string(url.Data))

	rawURL = "data:;base64," + base64.StdEncoding.EncodeToString([]byte("test"))
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", url.Mimetype)
	assert.Equal(t, "test", string(url.Data))

	rawURL = "data:application/json,test"
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "application/json", url.Mimetype)
	assert.Equal(t, "test", string(url.Data))

	rawURL = "data:application/json;base64," + base64.StdEncoding.EncodeToString([]byte("test"))
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "application/json", url.Mimetype)
	assert.Equal(t, "test", string(url.Data))

	rawURL = "data://test"
	url, err = newDataURL(rawURL)
	require.Nil(t, url)
	assert.Error(t, err)
	assert.EqualError(t, errDataURLInvalid, err.Error())
}

func TestNewDataURLEdgeCases(t *testing.T) {
	rawURL := "data:,"
	url, err := newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", url.Mimetype)
	assert.Equal(t, "", string(url.Data))

	rawURL = "data:;,test"
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", url.Mimetype)
	assert.Equal(t, "test", string(url.Data))

	rawURL = "data:;base64,@@invalid@@"
	_, err = newDataURL(rawURL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "illegal base64 data")

	rawURL = "data:application/json;charset=utf-8,{\"test\":true}"
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "application/json", url.Mimetype)
	assert.Equal(t, "{\"test\":true}", string(url.Data))

	rawURL = "data:text/plain,hello%20world"
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", url.Mimetype)
	assert.Equal(t, "hello%20world", string(url.Data))

	rawURL = "data:;base64," + base64.StdEncoding.EncodeToString([]byte(""))
	url, err = newDataURL(rawURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", url.Mimetype)
	assert.Equal(t, "", string(url.Data))
}
