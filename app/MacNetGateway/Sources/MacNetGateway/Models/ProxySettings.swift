import Foundation

struct ProxySettings: Codable {
    let listenAddress: String
    let cacheDirectory: String
    let generatedConfigPath: String
    let configPreview: String
}
