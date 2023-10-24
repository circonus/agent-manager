// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package inventory

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/platform"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	testAuthToken = "foo"
)

var initialized = false //nolint:gochecknoglobals

func inventoryFileName() string {
	return filepath.Join("testdata", "inventory.yaml")
}
func binaryFileName() string {
	return filepath.Join("testdata", "test_binary")
}
func confFileID() string {
	return "d81c7650-19ae-4bf3-98df-5d24d53f5756"
}
func confFileName() string {
	return filepath.Join("testdata", "test_conf")
}
func setupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	file := inventoryFileName()
	aa := Agents{
		platform.Get(): map[string]Agent{
			"foo": {
				Binary: binaryFileName(),
				Start:  "start foo",
				Stop:   "stop foo",
				ConfigFiles: map[string]string{
					confFileID(): confFileName(),
				},
			},
		},
	}

	data, err := yaml.Marshal(aa)
	if err != nil {
		log.Fatal("yaml marshal", err)
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		log.Fatal("write inv file", err)
	}

	initialized = true
}

func TestFetchAgents(t *testing.T) {
	setupTest()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/agent_type":
			switch r.Method {
			case http.MethodGet:
				authToken := r.Header.Get("Authorization")
				if authToken != testAuthToken {
					http.Error(w, "invalid auth token", http.StatusUnauthorized)

					return
				}

				c := APIAgents{
					APIAgent{
						Platforms: []Platform{
							{
								ID:          platform.Get(),
								AgentTypeID: "foo",
								Executable:  binaryFileName(),
								Start:       "start foo",
								Stop:        "stop foo",
								ConfigFiles: []ConfigFile{
									{
										ConfigFileID: "d81c7650-19ae-4bf3-98df-5d24d53f5756",
										Path:         confFileName(),
									},
								},
							},
						},
					},
				}
				data, err := json.Marshal(c)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(data)
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
			wantErr:  false,
		},
		{
			name:     "invalid (url)",
			reqURL:   "",
			apiToken: testAuthToken,
			invFile:  inventoryFileName(),
			wantErr:  true,
		},
		{
			name:     "invalid (token)",
			reqURL:   ts.URL,
			apiToken: "",
			invFile:  inventoryFileName(),
			wantErr:  true,
		},
		{
			name:     "invalid (inv file)",
			reqURL:   ts.URL,
			apiToken: testAuthToken,
			invFile:  "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.APIURL, tt.reqURL)
			viper.Set(keys.APIToken, tt.apiToken)
			viper.Set(keys.InventoryFile, tt.invFile)
			if err := FetchAgents(context.Background()); (err != nil) != tt.wantErr {
				t.Fatalf("FetchAgents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAgents(t *testing.T) {
	setupTest()

	tests := []struct {
		want    Agents
		invFile string
		name    string
		wantErr bool
	}{
		{
			name:    "valid",
			invFile: inventoryFileName(),
			want:    Agents{platform.Get(): map[string]Agent{"foo": {ConfigFiles: map[string]string{confFileID(): confFileName()}, Binary: binaryFileName(), Start: "start foo", Stop: "stop foo", Restart: "", Reload: "", Status: "", Version: ""}}},
			wantErr: false,
		},
		{
			name:    "invalid (inv file)",
			invFile: "",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.InventoryFile, tt.invFile)

			got, err := LoadAgents()
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadAgents() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("LoadAgents() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCheckForAgents(t *testing.T) {
	setupTest()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/agent/manager":
			switch r.Method {
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

				var agents InstalledAgents

				if err = json.Unmarshal(b, &agents); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"agent_id":"abc123"}]`))
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
		{
			name:     "invalid (url)",
			reqURL:   "",
			apiToken: testAuthToken,
			invFile:  inventoryFileName(),
			wantErr:  true,
		},
		{
			name:     "invalid (token)",
			reqURL:   ts.URL,
			apiToken: "",
			invFile:  inventoryFileName(),
			wantErr:  true,
		},
		{
			name:     "invalid (inv file)",
			reqURL:   ts.URL,
			apiToken: testAuthToken,
			invFile:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.APIURL, tt.reqURL)
			viper.Set(keys.APIToken, tt.apiToken)
			viper.Set(keys.InventoryFile, tt.invFile)

			if err := CheckForAgents(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("CheckForAgents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
