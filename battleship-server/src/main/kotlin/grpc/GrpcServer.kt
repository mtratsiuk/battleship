package dev.spris.battleship.server.grpc

import dev.spris.battleship.server.config.GrpcConfig
import io.github.oshai.kotlinlogging.KotlinLogging
import io.grpc.*
import io.grpc.netty.shaded.io.grpc.netty.NettyServerBuilder
import io.grpc.protobuf.services.ProtoReflectionService
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import java.net.InetSocketAddress
import org.springframework.stereotype.Component

private val logger = KotlinLogging.logger {}

@Component
class GrpcServer(
    private final val config: GrpcConfig,
    private final val grpcBattleshipServerService: GrpcBattleshipServerService,
) {
    private val server =
        NettyServerBuilder.forAddress(InetSocketAddress(config.server.host, config.server.port))
            .addService(grpcBattleshipServerService)
            .addService(ProtoReflectionService.newInstance())
            .intercept(
                object : ServerInterceptor {
                    override fun <ReqT : Any, RespT : Any> interceptCall(
                        call: ServerCall<ReqT, RespT>,
                        headers: Metadata,
                        next: ServerCallHandler<ReqT, RespT>
                    ): ServerCall.Listener<ReqT> {
                        return next.startCall(ExceptionTranslatingServerCall(call), headers)
                    }
                }
            )
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

    private class ExceptionTranslatingServerCall<ReqT, RespT>(delegate: ServerCall<ReqT, RespT>) :
        ForwardingServerCall.SimpleForwardingServerCall<ReqT, RespT>(delegate) {
        override fun close(status: Status, trailers: Metadata) {
            if (status.isOk) {
                return super.close(status, trailers)
            }

            val cause = status.cause
            val newStatus = status.withDescription(cause?.message).withCause(cause)

            logger.error(cause) { "Error while handling gRPC request" }

            super.close(newStatus, trailers)
        }
    }
}
