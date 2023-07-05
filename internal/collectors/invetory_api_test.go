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
)

func testAPIInventoryFileName() string {
	return filepath.Join("testdata", "api_inventory.json")
}
func TestParseAPICollectors(t *testing.T) {
	tests := []struct {
		want    Collectors
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid",
			file:    testAPIInventoryFileName(),
			want:    Collectors{"linux": map[string]Collector{"telegraf": {Binary: "telegraf", Start: "", Stop: "", Restart: "", Reload: "", Status: "", Version: "", ConfigFiles: map[string]string{"d81c7650-19ae-4bf3-98df-5d24d53f5756": "/etc/telegraf/telegraf.conf"}}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("reading data: %s", err)
			}

			got, err := ParseAPICollectors(data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseAPICollectors() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseAPICollectors() = %v, want %v", got, tt.want)
			}
		})
	}
}
