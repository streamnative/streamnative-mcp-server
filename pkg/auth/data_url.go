// Copyright (c) 2025, Inc.. All Rights Reserved.

package auth

import (
	"encoding/base64"
	"errors"
	"regexp"
)

var errDataURLInvalid = errors.New("invalid data URL")

// https://datatracker.ietf.org/doc/html/rfc2397
var dataURLRegex = regexp.MustCompile("^data:(?P<mimetype>[^;,]+)?(;(?P<charset>charset=[^;,]+))?" +
	"(;(?P<base64>base64))?,(?P<data>.+)")

type dataURL struct {
	url      string
	Mimetype string
	Data     []byte
}

func newDataURL(url string) (*dataURL, error) {
	if !dataURLRegex.Match([]byte(url)) {
		return nil, errDataURLInvalid
	}

	match := dataURLRegex.FindStringSubmatch(url)
	if len(match) != 7 {
		return nil, errDataURLInvalid
	}

	dataURL := &dataURL{
		url: url,
	}

	mimetype := match[dataURLRegex.SubexpIndex("mimetype")]
	if mimetype == "" {
		mimetype = "text/plain"
	}
	dataURL.Mimetype = mimetype

	data := match[dataURLRegex.SubexpIndex("data")]
	if match[dataURLRegex.SubexpIndex("base64")] == "" {
		dataURL.Data = []byte(data)
	} else {
		data, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}
		dataURL.Data = data
	}

	return dataURL, nil
}
