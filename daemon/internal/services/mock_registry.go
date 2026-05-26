package services

func NewMockRegistry() []ServiceStatus {
	return []ServiceStatus{
		{Name: "macnet-gatewayd", Status: "running", Message: "healthy"},
		{Name: "macnet-helper", Status: "stopped", Message: "not installed"},
		{Name: "squid", Status: "mock", Message: "integration pending"},
		{Name: "dns", Status: "mock", Message: "integration pending"},
	}
}
