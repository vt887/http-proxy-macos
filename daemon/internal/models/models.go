package models

type Health struct {
	Status string `json:"status"`
}

type Dashboard struct {
	ActiveDevices          int     `json:"active_devices"`
	ProxyRequestsPerMinute int     `json:"proxy_requests_per_minute"`
	DNSQueriesPerMinute    int     `json:"dns_queries_per_minute"`
	BlockedRequests        int     `json:"blocked_requests"`
	TrafficTodayMB         int     `json:"traffic_today_mb"`
	CacheHitRatio          float64 `json:"cache_hit_ratio"`
}

type ServiceStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Event struct {
	Time   string `json:"time"`
	Type   string `json:"type"`
	Target string `json:"target"`
	Action string `json:"action"`
}

type ProxySettings struct {
	ListenAddress       string `json:"listen_address"`
	CacheDirectory      string `json:"cache_directory"`
	GeneratedConfigPath string `json:"generated_config_path"`
	ConfigPreview       string `json:"config_preview"`
}

type DNSSettings struct {
	ListenAddress       string `json:"listen_address"`
	UpstreamMode        string `json:"upstream_mode"`
	Provider            string `json:"provider"`
	GeneratedConfigPath string `json:"generated_config_path"`
	ConfigPreview       string `json:"config_preview"`
}

type ActionResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
