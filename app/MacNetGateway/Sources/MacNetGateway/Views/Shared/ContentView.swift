import SwiftUI

struct ContentView: View {
    let selection: AppSection

    var body: some View {
        switch selection {
        case .dashboard:
            DashboardView()
        case .liveActivity:
            LiveActivityView()
        case .services:
            ServicesView()
        case .settings:
            SettingsView()
        default:
            PlaceholderView(title: selection.title)
        }
    }
}
