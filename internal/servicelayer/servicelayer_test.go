package servicelayer

import (
	"strings"
	"testing"
)

func Test_LiveLogs_errorHandle(t *testing.T) {
	tests := []struct {
		want            string
		statusCode      int
		resBody         []byte
		wantErr         bool
		name            string
	}{
		{
			name: "status 200",
			statusCode: 200,
			resBody: []byte("success"),
		},
		{
			name: "status 400",
			want:    "error_400_content",
			resBody: []byte("error_400_content"),
			wantErr: true,
			statusCode: 400,
		},
		{
			name: "status 429",
			want:    "error_429_content",
			resBody: []byte("error_429_content"),
			wantErr: true,
			statusCode: 429,
		},
		{
			name: "status 500",
			want:    "unexpected response",
			resBody: []byte("error_500_content"),
			wantErr: true,
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorHandle(tt.statusCode, tt.resBody)
			if (err != nil) && !strings.Contains(err.Error(),tt.want) {
				t.Errorf("errorHandle() error = %v, wantErr %v", err.Error(), tt.want)
				return
			}
		})
	}
}
