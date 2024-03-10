package dev.spris.battleship.server.repository

import dev.spris.battleship.core.BattleshipGameId
import dev.spris.battleship.core.BattleshipGameLog
import dev.spris.battleship.core.BattleshipPlayerId
import dev.spris.battleship.server.service.IdGenerator
import java.util.concurrent.ConcurrentHashMap
import org.springframework.stereotype.Repository

enum class GameState {
    IDLE,
    RUNNING,
    FINISHED,
}

data class Game(
    val id: BattleshipGameId,
    val player1Id: BattleshipPlayerId,
    val player2Id: BattleshipPlayerId,
    val state: GameState,
    val log: BattleshipGameLog,
)

interface GameRepository {
    suspend fun create(player1Id: BattleshipPlayerId, player2Id: BattleshipPlayerId): Game

    suspend fun update(game: Game): Game

    suspend fun findById(gameId: BattleshipGameId): Game?

    suspend fun findAll(): Collection<Game>
}

@Repository
class InMemoryGameRepository(
    private val idGenerator: IdGenerator,
) : GameRepository {
    private val games = ConcurrentHashMap<BattleshipGameId, Game>()

    override suspend fun create(
        player1Id: BattleshipPlayerId,
        player2Id: BattleshipPlayerId
    ): Game {
        val game =
            Game(
                id = BattleshipGameId(idGenerator.next()),
                player1Id = player1Id,
                player2Id = player2Id,
                state = GameState.IDLE,
                log = BattleshipGameLog(),
            )

        games[game.id] = game

        return game
    }

    override suspend fun update(game: Game): Game {
        require(games.containsKey(game.id)) { "Game ${game.id} doesn't exist" }

        games[game.id] = game

        return game
    }

    override suspend fun findById(gameId: BattleshipGameId): Game? {
        return games[gameId]
    }

    override suspend fun findAll(): Collection<Game> {
        return games.values
    }
}
