// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	_GOOGLE_API_SCOPE = []string{
		"https://www.googleapis.com/auth/service.management.readonly",
	}
	tokenCache = &oauth2.Token{}
)

func GenerateAccessTokenFromFile(serviceAccountKey string) (string, time.Duration, error) {
	data, err := ioutil.ReadFile(serviceAccountKey)
	if err != nil {
		return "", 0, err
	}

	return generateAccessToken(data)
}

func generateAccessToken(keyData []byte) (string, time.Duration, error) {
	now := time.Now()
	// Follow the similar logic as GCE metadata server, where returned token will be valid for at
	// least 60s.
	if tokenCache.AccessToken != "" && !now.After(tokenCache.Expiry.Add(-time.Second*60)) {
		return tokenCache.AccessToken, tokenCache.Expiry.Sub(now), nil
	}
	creds, err := google.CredentialsFromJSON(oauth2.NoContext, keyData, _GOOGLE_API_SCOPE...)
	if err != nil {
		return "", 0, err
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", 0, err
	}
	tokenCache = token
	return token.AccessToken, token.Expiry.Sub(time.Now()), nil
}
