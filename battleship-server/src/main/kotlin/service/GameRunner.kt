package dev.spris.battleship.server.service

import dev.spris.battleship.server.model.Player
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.channels.Channel.Factory.UNLIMITED
import org.springframework.stereotype.Service

@Service
class GameRunner() {
    private val scope = CoroutineScope(Dispatchers.Default)
    private val games = Channel<Pair<Player, Player>>(UNLIMITED)

    @PostConstruct
    private fun init() {
        scope.launch {
            for (game in games) {
                runGame(game.first, game.second)
            }
        }

    }

    @PreDestroy
    private fun destroy() {
        games.cancel()
        scope.cancel()
    }

    suspend fun addGame(player1: Player, player2: Player) {
        games.send(Pair(player1, player2))
    }

    private suspend fun runGame(player1: Player, player2: Player) {
        println("Starting game $player1 vs $player2")
    }
}
