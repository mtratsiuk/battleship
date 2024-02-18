package dev.spris.battleship.server.repository

import dev.spris.battleship.core.BattleshipPlayerId
import dev.spris.battleship.server.service.IdGenerator
import java.util.concurrent.ConcurrentHashMap
import org.springframework.stereotype.Repository

data class Player(
    val id: BattleshipPlayerId,
    val addr: String,
    val name: String,
)

interface PlayerRepository {
    suspend fun create(
        addr: String,
        name: String,
    ): Player

    suspend fun findAll(): List<Player>

    suspend fun findById(id: BattleshipPlayerId): Player?
}

@Repository
class InMemoryPlayerRepository(
    private val idGenerator: IdGenerator,
) : PlayerRepository {
    private val players = ConcurrentHashMap<BattleshipPlayerId, Player>()

    override suspend fun create(
        addr: String,
        name: String,
    ): Player {
        val player =
            Player(
                id = BattleshipPlayerId(idGenerator.next()),
                addr = addr,
                name = name,
            )

        players[player.id] = player
        return player
    }

    override suspend fun findAll(): List<Player> {
        return players.values.toList()
    }

    override suspend fun findById(id: BattleshipPlayerId): Player? {
        return players[id]
    }
}
