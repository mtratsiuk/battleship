package dev.spris.battleship.server.grpc

import io.github.oshai.kotlinlogging.KotlinLogging
import io.grpc.protobuf.services.ProtoReflectionService
import io.grpc.ServerBuilder
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import org.springframework.stereotype.Component

private val logger = KotlinLogging.logger {}
private const val GRPC_PORT = 6969

@Component
class GrpcServer {
    private val server = ServerBuilder
        .forPort(GRPC_PORT)
        .addService(GrpcBattleshipServerService())
        .addService(ProtoReflectionService.newInstance())
        .build()

    @PostConstruct
    private fun init() {
        logger.info { "Starting grpc server" }
        server.start()
        logger.info { "Grpc server listening on port $GRPC_PORT" }
    }

    @PreDestroy
    private fun destroy() {
        logger.info { "Shutting down grpc server" }
        server.shutdown()
        logger.info { "Grpc server shut down" }
    }
}
