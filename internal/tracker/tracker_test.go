// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package tracker

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func setup(t *testing.T) {
	t.Helper()

	defaults.EtcPath = "testdata"

	agents := registration.Agents{
		registration.Agent{
			AgentID:     "abc123",
			AgentTypeID: "test1",
		},
		registration.Agent{
			AgentID:     "def456",
			AgentTypeID: "test2",
		},
	}

	data, err := yaml.Marshal(agents)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join("testdata", "agents.yaml"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

func createTest1Conf(cfgFile string) error {
	return os.WriteFile(cfgFile, baseConfData(), 0600) //nolint:wrapcheck
}

func baseConfData() []byte {
	return []byte("test:1")
}

func updateTest1Conf(cfgFile string) error {
	return os.WriteFile(cfgFile, updatedConfData(), 0600) //nolint:wrapcheck
}

func updatedConfData() []byte {
	return []byte("test:100")
}

func Test_getBasePath(t *testing.T) {
	setup(t)

	tests := []struct {
		name    string
		agentID string
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			agentID: "test1",
			want:    filepath.Join("testdata", "configs", "test1"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBasePath(tt.agentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBasePath() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("getBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTrackerFile(t *testing.T) {
	setup(t)

	type args struct {
		agentName string
		cfgFile   string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				agentName: "test1",
				cfgFile:   filepath.Join("testdata", "test1.conf"),
			},
			want:    filepath.Join("testdata", "configs", "test1", "test1.conf.current.yaml"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTrackerFile(tt.args.agentName, tt.args.cfgFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTrackerFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("getTrackerFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateConfig(t *testing.T) {
	setup(t)

	if err := createTest1Conf(filepath.Join("testdata", "test1.conf")); err != nil {
		t.Fatal(err)
	}

	type args struct {
		agentName string
		configID  string
		cfgFile   string
		data      []byte
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				agentName: "test1",
				configID:  "123",
				cfgFile:   filepath.Join("testdata", "test1.conf"),
				data:      []byte("test:1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateConfig(tt.args.agentName, tt.args.configID, tt.args.cfgFile, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyConfig(t *testing.T) {
	setup(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL.String())
		switch r.URL.String() {
		case "/agent/abc123/config_assignment/123":
			t.Log("ack modified config")
			_, _ = io.WriteString(w, "all good")

			return
		default:
			http.Error(w, "not found", http.StatusNotFound)

			return
		}
	}))
	defer ts.Close()

	viper.Set(keys.APIURL, ts.URL)
	viper.Set(keys.APIToken, "abc123")

	if err := createTest1Conf(filepath.Join("testdata", "test1.conf")); err != nil {
		t.Fatal(err)
	}

	if err := UpdateConfig("test1", "123", filepath.Join("testdata", "test1.conf"), []byte("test:1")); err != nil {
		t.Fatal(err)
	}

	if err := updateTest1Conf(filepath.Join("testdata", "test1.conf")); err != nil {
		t.Fatal(err)
	}

	if err := VerifyConfig(context.Background(), "test1", filepath.Join("testdata", "test1.conf")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
