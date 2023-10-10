// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package inventory

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func testAPIInventoryFileName() string {
	return filepath.Join("testdata", "api_inventory.json")
}
func TestParseAPIAgents(t *testing.T) {
	tests := []struct {
		want    Agents
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid",
			file:    testAPIInventoryFileName(),
			want:    Agents{"linux": map[string]Agent{"telegraf": {Binary: "telegraf", Start: "", Stop: "", Restart: "", Reload: "", Status: "", Version: "", ConfigFiles: map[string]string{"d81c7650-19ae-4bf3-98df-5d24d53f5756": "/etc/telegraf/telegraf.conf"}}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("reading data: %s", err)
			}

			got, err := ParseAPIAgents(data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseAPIAgents() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseAPIAgents() = %v, want %v", got, tt.want)
			}
		})
	}
}
