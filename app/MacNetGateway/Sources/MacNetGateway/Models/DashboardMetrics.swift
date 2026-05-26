import Foundation

struct DashboardMetrics: Codable {
    let activeDevices: Int
    let proxyRequestsPerMinute: Int
    let dnsQueriesPerMinute: Int
    let blockedRequests: Int
    let trafficTodayMB: Int
    let cacheHitRatio: Double
}
