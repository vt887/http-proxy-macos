import Foundation

@MainActor
final class ServicesViewModel: ObservableObject {
    @Published private(set) var services: [ServiceStatus] = []
    @Published private(set) var loadState: LoadState = .idle

    private let apiClient: APIClient

    init(apiClient: APIClient = DaemonAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        loadState = .loading
        do {
            services = try await apiClient.fetchServices()
            loadState = .loaded
        } catch {
            loadState = .failed(error.localizedDescription)
        }
    }
}
