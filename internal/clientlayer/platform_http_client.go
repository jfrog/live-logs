package clientlayer

import (
	"context"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	cliCommands "github.com/jfrog/jfrog-cli-core/common/commands"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/live-logs/internal/constants"
	"net/http"
)

func newPlatformHttpClient(cliServerId string) (*platformHttpClient, error) {
	platformDetails, err := cliCommands.GetConfig(cliServerId, false)
	if err != nil {
		return nil, err
	}
	platform, err := utils.CreateServiceManager(platformDetails, false)
	if err != nil {
		return nil, err
	}

	return &platformHttpClient{
		platform: platform,
	},nil
}

type platformHttpClient struct {
	platform artifactory.ArtifactoryServicesManager
}

func SendGet(_ context.Context, cliServerId, endpoint, nodeId, baseUrl string, extraHeaders map[string]string) (*http.Response, []byte, error) {

	platformClient, err := newPlatformHttpClient(cliServerId)
	if err != nil {
		return nil, nil, err
	}
	client := platformClient.platform.Client()

	platformDetails, err := cliCommands.GetConfig(cliServerId, false)
	artAuth, err := platformDetails.CreateArtAuthConfig()
	httpClientDetails := artAuth.CreateHttpClientDetails()

	if nodeId != constants.EmptyNodeId && nodeId != "" {
		httpClientDetails.Headers[constants.NodeIdHeader] = nodeId
	}

	for key, value := range extraHeaders {
		httpClientDetails.Headers[key] = value
	}

	res, resBody, _, err := client.SendGet(baseUrl+endpoint, true, &httpClientDetails)

	if err != nil {
		return nil, nil, err
	}
	return res, resBody, nil
}
