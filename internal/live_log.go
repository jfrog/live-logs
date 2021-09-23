package livelog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cliCommands "github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/live-logs/internal/model"
	"github.com/jfrog/live-logs/internal/servicelayer"
	"github.com/jfrog/live-logs/internal/util"
	"io"
	"os"
	"time"
)

// Method initialised as a variable to improved unit test coverage
var getAllServiceIds = cliCommands.GetAllServerIds
var newServiceLayer = servicelayer.NewService

const (
	defaultLogsRefreshRate   = time.Second
)

type Data struct {
	productId       string
	serviceId       string
	serviceLayerClient servicelayer.ServiceLayer
	logsRefreshRate time.Duration
}

type LiveLogs interface {
	// The non interactive flow to display logs data.
	// The configured product id, server id, node id and log file name are used.
	// Any error during read or write is returned.
	LogNonInteractive(ctx context.Context, cliProductId, cliServerId, nodeId, logName string, isStreaming bool) error

	// Writes continuous or given single log data snapshots from the remote service into the passed io.Writer.
	// The configured product id, server id, node id and log file name are used.
	// Any error during read or write is returned.
	PrintLogs (ctx context.Context, nodeId, logName  string, isStreaming bool) error

	// Displays the list of available nodes and log files.
	ConfigNonInteractive(ctx context.Context, cliProductId, cliServerId string) error

	// A wrapper around service config api method call.
	GetConfigData (ctx context.Context, productId, serviceId string) (srvConfig *model.Config, err error)

	// Displays the list of available nodes and log files.
	DisplayConfig(ctx context.Context)  error

	// Sets the product id to use when querying the remote service for log data.
	SetProductId(productId string)
	GetProductId() (productId string)

	// Sets the service id to use when querying the remote service for log data.
	SetServiceId(serviceId string)
	GetServiceId() (serviceId string)

	// Sets the log refresh rate to use when querying the remote service for log data in tail mode.
	SetLogsRefreshRate(logsRefreshRate time.Duration)
	GetLogsRefreshRate() (logsRefreshRate time.Duration)

	// Sets and gets the a service layer.
	SetServiceLayer(productId string) error
	GetServiceLayer() servicelayer.ServiceLayer
}

func NewLiveLogs() LiveLogs {
	return &Data{
		logsRefreshRate: defaultLogsRefreshRate,
	}
}

func (s *Data) SetProductId(productId string) {
	s.productId = productId
}

func (s *Data) SetServiceLayer(productId string) error {
	var err error
	s.serviceLayerClient, err = newServiceLayer(productId)
	return err
}

func (s *Data) GetServiceLayer() servicelayer.ServiceLayer {
	return s.serviceLayerClient
}

func (s *Data) SetServiceId(serviceId string) {
	s.serviceId = serviceId
}

func (s *Data) GetLogsRefreshRate() (logsRefreshRate time.Duration) {
	return s.logsRefreshRate
}

func (s *Data) GetProductId()  string {
	return s.productId
}

func (s *Data) GetServiceId() string {
	return s.serviceId
}

func (s *Data) SetLogsRefreshRate(logsRefreshRate time.Duration) {
	s.logsRefreshRate = logsRefreshRate
}

func (s *Data) CatLog(ctx context.Context, output io.Writer) error {
	s.GetServiceLayer().SetLastPageMarker(0)
	logReader, err := s.doCatLog(ctx)
	if err != nil {
		return err
	}
	_, err = io.Copy(output, logReader)
	return err
}

func (s *Data) tailLog(ctx context.Context, output io.Writer) error {
	s.GetServiceLayer().SetLastPageMarker(0)
	curLogRefreshRate := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(curLogRefreshRate):
			if curLogRefreshRate == 0 {
				curLogRefreshRate = s.logsRefreshRate
			}
			var logReader io.Reader
			var err error

			logReader, err = s.doCatLog(ctx)

			if err != nil {
				return err
			}
			_, err = io.Copy(output, logReader)
			if err != nil {
				return err
			}
		}
	}
}

