import Foundation

@MainActor
final class ProxyViewModel: ObservableObject {
    @Published private(set) var status: ProxyStatus?
    @Published private(set) var settings: ProxySettings?
    @Published private(set) var loadState: LoadState = .idle
    @Published private(set) var actionState: LoadState = .idle
    @Published private(set) var actionMessage: String?

    private let apiClient: APIClient

    init(apiClient: APIClient = DaemonAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        loadState = .loading
        do {
            status = try await apiClient.fetchProxyStatus()
            settings = try await apiClient.fetchProxySettings()
            loadState = .loaded
        } catch {
            loadState = .failed(error.localizedDescription)
        }
    }

    func validateConfig() async {
        actionState = .loading
        do {
            let result = try await apiClient.validateProxyConfig()
            actionState = .loaded
            actionMessage = result.message
            await load()
        } catch {
            actionState = .failed(error.localizedDescription)
            actionMessage = nil
        }
    }

    func reloadConfig() async {
        actionState = .loading
        do {
            let result = try await apiClient.reloadProxyConfig()
            actionState = .loaded
            actionMessage = result.message
        } catch {
            actionState = .failed(error.localizedDescription)
            actionMessage = nil
        }
    }
}
