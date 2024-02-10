import dev.spris.battleship.core.BattleshipPos
import dev.spris.battleship.core.BattleshipType

fun createShip(shipType: BattleshipType, vararg coords: Pair<Int, Int>): ShipDef {
    return Pair(shipType, coords.map { BattleshipPos(it.first, it.second) })
}

fun createValidBattleshipField(): List<ShipDef> {
    return listOf(
        createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
        createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(1, 1), Pair(1, 2)),
        createShip(BattleshipType.DESTROYER, Pair(2, 0), Pair(2, 1), Pair(2, 2)),
        createShip(BattleshipType.BATTLESHIP, Pair(3, 0), Pair(3, 1), Pair(3, 2), Pair(3, 3)),
        createShip(BattleshipType.CARRIER, Pair(4, 0), Pair(4, 1), Pair(4, 2), Pair(4, 3), Pair(4, 4)),
    )
}
