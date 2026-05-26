import SwiftUI

struct SettingsView: View {
    @StateObject private var viewModel = SettingsViewModel()

    var body: some View {
        Form {
            Section("General") {
                Picker("Appearance", selection: $viewModel.selectedTheme) {
                    ForEach(ThemeOption.allCases) { option in
                        Text(option.title).tag(option)
                    }
                }
                Button("Save Settings") {
                    Task { await viewModel.save() }
                }
                .disabled(viewModel.saveState == .loading)
            }
            Section("Security") {
                Text("Daemon API is intended for 127.0.0.1 only.")
            }
            switch viewModel.loadState {
            case .failed(let message):
                Section("Load Error") {
                    Text(message).foregroundStyle(.red)
                    Button("Retry") {
                        Task { await viewModel.load() }
                    }
                }
            default:
                EmptyView()
            }
            switch viewModel.saveState {
            case .loading:
                Section("Status") {
                    ProgressView("Saving…")
                }
            case .loaded:
                Section("Status") {
                    Text("Saved")
                        .foregroundStyle(.green)
                }
            case .failed(let message):
                Section("Status") {
                    Text(message)
                        .foregroundStyle(.red)
                }
            default:
                EmptyView()
            }
        }
        .padding()
        .frame(width: 480)
        .task {
            await viewModel.load()
        }
    }
}
