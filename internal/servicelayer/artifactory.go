package servicelayer

import (
	"context"
	"encoding/json"
	"fmt"
	cliCommands "github.com/jfrog/jfrog-cli-core/common/commands"
	"github.com/jfrog/live-logs/internal/clientlayer"
	"github.com/jfrog/live-logs/internal/constants"
	"github.com/jfrog/live-logs/internal/model"
	"os"
	"strings"
	"time"
)
const (
	defaultRequestTimeout    = 15 * time.Second
	defaultLogRequestTimeout = time.Minute
	defaultLogsRefreshRate   = time.Second
	artifactoryVersionEndPoint = "api/system/version"
	artifactoryMinVersionSupport = "7.16.0"
	artifactoryConfigEndpoint = "api/system/logs/config"
	artifactoryDataEndpoint   = "api/system/logs/data"
	artifactoryProductName = "Artifactory"
)

type ArtifactoryData struct {
	nodeId          string
	logFileName     string
	lastPageMarker  int64
	logsRefreshRate time.Duration
}

type artifactoryVersionData struct {
	Version string `json:"version,omitempty"`
}

func (s *ArtifactoryData) GetConfig(ctx context.Context, serverId string) (*model.Config, error) {

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelTimeout()

	err := s.ArtifactoryValidations(ctx, serverId)
	if err != nil {
		return nil, err
	}
	baseUrl, err := s.getUrl(serverId)
	if err != nil {
		return nil, err
	}
	res, resBody, err := clientlayer.SendGet(timeoutCtx, serverId, artifactoryConfigEndpoint,constants.EmptyNodeId, baseUrl, nil)
	if err != nil {
		return nil, err
	}

	err = errorHandle(res.StatusCode, resBody)
	if err != nil {
		return nil, err
	}

	logConfig := model.Config{}
	err = json.Unmarshal(resBody, &logConfig)
	if err != nil {
		return nil, err
	}
	if len(logConfig.LogFileNames) == 0 {
		return nil, fmt.Errorf("no log file names were found")
	}
	if len(logConfig.Nodes) == 0 {
		return nil, fmt.Errorf("no node names were found")
	}
	return &logConfig, nil
}

func (s *ArtifactoryData) getVersion(ctx context.Context, serverId string) (string, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelTimeout()

	baseUrl, err := s.getUrl(serverId)
	if err != nil {
		return "", err
	}

	res, resBody, err := clientlayer.SendGet(timeoutCtx, serverId, artifactoryVersionEndPoint,constants.EmptyNodeId, baseUrl, nil)
	if err != nil {
		return "", err
	}

	err = errorHandle(res.StatusCode, resBody)
	if err != nil {
		return "", err
	}

	versionData := artifactoryVersionData{}
	err = json.Unmarshal(resBody, &versionData)
	if err != nil {
		return "", err
	}
	if versionData.Version == "" {
		return "", fmt.Errorf("could not retreive version information from Artifactory")
	}
	return strings.TrimSpace(versionData.Version), nil
}

func (s *ArtifactoryData) GetLogData(ctx context.Context, serverId string) (logData model.Data, err error) {
	if s.nodeId == "" {
		return logData, fmt.Errorf("node id must be set")
	}
	if s.logFileName == "" {
		return logData, fmt.Errorf("log file name must be set")
	}

	err = s.ArtifactoryValidations(ctx, serverId)
	if err != nil {
		return logData, err
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultLogRequestTimeout)
	defer cancelTimeout()

	var endpoint string
	endpoint = fmt.Sprintf("%s?file_size=%d&id=%s", artifactoryDataEndpoint, s.lastPageMarker, s.logFileName)
	baseUrl, err := s.getUrl(serverId)
	if err != nil {
		return logData, err
	}

	res, resBody, err := clientlayer.SendGet(timeoutCtx, serverId, endpoint, s.nodeId, baseUrl, nil)
	if err != nil {
		return logData, err
	}

	err = errorHandle(res.StatusCode, resBody)
	if err != nil {
		return logData, err
	}

	if err := json.Unmarshal(resBody, &logData); err != nil {
		return logData, err
	}

	return logData, nil
}

func (s *ArtifactoryData) getUrl(serverId string)(url string,_ error){
	confDetails, err := cliCommands.GetConfig(serverId, false)
	if err != nil {
		return "", err
	}
	url = confDetails.GetArtifactoryUrl()
	if url == "" {
		return "", fmt.Errorf("the Artifactory url was not found in the serverId : %s; verify that you are using the latest version of the JFrog CLI", serverId)
	}
	return url, nil
}

func (s *ArtifactoryData) ArtifactoryValidations(ctx context.Context, serverId string) error {
	if os.Getenv(constants.VersionCheckEnv) == "false" {
		return nil
	}

	currentVersion, err := s.getVersion(ctx, serverId)
	if err != nil {
		return err
	}

	return checkVersion(currentVersion, artifactoryMinVersionSupport, artifactoryProductName)
}

func (s *ArtifactoryData) SetNodeId(nodeId string) {
	s.nodeId = nodeId
}

func (s *ArtifactoryData) SetLogFileName(logFileName string) {
	s.logFileName = logFileName
}

func (s *ArtifactoryData) SetLogsRefreshRate(logsRefreshRate time.Duration) {
	s.logsRefreshRate = logsRefreshRate
}

func (s *ArtifactoryData) SetLastPageMarker(pageMarker int64) {
	s.lastPageMarker = pageMarker
}

func (s *ArtifactoryData) GetLastPageMarker() int64 {
	return s.lastPageMarker
}

func (s *ArtifactoryData) GetNodeId() string {
	return s.nodeId
}

func (s *ArtifactoryData) GetLogFileName() string {
	return s.logFileName
}

func (s *ArtifactoryData) GetLogsRefreshRate() time.Duration {
	return s.logsRefreshRate
}
