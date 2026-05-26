import Foundation

protocol APIClient {
    @MainActor func fetchDashboard() async throws -> DashboardMetrics
    @MainActor func fetchServices() async throws -> [ServiceStatus]
    @MainActor func fetchLiveActivity() async throws -> [LiveActivityEvent]
    @MainActor func fetchProxyStatus() async throws -> ProxyStatus
    @MainActor func fetchProxySettings() async throws -> ProxySettings
    @MainActor func validateProxyConfig() async throws -> ProxyActionResult
    @MainActor func reloadProxyConfig() async throws -> ProxyActionResult
    @MainActor func fetchDNSStatus() async throws -> DNSStatus
    @MainActor func fetchDNSSettings() async throws -> DNSSettings
    @MainActor func validateDNSConfig() async throws -> DNSActionResult
    @MainActor func reloadDNSConfig() async throws -> DNSActionResult
    @MainActor func fetchSettings() async throws -> [String: String]
    @MainActor func patchSettings(_ payload: [String: String]) async throws
}

enum APIClientError: Error {
    case invalidURL
    case invalidResponse(statusCode: Int)
    case decodeFailed
    case transport(Error)
}

extension APIClientError: LocalizedError {
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            "Invalid daemon API URL."
        case .invalidResponse(let statusCode):
            "Daemon API returned status \(statusCode)."
        case .decodeFailed:
            "Failed to decode daemon response."
        case .transport(let error):
            "Daemon API request failed: \(error.localizedDescription)"
        }
    }
}

struct DaemonAPIClient: APIClient {
    private let baseURL: URL
    private let session: URLSession
    private let decoder: JSONDecoder

    init(
        baseURL: URL = URL(string: "http://127.0.0.1:18080")!,
        session: URLSession = .shared
    ) {
        self.baseURL = baseURL
        self.session = session
        self.decoder = JSONDecoder()
        self.decoder.keyDecodingStrategy = .convertFromSnakeCase
    }

    func fetchDashboard() async throws -> DashboardMetrics {
        try await request(path: "/api/dashboard", method: "GET", body: Optional<[String: String]>.none)
    }

    func fetchServices() async throws -> [ServiceStatus] {
        try await request(path: "/api/services", method: "GET", body: Optional<[String: String]>.none)
    }

    func fetchLiveActivity() async throws -> [LiveActivityEvent] {
        try await request(path: "/api/live-activity", method: "GET", body: Optional<[String: String]>.none)
    }

    func fetchProxyStatus() async throws -> ProxyStatus {
        try await request(path: "/api/proxy/status", method: "GET", body: Optional<[String: String]>.none)
    }

    func fetchProxySettings() async throws -> ProxySettings {
        try await request(path: "/api/proxy/settings", method: "GET", body: Optional<[String: String]>.none)
    }

    func validateProxyConfig() async throws -> ProxyActionResult {
        try await request(path: "/api/proxy/validate", method: "POST", body: Optional<[String: String]>.none)
    }

    func reloadProxyConfig() async throws -> ProxyActionResult {
        try await request(path: "/api/proxy/reload", method: "POST", body: Optional<[String: String]>.none)
    }

    func fetchDNSStatus() async throws -> DNSStatus {
        try await request(path: "/api/dns/status", method: "GET", body: Optional<[String: String]>.none)
    }

    func fetchDNSSettings() async throws -> DNSSettings {
        try await request(path: "/api/dns/settings", method: "GET", body: Optional<[String: String]>.none)
    }

    func validateDNSConfig() async throws -> DNSActionResult {
        try await request(path: "/api/dns/validate", method: "POST", body: Optional<[String: String]>.none)
    }

    func reloadDNSConfig() async throws -> DNSActionResult {
        try await request(path: "/api/dns/reload", method: "POST", body: Optional<[String: String]>.none)
    }

    func fetchSettings() async throws -> [String: String] {
        try await request(path: "/api/settings", method: "GET", body: Optional<[String: String]>.none)
    }

    func patchSettings(_ payload: [String: String]) async throws {
        let _: [String: String] = try await request(path: "/api/settings", method: "PATCH", body: payload)
    }

    private func request<T: Decodable, Body: Encodable>(
        path: String,
        method: String,
        body: Body?
    ) async throws -> T {
        guard let url = URL(string: path, relativeTo: baseURL) else {
            throw APIClientError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        request.timeoutInterval = 5
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let body {
            request.httpBody = try JSONEncoder().encode(body)
        }

        do {
            let (data, response) = try await session.data(for: request)
            guard let http = response as? HTTPURLResponse else {
                throw APIClientError.invalidResponse(statusCode: -1)
            }
            guard (200...299).contains(http.statusCode) else {
                throw APIClientError.invalidResponse(statusCode: http.statusCode)
            }
            do {
                return try decoder.decode(T.self, from: data)
            } catch {
                throw APIClientError.decodeFailed
            }
        } catch let apiError as APIClientError {
            throw apiError
        } catch {
            throw APIClientError.transport(error)
        }
    }
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

    func fetchSettings() async throws -> [String: String] {
        ["ui.theme": "system"]
    }

    func fetchLiveActivity() async throws -> [LiveActivityEvent] {
        [
            LiveActivityEvent(time: "2026-01-01T12:00:00Z", type: "SERVICE_STARTED", target: "macnet-gatewayd", action: "allowed"),
        ]
    }

    func fetchProxyStatus() async throws -> ProxyStatus {
        ProxyStatus(name: "squid", status: "mock", message: "Mock proxy status")
    }

    func fetchProxySettings() async throws -> ProxySettings {
        ProxySettings(
            listenAddress: "127.0.0.1:3128",
            cacheDirectory: "/Library/Caches/MacNetGateway/squid",
            generatedConfigPath: "~/.macnet-gateway-dev/generated/squid",
            configPreview: "http_port 127.0.0.1:3128\ncache_dir ufs /Library/Caches/MacNetGateway/squid 4096 16 256\n"
        )
    }

    func validateProxyConfig() async throws -> ProxyActionResult {
        ProxyActionResult(status: "valid", message: "Mock validation succeeded")
    }

    func reloadProxyConfig() async throws -> ProxyActionResult {
        ProxyActionResult(status: "reloaded", message: "Mock reload accepted")
    }

    func fetchDNSStatus() async throws -> DNSStatus {
        DNSStatus(name: "dns", status: "mock", message: "Mock DNS status")
    }

    func fetchDNSSettings() async throws -> DNSSettings {
        DNSSettings(
            listenAddress: "127.0.0.1:53",
            upstreamMode: "doh",
            provider: "cloudflare",
            generatedConfigPath: "~/.macnet-gateway-dev/generated/dns",
            configPreview: "listen-address=127.0.0.1\nport=53\nserver=1.1.1.1\ncache-size=1000\n"
        )
    }

    func validateDNSConfig() async throws -> DNSActionResult {
        DNSActionResult(status: "valid", message: "Mock DNS validation succeeded")
    }

    func reloadDNSConfig() async throws -> DNSActionResult {
        DNSActionResult(status: "reloaded", message: "Mock DNS reload accepted")
    }

    func patchSettings(_ payload: [String: String]) async throws {
        _ = payload
    }
}
