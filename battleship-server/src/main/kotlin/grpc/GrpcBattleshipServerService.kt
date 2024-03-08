package dev.spris.battleship.server.grpc

import dev.spris.battleship.proto.v1.BattleshipServerServiceGrpcKt.BattleshipServerServiceCoroutineImplBase
import dev.spris.battleship.proto.v1.JoinLobbyRequest
import dev.spris.battleship.proto.v1.JoinLobbyResponse
import io.github.oshai.kotlinlogging.KotlinLogging

private val logger = KotlinLogging.logger {}

class GrpcBattleshipServerService : BattleshipServerServiceCoroutineImplBase() {
    override suspend fun joinLobby(request: JoinLobbyRequest): JoinLobbyResponse {
        logger.info { "joinLobby: $request" }

        return super.joinLobby(request)
    }
}
