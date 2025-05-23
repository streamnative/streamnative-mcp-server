// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"k8s.io/utils/clock"
)

const (
	ClaimNameUserName = "https://streamnative.io/username"
)

// Flow abstracts an OAuth 2.0 authentication and authorization flow
type Flow interface {
	// Authorize obtains an authorization grant based on an OAuth 2.0 authorization flow.
	// The method returns a grant which may contain an initial access token.
	Authorize() (*AuthorizationGrant, error)
}

// AuthorizationGrantRefresher refreshes OAuth 2.0 authorization grant
type AuthorizationGrantRefresher interface {
	// Refresh refreshes an authorization grant to contain a fresh access token
	Refresh(grant *AuthorizationGrant) (*AuthorizationGrant, error)
}

type AuthorizationGrantType string

const (
	// GrantTypeClientCredentials represents a client credentials grant
	GrantTypeClientCredentials AuthorizationGrantType = "client_credentials"
)

// AuthorizationGrant is a credential representing the resource owner's authorization
// to access its protected resources, and is used by the client to obtain an access token
type AuthorizationGrant struct {
	// Type describes the type of authorization grant represented by this structure
	Type AuthorizationGrantType `json:"type"`

	// ClientCredentials is credentials data for the client credentials grant type
	ClientCredentials *KeyFile `json:"client_credentials,omitempty"`

	// Token contains an access token in the client credentials grant type,
	// and a refresh token in the device authorization grant type
	Token *oauth2.Token `json:"token,omitempty"`
}

// TokenResult holds token information
type TokenResult struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// Issuer holds information about the issuer of tokens
type Issuer struct {
	IssuerEndpoint string
	ClientID       string
	Audience       string
}

func convertToOAuth2Token(token *TokenResult, clock clock.Clock) oauth2.Token {
	return oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       clock.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}
}

// ExtractUserName extracts the username claim from an authorization grant
func ExtractUserName(token oauth2.Token) (string, error) {
	p := jwt.Parser{}
	claims := jwt.MapClaims{}

	// First try to extract from access token
	_, _, err := p.ParseUnverified(token.AccessToken, claims)

	// If access token parsing fails, try to extract from id_token
	if err != nil {
		// If id_token exists, try to parse it
		idToken, ok := token.Extra("id_token").(string)
		if ok && idToken != "" {
			if _, _, err := p.ParseUnverified(idToken, claims); err == nil {
				// Successfully parsed id_token
				if username, ok := claims[ClaimNameUserName].(string); ok {
					return username, nil
				}

				// Try other standard claims
				if email, ok := claims["email"].(string); ok {
					return email, nil
				}
				if sub, ok := claims["sub"].(string); ok {
					return sub, nil
				}
			}
		}

		// If unable to parse any token, return error
		return "", fmt.Errorf("unable to decode the token: %v", err)
	}

	// If parsing successful, try to get username
	if username, ok := claims[ClaimNameUserName]; ok {
		if v, ok := username.(string); ok {
			return v, nil
		}
		return "", fmt.Errorf("access token contains an unsupported username claim")
	}

	// Try other standard claims
	if email, ok := claims["email"].(string); ok {
		return email, nil
	}
	if sub, ok := claims["sub"].(string); ok {
		return sub, nil
	}

	return "", fmt.Errorf("access token doesn't contain a recognizable user claim")
}

func DumpToken(out io.Writer, token oauth2.Token) {
	p := jwt.Parser{}
	claims := jwt.MapClaims{}
	if _, _, err := p.ParseUnverified(token.AccessToken, claims); err != nil {
		fmt.Fprintf(out, "Unable to parse token.  Err: %v\n", err)
		return
	}

	text, err := json.MarshalIndent(claims, "", "  ")
	if err != nil {
		fmt.Fprintf(out, "Unable to print token.  Err: %v\n", err)
	}
	_, _ = out.Write(text)
	_, _ = fmt.Fprintln(out, "")
}
