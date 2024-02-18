package dev.spris.battleship.server.service

import dev.spris.battleship.core.BattleshipField
import dev.spris.battleship.core.BattleshipGameId
import dev.spris.battleship.core.BattleshipPos
import dev.spris.battleship.core.BattleshipType
import dev.spris.battleship.server.repository.Player
import kotlin.random.Random
import kotlinx.coroutines.delay
import org.springframework.stereotype.Service

@Service
class PlayerDriverFactory {
    suspend fun create(player: Player): PlayerDriver {
        return InProcessRandomPlayerDriver(player)
    }
}

interface PlayerDriver {
    suspend fun requestField(gameId: BattleshipGameId): BattleshipField

    suspend fun requestStrike(
        gameId: BattleshipGameId,
        ownField: BattleshipField,
        otherField: BattleshipField,
    ): BattleshipPos
}

class InProcessRandomPlayerDriver(
    val player: Player,
) : PlayerDriver {
    override suspend fun requestField(gameId: BattleshipGameId): BattleshipField {
        delay(500)

        return BattleshipField.fromShips(
            BattleshipType.PATROL_BOAT to listOf(BattleshipPos(0, 0), BattleshipPos(0, 1)),
            BattleshipType.SUBMARINE to
                listOf(BattleshipPos(1, 0), BattleshipPos(1, 1), BattleshipPos(1, 2)),
            BattleshipType.DESTROYER to
                listOf(BattleshipPos(2, 0), BattleshipPos(2, 1), BattleshipPos(2, 2)),
            BattleshipType.BATTLESHIP to
                listOf(
                    BattleshipPos(3, 0),
                    BattleshipPos(3, 1),
                    BattleshipPos(3, 2),
                    BattleshipPos(3, 3),
                ),
            BattleshipType.CARRIER to
                listOf(
                    BattleshipPos(4, 0),
                    BattleshipPos(4, 1),
                    BattleshipPos(4, 2),
                    BattleshipPos(4, 3),
                    BattleshipPos(4, 4),
                ),
        )
    }

    override suspend fun requestStrike(
        gameId: BattleshipGameId,
        ownField: BattleshipField,
        otherField: BattleshipField,
    ): BattleshipPos {
        delay(50)

        return BattleshipPos(
            Random.nextInt(otherField.field.size),
            Random.nextInt(otherField.field.size)
        )
    }
}
