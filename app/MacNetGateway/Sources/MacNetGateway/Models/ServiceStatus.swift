import Foundation

struct ServiceStatus: Codable, Identifiable {
    let name: String
    let status: String
    let message: String

    var id: String { name }
}
