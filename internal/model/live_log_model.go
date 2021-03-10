package model

type Config struct {
	LogFileNames      []string `json:"logs,omitempty"`
	Nodes             []string `json:"nodes"`
	RefreshRateMillis int64    `json:"refresh_rate_millis,omitempty"`
}

type Data struct {
	Content    string `json:"log_content,omitempty"`
	PageMarker int64  `json:"file_size,omitempty"`
}

type ServiceNode struct {
	NodeId string `json:"node-id"`
}

type ConfigDisplayData struct {
	Logs      []string `json:"logs,omitempty"`
	Nodes     []string `json:"nodes"`
}
