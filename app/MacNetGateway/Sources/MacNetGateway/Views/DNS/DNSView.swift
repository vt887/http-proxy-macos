import SwiftUI

struct DNSView: View {
    @StateObject private var viewModel = DNSViewModel()

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("DNS")
                .font(.largeTitle)
                .bold()

            switch viewModel.loadState {
            case .idle, .loading:
                ProgressView("Loading DNS settings…")
            case .failed(let message):
                VStack(alignment: .leading, spacing: 8) {
                    Text(message)
                        .foregroundStyle(.red)
                    Button("Retry") {
                        Task { await viewModel.load() }
                    }
                }
            case .loaded:
                if let status = viewModel.status, let settings = viewModel.settings {
                    GroupBox("Status") {
                        VStack(alignment: .leading, spacing: 6) {
                            Text("Service: \(status.name)")
                            Text("State: \(status.status)")
                            Text(status.message)
                                .foregroundStyle(.secondary)
                        }
                        .frame(maxWidth: .infinity, alignment: .leading)
                    }

                    GroupBox("Configuration") {
                        VStack(alignment: .leading, spacing: 8) {
                            Text("Listen: \(settings.listenAddress)")
                            Text("Upstream Mode: \(settings.upstreamMode)")
                            Text("Provider: \(settings.provider)")
                            Text("Generated: \(settings.generatedConfigPath)")
                            ScrollView {
                                Text(settings.configPreview)
                                    .font(.system(.body, design: .monospaced))
                                    .frame(maxWidth: .infinity, alignment: .leading)
                            }
                            .frame(minHeight: 160)
                        }
                        .frame(maxWidth: .infinity, alignment: .leading)
                    }

                    HStack(spacing: 12) {
                        Button("Validate Config") {
                            Task { await viewModel.validateConfig() }
                        }
                        Button("Reload") {
                            Task { await viewModel.reloadConfig() }
                        }
                        .buttonStyle(.borderedProminent)
                    }
                    .disabled(viewModel.actionState == .loading)

                    switch viewModel.actionState {
                    case .loading:
                        ProgressView("Running action…")
                    case .failed(let message):
                        Text(message).foregroundStyle(.red)
                    case .loaded:
                        if let actionMessage = viewModel.actionMessage {
                            Text(actionMessage).foregroundStyle(.green)
                        }
                    default:
                        EmptyView()
                    }
                }
            }
        }
        .padding()
        .task {
            await viewModel.load()
        }
    }
}
