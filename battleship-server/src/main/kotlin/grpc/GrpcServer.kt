package dev.spris.battleship.server.grpc

import dev.spris.battleship.server.config.GrpcConfig
import io.github.oshai.kotlinlogging.KotlinLogging
import io.grpc.protobuf.services.ProtoReflectionService
import io.grpc.ServerBuilder
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import org.springframework.stereotype.Component

private val logger = KotlinLogging.logger {}

@Component
class GrpcServer(
    private final val config: GrpcConfig,
    private final val grpcBattleshipServerService: GrpcBattleshipServerService,
) {
    private val server = ServerBuilder
        .forPort(config.server.port)
        .addService(grpcBattleshipServerService)
        .addService(ProtoReflectionService.newInstance())
        .build()

    @PostConstruct
    private fun init() {
        logger.info { "Starting grpc server" }
        server.start()
        logger.info { "Grpc server listening on port ${config.server.port}" }
    }

    @PreDestroy
    private fun destroy() {
        logger.info { "Shutting down grpc server" }
        server.shutdown()
        logger.info { "Grpc server shut down" }
    }
}