func (s *Data) doCatLog(ctx context.Context) (logReader io.Reader, err error) {
	if s.GetServiceLayer().GetNodeId() == "" {
		return nil, fmt.Errorf("node id must be set")
	}
	if s.GetServiceLayer().GetLogFileName() == "" {
		return nil, fmt.Errorf("log file name must be set")
	}
	logData := model.Data{}
	logData, err = s.GetServiceLayer().GetLogData(ctx,s.GetServiceId())
	if err != nil {
		return nil, err
	}
	s.GetServiceLayer().SetLastPageMarker(logData.PageMarker)
	logDataBuf := bytes.NewBufferString(logData.Content)
	return logDataBuf, nil
}

func (s *Data) LogNonInteractive(ctx context.Context, cliProductId, cliServerId, nodeId, logName string, isStreaming bool) error {
	productIds := util.FetchAllProductIds()
	err := util.ValidateArgument("product id", cliProductId, productIds)

	s.SetProductId(cliProductId)
	serverIds := getAllServiceIds()
	s.SetServiceId(cliServerId)
	err = util.ValidateArgument("server id", cliServerId, serverIds)
	if err != nil {
		return err
	}

	err = s.SetServiceLayer(s.GetProductId())
	if err != nil {
		return err
	}

	var logsRefreshRate time.Duration
	srvConfig, fetchErr := s.GetServiceLayer().GetConfig(ctx,s.GetServiceId())
	if fetchErr != nil {
		return fetchErr
	}

	logsRefreshRate = util.MillisToDuration(srvConfig.RefreshRateMillis)
	err = util.ValidateArgument("log name", logName, srvConfig.LogFileNames)

	if err != nil {
		return err
	}
	err = util.ValidateArgument("node id", nodeId, srvConfig.Nodes)
	if err != nil {
		return err
	}

	s.SetLogsRefreshRate(logsRefreshRate)
	return s.PrintLogs(ctx , nodeId, logName, isStreaming)
}

func (s *Data) GetConfigData (ctx context.Context, productId, serviceId string) (srvConfig *model.Config, err error) {
	s.SetServiceLayer(productId)
	if err != nil {
		return nil, err
	}

	srvConfig,err = s.GetServiceLayer().GetConfig(ctx,s.GetServiceId())
	return srvConfig, err
}

func (s *Data) PrintLogs (ctx context.Context, nodeId, logName  string, isStreaming bool) error {
	var err error
	if s.GetServiceLayer() == nil {
		err = s.SetServiceLayer(s.GetProductId())
		if err != nil {
			return err
		}
	}
	if isStreaming == true {
		s.SetLogsRefreshRate(s.GetLogsRefreshRate())
	}
	s.GetServiceLayer().SetLogFileName(logName)
	s.GetServiceLayer().SetNodeId(nodeId)

	if isStreaming {
		return s.tailLog(ctx, os.Stdout)
	}
	return s.CatLog(ctx, os.Stdout)
}

func (s *Data) DisplayConfig(ctx context.Context)  error {
	var err error

	err = s.SetServiceLayer(s.GetProductId())
	if err != nil {
		return err
	}

	srvConfig, fetchErr := s.GetConfigData(ctx, s.GetProductId(), s.GetServiceId())

	if fetchErr != nil {
		return fetchErr
	}

	var displayData model.ConfigDisplayData

	displayData.Logs = srvConfig.LogFileNames
	displayData.Nodes = srvConfig.Nodes

	data, err := json.MarshalIndent(displayData, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func (s *Data)  ConfigNonInteractive(ctx context.Context, cliProductId, cliServerId string) error {
	productIds := util.FetchAllProductIds()

	err := util.ValidateArgument("product id", cliProductId, productIds)
	if err != nil {
		return err
	}
	s.SetProductId(cliProductId)

	serverIds := getAllServiceIds()

	err = util.ValidateArgument("server id", cliServerId, serverIds)
	s.SetServiceId(cliServerId)

	return s.DisplayConfig(ctx)
}
