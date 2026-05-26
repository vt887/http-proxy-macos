import Foundation

@MainActor
final class DashboardViewModel: ObservableObject {
    @Published private(set) var metrics: DashboardMetrics?
    @Published private(set) var errorMessage: String?

    private let apiClient: APIClient

    init(apiClient: APIClient = MockAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        do {
            metrics = try await apiClient.fetchDashboard()
            errorMessage = nil
        } catch {
            errorMessage = "Failed to load dashboard metrics."
        }
    }
}
