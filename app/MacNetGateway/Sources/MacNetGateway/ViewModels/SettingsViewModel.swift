import Foundation

@MainActor
final class SettingsViewModel: ObservableObject {
    @Published var selectedTheme: ThemeOption = .system
    @Published private(set) var loadState: LoadState = .idle
    @Published private(set) var saveState: LoadState = .idle

    private let apiClient: APIClient

    init(apiClient: APIClient = DaemonAPIClient()) {
        self.apiClient = apiClient
    }

    func load() async {
        loadState = .loading
        do {
            let settings = try await apiClient.fetchSettings()
            let themeValue = settings["ui.theme"] ?? ThemeOption.system.rawValue
            selectedTheme = ThemeOption(rawValue: themeValue) ?? .system
            loadState = .loaded
        } catch {
            loadState = .failed(error.localizedDescription)
        }
    }

    func save() async {
        saveState = .loading
        do {
            try await apiClient.patchSettings(["ui.theme": selectedTheme.rawValue])
            saveState = .loaded
        } catch {
            saveState = .failed(error.localizedDescription)
        }
    }
}
