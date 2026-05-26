import Foundation

struct LiveActivityEvent: Codable, Identifiable {
    let time: String
    let type: String
    let target: String
    let action: String

    var id: String { "\(time)-\(type)-\(target)" }
}
