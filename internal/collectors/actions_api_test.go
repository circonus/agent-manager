// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func testAPIActionsFileName() string {
	return filepath.Join("testdata", "api_test_actions.json")
}

func TestParseAPIActions(t *testing.T) {
	setupTest()
	zerolog.SetGlobalLevel(zerolog.Disabled)
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

	for _, tt := range tests {
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
