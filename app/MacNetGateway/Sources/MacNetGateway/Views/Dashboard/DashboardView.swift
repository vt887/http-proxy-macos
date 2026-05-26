import SwiftUI

struct DashboardView: View {
    @StateObject private var viewModel = DashboardViewModel()

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Dashboard")
                .font(.largeTitle)
                .bold()
            if let metrics = viewModel.metrics {
                LazyVGrid(columns: [GridItem(.flexible()), GridItem(.flexible())], spacing: 12) {
                    MetricCard(title: "Active Devices", value: "\(metrics.activeDevices)")
                    MetricCard(title: "Proxy RPM", value: "\(metrics.proxyRequestsPerMinute)")
                    MetricCard(title: "DNS QPM", value: "\(metrics.dnsQueriesPerMinute)")
                    MetricCard(title: "Blocked", value: "\(metrics.blockedRequests)")
                    MetricCard(title: "Traffic Today", value: "\(metrics.trafficTodayMB) MB")
                    MetricCard(title: "Cache Hit Ratio", value: "\(Int(metrics.cacheHitRatio * 100))%")
                }
            } else if let errorMessage = viewModel.errorMessage {
                Text(errorMessage)
                    .foregroundStyle(.red)
            } else {
                ProgressView("Loading metrics…")
            }
        }
        .padding()
        .task {
            await viewModel.load()
        }
    }
}
