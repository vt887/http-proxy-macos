import Foundation

@MainActor
final class DashboardViewModel: ObservableObject {
    @Published private(set) var metrics: DashboardMetrics?
    @Published private(set) var loadState: LoadState = .idle

    private let apiClient: APIClient

    init(apiClient: APIClient = DaemonAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        loadState = .loading
        do {
            metrics = try await apiClient.fetchDashboard()
            loadState = .loaded
        } catch {
            loadState = .failed(error.localizedDescription)
        }
    }
}
