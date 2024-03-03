package dev.spris.battleship.core

sealed interface BattleshipGameLogEntry

data class BattleshipGameLogActionEntry(val action: BattleshipAction) : BattleshipGameLogEntry

data class BattleshipGameLogErrorEntry(val error: Exception) : BattleshipGameLogEntry

class BattleshipGameLog {
    private val entries = mutableListOf<BattleshipGameLogEntry>()

    fun read() = listOf(entries)

    fun append(entry: BattleshipGameLogEntry) {
        entries.addLast(entry)
    }

    override fun toString() = "GameLog[${entries.joinToString("\n")}]"
}
