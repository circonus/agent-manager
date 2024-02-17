// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func Test_getJWT(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/manager/register":
			switch r.Method {
			case http.MethodGet:
				http.Error(w, "not found", http.StatusNotFound)

				return
			case http.MethodPost:
				regToken := r.Header.Get("Authorization")
				if regToken == "" {
					http.Error(w, "missing token", http.StatusUnauthorized)

					return
				}

				defer r.Body.Close()

				b, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				var reg Registration
				if err = json.Unmarshal(b, &reg); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				if reg.MachineID == "" {
					http.Error(w, "bad machine id", http.StatusBadRequest)

					return
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					Subject: reg.MachineID,
				})

				tokenString, err := token.SignedString([]byte("secret"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				r := Response{
					AuthToken:    tokenString,
					RefreshToken: "abc",
					ManagerID:    "test",
				}

				data, err := json.Marshal(r)
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
	viper.Set(keys.APIURL, ts.URL)

	type args struct {
		token string
		reg   Registration
	}

	tests := []struct {
		name    string
		want    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid (no token)",
			args:    args{token: "", reg: Registration{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid (no subject)",
			args:    args{token: "foo", reg: Registration{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "valid",
			args:    args{token: "foo", reg: Registration{MachineID: "bar", Hostname: "foo"}},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiYXIifQ.ST4yLHEt-g5qTE6NW5gAp6omAfVezv8dwUPTVtM2rKs",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJWT(context.Background(), tt.args.token, tt.args.reg)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJWT() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != nil {
				if got.AuthToken != tt.want {
					t.Errorf("getJWT() = '%v', want '%v'", got.AuthToken, tt.want)
				}
			}
		})
	}
}
