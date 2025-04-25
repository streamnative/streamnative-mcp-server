// Copyright (c) 2020 StreamNative, Inc.. All Rights Reserved.

package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	KeyFileTypeServiceAccount = "sn_service_account"
	FILE                      = "file://"
	DATA                      = "data://"
)

type KeyFileProvider struct {
	KeyFile string
}

type KeyFile struct {
	Type         string `json:"type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ClientEmail  string `json:"client_email"`
}

func NewClientCredentialsProviderFromKeyFile(keyFile string) *KeyFileProvider {
	return &KeyFileProvider{
		KeyFile: keyFile,
	}
}

var _ ClientCredentialsProvider = &KeyFileProvider{}

func (k *KeyFileProvider) GetClientCredentials() (*KeyFile, error) {
	var keyFile []byte
	var err error
	switch {
	case strings.HasPrefix(k.KeyFile, FILE):
		filename := strings.TrimPrefix(k.KeyFile, FILE)
		keyFile, err = os.ReadFile(filename)
	case strings.HasPrefix(k.KeyFile, DATA):
		keyFile = []byte(strings.TrimPrefix(k.KeyFile, DATA))
	case strings.HasPrefix(k.KeyFile, "data:"):
		url, err := newDataURL(k.KeyFile)
		if err != nil {
			return nil, err
		}
		keyFile = url.Data
	default:
		keyFile, err = os.ReadFile(k.KeyFile)
	}
	if err != nil {
		return nil, err
	}

	var v KeyFile
	err = json.Unmarshal(keyFile, &v)
	if err != nil {
		return nil, err
	}
	if v.Type != KeyFileTypeServiceAccount {
		return nil, fmt.Errorf("open %s: unsupported format", k.KeyFile)
	}

	return &v, nil
}

type KeyFileStructProvider struct {
	KeyFile *KeyFile
}

func (k *KeyFileStructProvider) GetClientCredentials() (*KeyFile, error) {
	if k.KeyFile == nil {
		return nil, fmt.Errorf("key file is nil")
	}
	return k.KeyFile, nil
}

func NewClientCredentialsProviderFromKeyFileStruct(keyFile *KeyFile) *KeyFileStructProvider {
	return &KeyFileStructProvider{
		KeyFile: keyFile,
	}
}

var _ ClientCredentialsProvider = &KeyFileStructProvider{}
