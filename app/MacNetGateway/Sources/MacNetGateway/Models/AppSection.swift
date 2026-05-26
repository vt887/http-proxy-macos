import Foundation

enum AppSection: String, CaseIterable, Identifiable {
    case dashboard
    case liveActivity
    case devices
    case groups
    case dns
    case proxy
    case rules
    case quotas
    case speedLimits
    case blocklists
    case services
    case settings

    var id: String { rawValue }

    var title: String {
        switch self {
        case .dashboard: "Dashboard"
        case .liveActivity: "Live Activity"
        case .devices: "Devices"
        case .groups: "Groups"
        case .dns: "DNS"
        case .proxy: "Proxy"
        case .rules: "Rules"
        case .quotas: "Quotas"
        case .speedLimits: "Speed Limits"
        case .blocklists: "Blocklists"
        case .services: "Services"
        case .settings: "Settings"
        }
    }

    var systemImage: String {
        switch self {
        case .dashboard: "speedometer"
        case .liveActivity: "waveform.path.ecg"
        case .devices: "desktopcomputer"
        case .groups: "person.3"
        case .dns: "network"
        case .proxy: "arrow.triangle.branch"
        case .rules: "checklist"
        case .quotas: "gauge.with.dots.needle.33percent"
        case .speedLimits: "tortoise"
        case .blocklists: "hand.raised"
        case .services: "server.rack"
        case .settings: "gearshape"
        }
    }
}
