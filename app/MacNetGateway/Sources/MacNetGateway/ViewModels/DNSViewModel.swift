import Foundation

@MainActor
final class DNSViewModel: ObservableObject {
    @Published private(set) var status: DNSStatus?
    @Published private(set) var settings: DNSSettings?
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
            status = try await apiClient.fetchDNSStatus()
            settings = try await apiClient.fetchDNSSettings()
            loadState = .loaded
        } catch {
            loadState = .failed(error.localizedDescription)
        }
    }

    func validateConfig() async {
        actionState = .loading
        do {
            let result = try await apiClient.validateDNSConfig()
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
            let result = try await apiClient.reloadDNSConfig()
            actionState = .loaded
            actionMessage = result.message
        } catch {
            actionState = .failed(error.localizedDescription)
            actionMessage = nil
        }
    }
}
