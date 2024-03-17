package dev.spris.battleship.server.config

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(prefix = "grpc")
data class GrpcConfig(
    val server: GrpcServer,
) {
    data class GrpcServer(
        val port: Int,
        val host: String,
    )
}
