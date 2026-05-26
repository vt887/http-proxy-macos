import Foundation

@MainActor
protocol APIClient {
    func fetchDashboard() async throws -> DashboardMetrics
    func fetchServices() async throws -> [ServiceStatus]
}

enum APIClientError: Error {
    case invalidResponse
}

struct MockAPIClient: APIClient {
    func fetchDashboard() async throws -> DashboardMetrics {
        DashboardMetrics(
            activeDevices: 8,
            proxyRequestsPerMinute: 143,
            dnsQueriesPerMinute: 322,
            blockedRequests: 27,
            trafficTodayMB: 1248,
            cacheHitRatio: 0.63
        )
    }

    func fetchServices() async throws -> [ServiceStatus] {
        [
            ServiceStatus(name: "macnet-gatewayd", status: "running", message: "Local API healthy"),
            ServiceStatus(name: "macnet-helper", status: "stopped", message: "Not installed"),
            ServiceStatus(name: "squid", status: "mock", message: "Integration not enabled in PR-1"),
            ServiceStatus(name: "dns", status: "mock", message: "Integration not enabled in PR-1"),
        ]
    }
}
