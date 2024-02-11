package dev.spris.battleship.server.repository

import dev.spris.battleship.core.BattleshipPlayerId
import dev.spris.battleship.server.model.Player
import org.springframework.stereotype.Repository
import java.util.UUID
import java.util.concurrent.ConcurrentHashMap

interface PlayerRepository {
    suspend fun create(addr: String, name: String): Player

    suspend fun findAll(): List<Player>

    suspend fun findById(id: BattleshipPlayerId): Player?
}

@Repository
class InMemoryPlayerRepository : PlayerRepository {
    private val players = ConcurrentHashMap<BattleshipPlayerId, Player>()

    override suspend fun create(addr: String, name: String): Player {
        require(!players.any { it.value.name == name }) { "Name $name is already taken" }

        val player = Player(
            id = BattleshipPlayerId(UUID.randomUUID().toString()),
            addr=addr,
            name=name,
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
