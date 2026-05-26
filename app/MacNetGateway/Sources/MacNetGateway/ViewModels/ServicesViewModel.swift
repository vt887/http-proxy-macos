import Foundation

@MainActor
final class ServicesViewModel: ObservableObject {
    @Published private(set) var services: [ServiceStatus] = []
    @Published private(set) var errorMessage: String?

    private let apiClient: APIClient

    init(apiClient: APIClient = MockAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        do {
            services = try await apiClient.fetchServices()
            errorMessage = nil
        } catch {
            errorMessage = "Failed to load services."
        }
    }
}
