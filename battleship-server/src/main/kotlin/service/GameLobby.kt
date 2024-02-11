package dev.spris.battleship.server.service

import dev.spris.battleship.server.repository.PlayerRepository
import org.springframework.stereotype.Service

const val ROUNDS_COUNT = 3

@Service
class GameLobby(
    val playerRepository: PlayerRepository,
    val gameRunner: GameRunner,
) {
    suspend fun join(addr: String, name: String) {
        val players = playerRepository.findAll()
        val newPlayer = playerRepository.create(addr, name)

        for (player in players) {
            for (i in 0..<ROUNDS_COUNT) {
                val even = i % 2 == 0
                val firstPlayer = if (even) player else newPlayer
                val secondPlayer = if (even) newPlayer else player

                gameRunner.addGame(firstPlayer, secondPlayer)
            }
        }
    }
}
