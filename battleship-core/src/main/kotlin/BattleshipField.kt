package dev.spris.battleship.core

private const val BattleshipFieldSize = 10
private val BattleshipTypesCount = enumValues<BattleshipType>().size
private val BattleshipTilesToHitCount = enumValues<BattleshipType>().fold(0) { r, c -> r + c.size }

enum class BattleshipType(val size: Int) {
    PATROL_BOAT(2),
    SUBMARINE(3),
    DESTROYER(3),
    BATTLESHIP(4),
    CARRIER(5);

    companion object {
        fun fromConsoleView(view: String): BattleshipType {
            return when (view) {
                "P" -> PATROL_BOAT
                "S" -> SUBMARINE
                "D" -> DESTROYER
                "B" -> BATTLESHIP
                "C" -> CARRIER
                else -> throw IllegalArgumentException("Unexpected console view: $view")
            }
        }
    }
}

fun BattleshipType.toConsoleView(): String {
    return when (this) {
        BattleshipType.PATROL_BOAT -> "P"
        BattleshipType.SUBMARINE -> "S"
        BattleshipType.DESTROYER -> "D"
        BattleshipType.BATTLESHIP -> "B"
        BattleshipType.CARRIER -> "C"
    }
}

sealed interface BattleshipTile {
    companion object {
        fun fromConsoleView(view: String): BattleshipTile {
            return when (view) {
                "." -> BattleshipEmptyTile
                else -> BattleshipShipTile(BattleshipType.fromConsoleView(view))
            }
        }
    }
}

data object BattleshipEmptyTile : BattleshipTile

class BattleshipShipTile(val shipType: BattleshipType) : BattleshipTile

data class BattleshipPos(val x: Int, val y: Int) {
    override fun toString() = "Pos[$x,$y]"
}

class BattleshipField(
    val field: Array<Array<BattleshipTile>>,
    val hits: MutableSet<BattleshipPos>,
    val misses: MutableSet<BattleshipPos>,
) {
    companion object {
        fun fieldArrayFromString(str: String): Array<Array<BattleshipTile>> {
            return str.split("\n").map { line ->
                Array(BattleshipFieldSize) {
                    BattleshipTile.fromConsoleView(line[it].toString())
                }
            }.toTypedArray()
        }

        fun fromShips(vararg ships: Pair<BattleshipType, List<BattleshipPos>>): BattleshipField {
            val field =
                Array(BattleshipFieldSize) {
                    Array<BattleshipTile>(BattleshipFieldSize) { BattleshipEmptyTile }
                }

            for ((shipType, positions) in ships) {
                for ((x, y) in positions) {
                    field[y][x] = BattleshipShipTile(shipType)
                }
            }

            return BattleshipField(
                field,
                hits = mutableSetOf(),
                misses = mutableSetOf(),
            )
        }
    }

    init {
        require(field.size == BattleshipFieldSize) {
            "Expected field rows count to be ${BattleshipFieldSize}, got ${field.size}"
        }

        val ships = mutableMapOf<BattleshipType, List<BattleshipPos>>()

        for (y in 0..<BattleshipFieldSize) {
            require(field[y].size == BattleshipFieldSize) {
                "Expected field columns count to be ${BattleshipFieldSize}, got ${field.size}"
            }

            for (x in 0..<BattleshipFieldSize) {
                val tile = field[y][x]

                if (tile is BattleshipShipTile) {
                    ships.getOrPut(tile.shipType) { mutableListOf() }.addLast(BattleshipPos(x, y))
                }
            }
        }

        require(enumValues<BattleshipType>().all { ships.contains(it) }) {
            "Expected ${enumValues<BattleshipType>().asList()} to be present, got ${ships.keys}"
        }
        require(ships.size == BattleshipTypesCount) {
            "Expected $BattleshipTypesCount ships to be present"
        }

        for ((shipType, positions) in ships) {
            require(positions.size == shipType.size) {
                "Expected $shipType to have ${shipType.size} tiles, got ${positions.size}"
            }

            when {
                positions[0].x + 1 == positions[1].x -> {
                    require(
                        positions.dropLast(1).withIndex().all {
                            it.value.x + 1 == positions[it.index + 1].x
                        }
                    ) {
                        "Expected all $shipType's tiles to be sequential"
                    }

                    require(positions.all { it.y == positions[0].y }) {
                        "Expected all $shipType's tiles to be horizontal"
                    }
                }

                positions[0].y + 1 == positions[1].y -> {
                    require(
                        positions.dropLast(1).withIndex().all {
                            it.value.y + 1 == positions[it.index + 1].y
                        }
                    ) {
                        "Expected all $shipType's tiles to be sequential"
                    }

                    require(positions.all { it.x == positions[0].x }) {
                        "Expected all $shipType's tiles to be vertical"
                    }
                }

                else ->
                    require(false) {
                        "Expected $shipType tiles direction to be either horizontal or vertical"
                    }
            }
        }
    }

    fun hasAliveShips(): Boolean = hits.size < BattleshipTilesToHitCount

    fun strikeAt(pos: BattleshipPos) {
        require(pos.x in 0..<BattleshipFieldSize && pos.y in 0..<BattleshipFieldSize) {
            "Strike position is out of bounds: $pos"
        }

        when (field[pos.y][pos.x]) {
            is BattleshipShipTile -> hits.add(pos)
            is BattleshipEmptyTile -> misses.add(pos)
        }
    }

    override fun toString(): String {
        val sb = StringBuilder()

        for (y in 0..<BattleshipFieldSize) {
            for (x in 0..<BattleshipFieldSize) {
                val tile = field[y][x]

                val char =
                    when (tile) {
                        is BattleshipShipTile -> tile.shipType.toConsoleView()
                        is BattleshipEmptyTile -> "."
                    }

                sb.append(char)
            }

            sb.append("\n")
        }

        return sb.toString().trimEnd()
    }
}
