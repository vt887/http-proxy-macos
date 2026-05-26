package services

import "context"

type ServiceStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ProxyEvent struct{}
type DNSEvent struct{}
type HelperRequest struct{}
type HelperResponse struct{}

type ProxySettings struct {
	ListenAddress       string `json:"listen_address"`
	CacheDirectory      string `json:"cache_directory"`
	GeneratedConfigPath string `json:"generated_config_path"`
}

type SquidService interface {
	Status(ctx context.Context) (ServiceStatus, error)
	Settings(ctx context.Context) (ProxySettings, error)
	ConfigPreview(ctx context.Context) (string, error)
	RenderConfig(ctx context.Context) error
	ValidateConfig(ctx context.Context) error
	Reload(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	TailAccessLog(ctx context.Context) (<-chan ProxyEvent, error)
}

type DNSService interface {
	Status(ctx context.Context) (ServiceStatus, error)
	RenderConfig(ctx context.Context) error
	ValidateConfig(ctx context.Context) error
	Reload(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	TailQueryLog(ctx context.Context) (<-chan DNSEvent, error)
}

type LaunchdService interface {
	Install(ctx context.Context, plistPath string) error
	Uninstall(ctx context.Context, label string) error
	Start(ctx context.Context, label string) error
	Stop(ctx context.Context, label string) error
	Restart(ctx context.Context, label string) error
	Status(ctx context.Context, label string) (ServiceStatus, error)
}

type HelperClient interface {
	IsAvailable(ctx context.Context) bool
	Execute(ctx context.Context, req HelperRequest) (HelperResponse, error)
}
