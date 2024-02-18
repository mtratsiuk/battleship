package dev.spris.battleship.server.service

import dev.spris.battleship.server.repository.Player
import dev.spris.battleship.server.repository.PlayerRepository
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import org.springframework.stereotype.Service

const val ROUNDS_COUNT = 3

@Service
class GameLobby(
    private val playerRepository: PlayerRepository,
    private val gameRunner: GameRunner,
) {
    private val joinMutex = Mutex()

    suspend fun join(
        addr: String,
        name: String,
    ) {
        var players: List<Player>
        var newPlayer: Player

        joinMutex.withLock {
            players = playerRepository.findAll()

            require(!players.any { it.name == name }) { "Player name $name is already taken" }

            newPlayer = playerRepository.create(addr, name)
        }

        for (player in players) {
            for (i in 0 ..< ROUNDS_COUNT) {
                val even = i % 2 == 0
                val firstPlayer = if (even) player else newPlayer
                val secondPlayer = if (even) newPlayer else player

                gameRunner.addGame(firstPlayer, secondPlayer)
            }
        }
    }
}
