package dev.spris.battleship.server.model

import dev.spris.battleship.core.BattleshipField
import dev.spris.battleship.core.BattleshipGameId
import dev.spris.battleship.core.BattleshipPos

interface PlayerDriver {
    suspend fun requestField(
        gameId: BattleshipGameId,
    ): BattleshipField

    suspend fun requestStrike(
        gameId: BattleshipGameId,
        ownField: BattleshipField,
        otherField: BattleshipField,
    ): BattleshipPos
}
