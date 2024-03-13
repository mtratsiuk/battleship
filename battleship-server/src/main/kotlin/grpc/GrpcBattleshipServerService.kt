package dev.spris.battleship.server.grpc

import dev.spris.battleship.core.BattleshipActionGameOver
import dev.spris.battleship.core.BattleshipGameId
import dev.spris.battleship.core.BattleshipGameLogActionEntry
import dev.spris.battleship.proto.server.v1.*
import dev.spris.battleship.proto.server.v1.BattleshipServerServiceGrpcKt.BattleshipServerServiceCoroutineImplBase
import dev.spris.battleship.server.repository.Game
import dev.spris.battleship.server.repository.GameRepository
import dev.spris.battleship.server.repository.GameState
import dev.spris.battleship.server.repository.PlayerRepository
import dev.spris.battleship.server.service.GameLobby
import io.github.oshai.kotlinlogging.KotlinLogging
import org.springframework.stereotype.Component
import java.util.concurrent.atomic.AtomicInteger

private val logger = KotlinLogging.logger {}

@Component
class GrpcBattleshipServerService(
    private val gameLobby: GameLobby,
    private val gameRepository: GameRepository,
    private val playerRepository: PlayerRepository,
) : BattleshipServerServiceCoroutineImplBase() {
    private var randomBotCounter = AtomicInteger(0)

    override suspend fun joinLobby(request: JoinLobbyRequest): JoinLobbyResponse {
        logger.info { "joinLobby: $request" }

        gameLobby.join(request.addr, request.name)

        return joinLobbyResponse {}
    }

    override suspend fun getGames(request: GetGamesRequest): GetGamesResponse {
        logger.info { "getGames: $request" }

        val players = playerRepository.findAll().associateBy { it.id }
        val gs = gameRepository.findAll().map { g ->
            getGamesResponseEntry {
                id = g.id.id
                state = g.state.toProto()
                player1 = players[g.player1Id]!!.toProto()
                player2 = players[g.player2Id]!!.toProto()

                g.maybeWinnerId()?.let {
                    winnerId = it
                }
            }
        }

        return getGamesResponse {
            games.addAll(gs)
        }
    }

    override suspend fun getGame(request: GetGameRequest): GetGameResponse {
        logger.info { "getGame: $request" }

        val game = gameRepository.findById(BattleshipGameId(request.id))

        require(game != null) { "Game ${request.id} not found" }

        val player1 = playerRepository.findById(game.player1Id)
        val player2 = playerRepository.findById(game.player2Id)

        return getGameResponse {
            this.game = gameProto {
                id = game.id.id
                this.player1 = player1!!.toProto()
                this.player2 = player2!!.toProto()
                state = game.state.toProto()
                log.addAll(game.log.entries.map { it.toProto() })
            }
        }
    }

    override suspend fun addRandomBot(request: AddRandomBotRequest): AddRandomBotResponse {
        logger.info { "addRandomBot: $request" }

        val ps = listOf("Mighty", "Funny", "Clever", "Handsome", "Tiny", "Pink")
        val ss = listOf("Cat", "Puppy", "Parrot", "Pony", "Bear", "Mouse", "Snake")
        val name = "${ps.random()} ${ss.random()} ${randomBotCounter.incrementAndGet()}"
        val addr = "inprocess://$name"

        gameLobby.join(addr, name)

        return addRandomBotResponse { }
    }
}

fun Game.maybeWinnerId(): String? {
    if (this.state != GameState.FINISHED || this.log.entries.isEmpty()) {
        return null
    }

    val lastLog = this.log.entries.last()

    if (lastLog is BattleshipGameLogActionEntry && lastLog.action is BattleshipActionGameOver) {
        return (lastLog.action as BattleshipActionGameOver).winnerId.id
    }

    return null
}
