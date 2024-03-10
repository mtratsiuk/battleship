package dev.spris.battleship.core

sealed interface BattleshipGameLogEntry

data class BattleshipGameLogActionEntry(val action: BattleshipAction) : BattleshipGameLogEntry

data class BattleshipGameLogErrorEntry(val error: Exception) : BattleshipGameLogEntry

class BattleshipGameLog {
    private val log = mutableListOf<BattleshipGameLogEntry>()

    val entries: List<BattleshipGameLogEntry>
        get() = log

    fun append(entry: BattleshipGameLogEntry) {
        log.addLast(entry)
    }

    override fun toString() = "GameLog[${log.joinToString("\n")}]"
}
