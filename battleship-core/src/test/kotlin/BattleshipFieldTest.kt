import dev.spris.battleship.core.BattleshipField
import dev.spris.battleship.core.BattleshipPos
import dev.spris.battleship.core.BattleshipShipTile
import dev.spris.battleship.core.BattleshipType
import kotlin.test.*

typealias ShipDef = Pair<BattleshipType, List<BattleshipPos>>

class BattleshipFieldTest {
    @Test
    fun `should create valid battleship field`() {
        val ships = createValidBattleshipField()
        val field = BattleshipField.fromShips(*ships.toTypedArray())

        assertEquals(
            field.toString(), """
        PSDBC.....
        PSDBC.....
        .SDBC.....
        ...BC.....
        ....C.....
        ..........
        ..........
        ..........
        ..........
        ..........
        """.trimIndent()
        )
    }

    @Test
    fun `should throw when trying to create invalid battleship field`() {
        for (ships in getInvalidBattleshipFields()) {
            assertFails("Expected failure for ${ships.fold("") { r, c -> "$r\n$c" }}") {
                BattleshipField.fromShips(*ships.toTypedArray())
            }
        }
    }

    @Test
    fun `#strike pos should be added to hits set`() {
        val ships = createValidBattleshipField()
        val field = BattleshipField.fromShips(*ships.toTypedArray())

        field.strikeAt(BattleshipPos(0, 0))

        assertContains(field.hits, BattleshipPos(0, 0))
        assertEquals(field.hits.size, 1)
        assertEquals(field.misses.size, 0)
    }

    @Test
    fun `#strike pos should be added to misses set`() {
        val ships = createValidBattleshipField()
        val field = BattleshipField.fromShips(*ships.toTypedArray())

        field.strikeAt(BattleshipPos(9, 9))

        assertContains(field.misses, BattleshipPos(9, 9))
        assertEquals(field.hits.size, 0)
        assertEquals(field.misses.size, 1)
    }

    @Test
    fun `#hasAliveShips should return true`() {
        val ships = createValidBattleshipField()
        val field = BattleshipField.fromShips(*ships.toTypedArray())

        assertTrue { field.hasAliveShips() }

        field.strikeAt(BattleshipPos(0, 0))
        assertTrue { field.hasAliveShips() }
    }

    @Test
    fun `#hasAliveShips should return false when all ship tiles have been hit`() {
        val ships = createValidBattleshipField()
        val field = BattleshipField.fromShips(*ships.toTypedArray())

        for ((y, row) in field.field.withIndex()) {
            for ((x, tile) in row.withIndex()) {
                if (tile is BattleshipShipTile) {
                    field.strikeAt(BattleshipPos(x, y))
                }
            }
        }

        assertFalse { field.hasAliveShips() }
    }

    fun getInvalidBattleshipFields(): List<List<ShipDef>> {
        return listOf(
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
            ),
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
                createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(1, 1), Pair(1, 2)),
            ),
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
                createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(1, 1), Pair(1, 2)),
                createShip(BattleshipType.DESTROYER, Pair(2, 0), Pair(2, 1), Pair(2, 2)),
            ),
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
                createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(1, 1), Pair(1, 2)),
                createShip(BattleshipType.DESTROYER, Pair(2, 0), Pair(2, 1), Pair(2, 2)),
                createShip(BattleshipType.BATTLESHIP, Pair(3, 0), Pair(3, 1), Pair(3, 2), Pair(3, 3)),
            ),
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(9, 0), Pair(0, 1)),
                createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(1, 1), Pair(1, 2)),
                createShip(BattleshipType.DESTROYER, Pair(2, 0), Pair(2, 1), Pair(2, 2)),
                createShip(BattleshipType.BATTLESHIP, Pair(3, 0), Pair(3, 1), Pair(3, 2), Pair(3, 3)),
                createShip(BattleshipType.CARRIER, Pair(4, 0), Pair(4, 1), Pair(4, 2), Pair(4, 3), Pair(4, 4)),
            ),
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
                createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(1, 1), Pair(1, 2)),
                createShip(BattleshipType.DESTROYER, Pair(2, 0), Pair(2, 1), Pair(2, 2)),
                createShip(BattleshipType.BATTLESHIP, Pair(3, 0), Pair(3, 1), Pair(3, 2), Pair(3, 3)),
                createShip(
                    BattleshipType.CARRIER,
                    Pair(4, 0),
                    Pair(4, 1),
                    Pair(4, 2),
                    Pair(4, 3),
                    Pair(4, 4),
                    Pair(4, 5)
                ),
            ),
            listOf(
                createShip(BattleshipType.PATROL_BOAT, Pair(0, 0), Pair(0, 1)),
                createShip(BattleshipType.SUBMARINE, Pair(1, 0), Pair(2, 1), Pair(3, 2)),
                createShip(BattleshipType.DESTROYER, Pair(2, 0), Pair(2, 1), Pair(2, 2)),
                createShip(BattleshipType.BATTLESHIP, Pair(3, 0), Pair(3, 1), Pair(3, 2), Pair(3, 3)),
                createShip(BattleshipType.CARRIER, Pair(4, 0), Pair(4, 1), Pair(4, 2), Pair(4, 3), Pair(4, 4)),
            ),
        )
    }
}
