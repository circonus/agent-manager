package registration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func Test_getJWT(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/registration":
			switch r.Method {
			case "GET":
				http.Error(w, "not found", http.StatusNotFound)
				return
			case "POST":
				regToken := r.Header.Get("X-Circonus-Reg-Token")
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

				var request Request
				if err = json.Unmarshal(b, &request); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if request.Meta.MachineID == "" {
					http.Error(w, "bad machine id", http.StatusBadRequest)
					return
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					Subject: request.Meta.MachineID,
				})
				tokenString, err := token.SignedString([]byte("secret"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"token":"` + tokenString + `"}`))
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
		token   string
		request Request
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "invalid (no token)",
			args:    args{token: "", request: Request{Meta: Meta{}}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid (no subject)",
			args:    args{token: "foo", request: Request{Meta: Meta{}}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid",
			args:    args{token: "foo", request: Request{Meta: Meta{MachineID: "bar"}}},
			want:    []byte(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiYXIifQ.ST4yLHEt-g5qTE6NW5gAp6omAfVezv8dwUPTVtM2rKs"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJWT(context.Background(), tt.args.token, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getJWT() = '%v', want '%v'", string(got), string(tt.want))
			}
		})
	}
}
