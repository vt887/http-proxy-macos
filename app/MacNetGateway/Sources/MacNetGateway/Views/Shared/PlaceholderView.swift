import SwiftUI

struct PlaceholderView: View {
    let title: String

    var body: some View {
        VStack(spacing: 12) {
            Image(systemName: "hammer")
                .font(.largeTitle)
            Text(title)
                .font(.title2)
            Text("Scaffolded for PR-1. Functional implementation lands in future PRs.")
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}
