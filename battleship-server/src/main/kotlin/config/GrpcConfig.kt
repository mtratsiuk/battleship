package dev.spris.battleship.server.config

import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Configuration

@ConfigurationProperties(prefix = "grpc")
data class GrpcConfig(
    val server: GrpcServer,
) {
    data class GrpcServer(
        val port: Int,
    )
}
