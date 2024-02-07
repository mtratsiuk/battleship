package dev.spris.battleship.cli

import dev.spris.battleship.core.*

fun main() {
    val field = BattleshipField.fromShips(
        Pair(
            BattleshipType.PATROL_BOAT, listOf(
                BattleshipPos(0, 0),
                BattleshipPos(0, 1),
            )
        ),
        Pair(
            BattleshipType.SUBMARINE, listOf(
                BattleshipPos(1, 0),
                BattleshipPos(1, 1),
                BattleshipPos(1, 2),
            )
        ),
        Pair(
            BattleshipType.DESTROYER, listOf(
                BattleshipPos(2, 0),
                BattleshipPos(2, 1),
                BattleshipPos(2, 2),
            )
        ),
        Pair(
            BattleshipType.BATTLESHIP, listOf(
                BattleshipPos(3, 0),
                BattleshipPos(3, 1),
                BattleshipPos(3, 2),
                BattleshipPos(3, 3),
            )
        ),
        Pair(
            BattleshipType.CARRIER, listOf(
                BattleshipPos(4, 0),
                BattleshipPos(4, 1),
                BattleshipPos(4, 2),
                BattleshipPos(4, 3),
                BattleshipPos(4, 4),
            )
        ),
    )

    println(field)
}
