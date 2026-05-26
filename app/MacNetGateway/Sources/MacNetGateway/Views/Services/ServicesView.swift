import SwiftUI

struct ServicesView: View {
    @StateObject private var viewModel = ServicesViewModel()

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Services")
                .font(.largeTitle)
                .bold()
            if let errorMessage = viewModel.errorMessage {
                Text(errorMessage)
                    .foregroundStyle(.red)
            } else if viewModel.services.isEmpty {
                ProgressView("Loading services…")
            } else {
                Table(viewModel.services) {
                    TableColumn("Name", value: \.name)
                    TableColumn("Status", value: \.status)
                    TableColumn("Message", value: \.message)
                }
            }
        }
        .padding()
        .task {
            await viewModel.load()
        }
    }
}
