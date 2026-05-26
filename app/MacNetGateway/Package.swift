// swift-tools-version: 6.3

import PackageDescription

let package = Package(
    name: "MacNetGateway",
    platforms: [.macOS(.v13)],
    products: [
        .executable(name: "MacNetGateway", targets: ["MacNetGateway"]),
    ],
    targets: [
        .executableTarget(
            name: "MacNetGateway"
        ),
        .testTarget(
            name: "MacNetGatewayTests",
            dependencies: ["MacNetGateway"]
        ),
    ],
    swiftLanguageModes: [.v6]
)
