// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/env"
	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	aa := inventory.Agents{
		env.GetPlatform(): map[string]inventory.Agent{
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

	// zerolog.SetGlobalLevel(zerolog.Disabled)

	data, err := yaml.Marshal(aa)
	if err != nil {
		log.Fatal("yaml marshal", err)
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		log.Fatal("write inv file", err)
	}

	initialized = true
}

func testAPIActionsFileName() string {
	return filepath.Join("testdata", "api_test_actions.json")
}

func TestParseAPIActions(t *testing.T) {
	setupTest()
	// zerolog.SetGlobalLevel(zerolog.Disabled)

	tests := []struct {
		name    string
		actFile string
		invFile string
		want    Actions
		wantErr bool
	}{
		{
			name:    "valid",
			actFile: testAPIActionsFileName(),
			invFile: inventoryFileName(),
			want: Actions{
				Action{
					Configs: map[string][]Config{
						"foo": {
							{
								ID:       "c3ef3233-2792-48be-aaab-745aaf02f5e9",
								Path:     "testdata/test_conf",
								Contents: "dGVzdAo=",
							},
						},
					},
					ID:       "",
					Type:     "config",
					Commands: []Command(nil),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests { //nolint:varnamelen
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.InventoryFile, tt.invFile)
			data, err := os.ReadFile(tt.actFile)
			if err != nil {
				t.Fatalf("reading test file: %s", err)
			}
			got, err := ParseAPIActions(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAPIActions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAPIActions() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
