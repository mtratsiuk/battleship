package dev.spris.battleship.server.service

import dev.spris.battleship.core.*
import dev.spris.battleship.server.repository.Player
import io.github.oshai.kotlinlogging.KotlinLogging
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.channels.Channel.Factory.UNLIMITED
import org.springframework.stereotype.Service

const val GAME_TURNS_LIMIT = 10_000

private val logger = KotlinLogging.logger {}

@Service
class GameRunner(
    val playerDriverFactory: PlayerDriverFactory,
    val idGenerator: IdGenerator,
) {
    private val scope = CoroutineScope(Dispatchers.IO)
    private val games = Channel<Pair<Player, Player>>(UNLIMITED)

    @PostConstruct
    private fun init() {
        scope.launch {
            for (game in games) {
                launch { runGame(game.first, game.second) }
            }
        }
    }

    @PreDestroy
    private fun destroy() {
        games.cancel()
        scope.cancel()
    }

    suspend fun addGame(
        player1: Player,
        player2: Player,
    ) {
        games.send(player1 to player2)
    }

    private suspend fun runGame(
        player1: Player,
        player2: Player,
    ) {
        val players =
            mapOf(
                player1.id to playerDriverFactory.create(player1),
                player2.id to playerDriverFactory.create(player2),
            )

        val game =
            BattleshipGame(
                BattleshipGameId(idGenerator.next()),
                player1.id,
                player2.id,
            )

        logger.info { "Game ${game.gameId} started: $player1 vs $player2" }

        for (turn in 0..GAME_TURNS_LIMIT) {
            when (val state = game.state) {
                is BattleshipStateAwaitingField -> {
                    val field = players[state.playerId]!!.requestField(game.gameId)
                    game.accept(BattleshipActionField(state.playerId, field))
                }
                is BattleshipStateAwaitingStrike -> {
                    val strike =
                        players[state.attackerId]!!.requestStrike(
                            gameId = game.gameId,
                            ownField = game.playerField(state.attackerId),
                            otherField = game.playerField(game.otherPlayerId(state.attackerId)),
                        )

                    game.accept(BattleshipActionStrike(state.attackerId, strike))
                }
                is BattleshipStateGameOver -> break
            }
        }

        // TODO: Save game log to game history

        if (game.state !is BattleshipStateGameOver) {
            throw IllegalStateException("Game is taking too long to finish, aborted")
        }

        logger.info { "Game ${game.gameId} finished: $player1 vs $player2" }
    }
}
