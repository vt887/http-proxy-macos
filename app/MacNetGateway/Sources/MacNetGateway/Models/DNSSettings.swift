import Foundation

struct DNSSettings: Codable {
    let listenAddress: String
    let upstreamMode: String
    let provider: String
    let generatedConfigPath: String
    let configPreview: String
}

struct DNSActionResult: Codable {
    let status: String
    let message: String
}
