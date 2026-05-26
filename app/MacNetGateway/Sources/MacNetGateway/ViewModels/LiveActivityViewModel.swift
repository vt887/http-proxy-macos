import Foundation

@MainActor
final class LiveActivityViewModel: ObservableObject {
    @Published private(set) var events: [LiveActivityEvent] = []
    @Published private(set) var loadState: LoadState = .idle

    private let apiClient: APIClient

    init(apiClient: APIClient = DaemonAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        loadState = .loading
        do {
            events = try await apiClient.fetchLiveActivity()
            loadState = .loaded
        } catch {
            loadState = .failed(error.localizedDescription)
        }
    }
}
