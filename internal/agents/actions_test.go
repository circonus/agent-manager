// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/spf13/viper"
)

const (
	testAuthToken = "foo"
)

func Test_getActions(t *testing.T) {
	setupTest()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/agent/update":
			switch r.Method {
			case http.MethodGet:
				authToken := r.Header.Get("Authorization")
				if authToken != testAuthToken {
					http.Error(w, "invalid auth token", http.StatusUnauthorized)

					return
				}

				a := APIActions{
					APIAction{
						ConfigAssignmentID: "c3ef3233-2792-48be-aaab-745aaf02f5e9",
						Config: APIConfig{
							FileID:   "d81c7650-19ae-4bf3-98df-5d24d53f5756",
							Contents: "dGVzdAo=",
						},
						Agent: APIConfigAgent{
							ID: "foo",
						},
					},
				}

				data, err := json.Marshal(a)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(data)

			case http.MethodPost:
				authToken := r.Header.Get("Authorization")
				if authToken != testAuthToken {
					http.Error(w, "invalid auth token", http.StatusUnauthorized)

					return
				}
				defer r.Body.Close()
				b, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				var c ConfigResult
				if err := json.Unmarshal(b, &c); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			default:
				http.Error(w, "not found", http.StatusNotFound)

				return
			}
		default:
			http.Error(w, "not found", http.StatusNotFound)

			return
		}
	}))

	defer ts.Close()

	tests := []struct {
		name     string
		reqURL   string
		apiToken string
		invFile  string
		wantErr  bool
	}{
		{
			name:     "valid",
			reqURL:   ts.URL,
			apiToken: testAuthToken,
			invFile:  inventoryFileName(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.APIURL, tt.reqURL)
			viper.Set(keys.APIToken, tt.apiToken)
			viper.Set(keys.InventoryFile, tt.invFile)

			if err := getActions(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("getActions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
