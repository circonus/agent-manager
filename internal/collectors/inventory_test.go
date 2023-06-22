package collectors

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
	"runtime"
	"testing"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func inventoryFileName() string {
	return filepath.Join("testdata", "inventory.yaml")
}
func binaryFileName() string {
	return filepath.Join("testdata", "test_binary")
}
func ConfFileName() string {
	return filepath.Join("testdata", "test_conf")
}
func setupTest() {
	file := inventoryFileName()
	c := Collectors{
		runtime.GOOS: map[string]Collector{
			"foo": {
				Binary: binaryFileName(),
				Start:  "start foo",
				Stop:   "stop foo",
				ConfigFiles: map[string]string{
					"000": ConfFileName(),
				},
			},
		},
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		log.Fatal("yaml marshal", err)
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		log.Fatal("write inv file", err)
	}
}

func TestFetchCollectors(t *testing.T) {
	testAuthToken := "foo"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/collector_type":
			switch r.Method {
			case "GET":
				authToken := r.Header.Get("X-Circonus-Auth-Token")
				if authToken != testAuthToken {
					http.Error(w, "invalid auth token", http.StatusUnauthorized)
					return
				}

				c := APICollectors{
					APICollector{
						Platforms: []Platform{
							{
								CollectorPlatformID: runtime.GOOS,
								CollectorTypeID:     "foo",
								Executable:          filepath.Join("testdata", "test_binary"),
								Start:               "start foo",
								Stop:                "stop foo",
								ConfigFiles: []ConfigFile{
									{
										ConfigFileID: "000",
										Path:         filepath.Join("testdata", "test_conf"),
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
			if err := FetchCollectors(context.Background()); (err != nil) != tt.wantErr {
				t.Fatalf("FetchCollectors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadCollectors(t *testing.T) {
	tests := []struct {
		want    Collectors
		invFile string
		name    string
		wantErr bool
	}{
		{
			name:    "valid",
			invFile: inventoryFileName(),
			want:    Collectors{runtime.GOOS: map[string]Collector{"foo": {ConfigFiles: map[string]string{"000": filepath.Join("testdata", "test_conf")}, Binary: filepath.Join("testdata", "test_binary"), Start: "start foo", Stop: "stop foo", Restart: "", Reload: "", Status: "", Version: ""}}},
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

			got, err := LoadCollectors()
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadCollectors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("LoadCollectors() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCheckForCollectors(t *testing.T) {
	testAuthToken := "foo"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/collector/agent":
			switch r.Method {
			case "POST":
				authToken := r.Header.Get("X-Circonus-Auth-Token")
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

				var collectors InstalledCollectors
				if err = json.Unmarshal(b, &collectors); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(""))
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

			if err := CheckForCollectors(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("CheckForCollectors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}