import SwiftUI

struct LiveActivityView: View {
    @StateObject private var viewModel = LiveActivityViewModel()

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Live Activity")
                .font(.largeTitle)
                .bold()
            switch viewModel.loadState {
            case .idle, .loading:
                ProgressView("Loading events…")
            case .loaded:
                Table(viewModel.events) {
                    TableColumn("Time", value: \.time)
                    TableColumn("Type", value: \.type)
                    TableColumn("Target", value: \.target)
                    TableColumn("Action", value: \.action)
                }
            case .failed(let message):
                VStack(alignment: .leading, spacing: 8) {
                    Text(message)
                        .foregroundStyle(.red)
                    Button("Retry") {
                        Task { await viewModel.load() }
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
