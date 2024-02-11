package dev.spris.battleship.server.model

import dev.spris.battleship.core.BattleshipPlayerId

data class Player(
    val id: BattleshipPlayerId,
    val addr: String,
    val name: String,
)
