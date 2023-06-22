package collectors

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseAPICollectors(t *testing.T) {
	tests := []struct {
		want    Collectors
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid",
			file:    "api_inventory.json",
			want:    Collectors{"linux": map[string]Collector{"telegraf": {Binary: "telegraf", Start: "", Stop: "", Restart: "", Reload: "", Status: "", Version: "", ConfigFiles: map[string]string{"d81c7650-19ae-4bf3-98df-5d24d53f5756": "/etc/telegraf/telegraf.conf"}}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tt.file))
			if err != nil {
				t.Fatalf("reading data: %s", err)
			}

			got, err := ParseAPICollectors(data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseAPICollectors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%#v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseAPICollectors() = %v, want %v", got, tt.want)
			}
		})
	}
}