package servicelayer

import (
	"github.com/jfrog/live-logs/internal/constants"
	"os"
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

func Test_servicelayer_checkVersion(t *testing.T) {
	tests := []struct {
		name            string
		want            string
		currentVersion  string
		minVersion      string
		wantErr         bool
		productName     string
		versionCheck    string
	}{
		{
			name: "equal_to_version",
			currentVersion: "7.16.0",
			minVersion: "7.16.0",
			productName: "artifactory",
			versionCheck: "true",
			wantErr:       false,
		},
		{
			name: "greater_than_version",
			currentVersion: "7.18.0",
			minVersion: "7.16.0",
			productName: "artifactory",
			versionCheck: "true",
			wantErr:       false,
		},
		{
			name: "less_than_version",
			currentVersion: "6.23.3",
			minVersion: "7.16.0",
			want: "minimum supported version is 7.16.0",
			productName: "artifactory",
			versionCheck: "true",
			wantErr:       true,
		},
		{
			name: "less_than_version_with_check_disabled",
			currentVersion: "6.23.3",
			minVersion: "7.16.0",
			productName: "artifactory",
			versionCheck: "false",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempEnvValue := os.Getenv(constants.VersionCheckEnv)

			os.Setenv(constants.VersionCheckEnv, tt.versionCheck)
			err := checkVersion(tt.currentVersion, tt.minVersion, tt.productName)
			os.Setenv(constants.VersionCheckEnv, tempEnvValue)

			if (err == nil) && tt.wantErr {
				t.Errorf("checkVersion() expecting error, no error recieved")
				return
			}
			if (err != nil) && !tt.wantErr {
				t.Errorf("checkVersion() expected no error, recieved err: %s",err.Error())
				return
			}
			if (err != nil) && tt.wantErr && !strings.Contains(err.Error(),tt.want) {
				t.Errorf("checkVersion() error = %s, want %s", err.Error(), tt.want)
				return
			}
		})
	}
}
