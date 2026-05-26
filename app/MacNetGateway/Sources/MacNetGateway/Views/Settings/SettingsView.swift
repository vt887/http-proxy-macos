import SwiftUI

struct SettingsView: View {
    var body: some View {
        Form {
            Section("General") {
                Text("PR-1 uses local mock data.")
            }
            Section("Security") {
                Text("Daemon API is intended for 127.0.0.1 only.")
            }
        }
        .padding()
        .frame(width: 480)
    }
}
