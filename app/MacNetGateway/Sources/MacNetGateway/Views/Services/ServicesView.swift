import SwiftUI

struct ServicesView: View {
    @StateObject private var viewModel = ServicesViewModel()

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Services")
                .font(.largeTitle)
                .bold()
            switch viewModel.loadState {
            case .idle, .loading:
                ProgressView("Loading services…")
            case .loaded:
                Table(viewModel.services) {
                    TableColumn("Name", value: \.name)
                    TableColumn("Status", value: \.status)
                    TableColumn("Message", value: \.message)
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
