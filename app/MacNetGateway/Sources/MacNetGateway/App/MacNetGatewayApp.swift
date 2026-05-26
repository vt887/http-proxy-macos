import SwiftUI

@main
struct MacNetGatewayApp: App {
    @State private var selection: AppSection? = .dashboard

    var body: some Scene {
        WindowGroup("MacNet Gateway") {
            NavigationSplitView {
                List(AppSection.allCases, selection: $selection) { section in
                    Label(section.title, systemImage: section.systemImage)
                        .tag(section)
                }
                .navigationTitle("MacNet Gateway")
            } detail: {
                ContentView(selection: selection ?? .dashboard)
            }
        }
        Settings {
            SettingsView()
        }
    }
}
