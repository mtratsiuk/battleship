package dev.spris.battleship.core

import java.util.*
import kotlin.reflect.typeOf

@JvmInline
value class BattleshipGameId(val id: String) {
    override fun toString() = "GameId[$id]"
}

@JvmInline
value class BattleshipPlayerId(val id: String) {
    override fun toString() = "PlayerId[$id]"
}

sealed interface BattleshipAction

data class BattleshipActionField(
    val playerId: BattleshipPlayerId,
    val field: BattleshipField,
) : BattleshipAction

data class BattleshipActionStrike(
    val attackerId: BattleshipPlayerId,
    val position: BattleshipPos,
) : BattleshipAction

data class BattleshipActionGameOver(
    val winnerId: BattleshipPlayerId,
) : BattleshipAction

sealed interface BattleshipState

data class BattleshipStateAwaitingField(val playerId: BattleshipPlayerId) : BattleshipState

data class BattleshipStateAwaitingStrike(val attackerId: BattleshipPlayerId) : BattleshipState

data class BattleshipStateGameOver(val winnerId: BattleshipPlayerId) : BattleshipState

class BattleshipGame(
    val gameId: BattleshipGameId,
    val player1Id: BattleshipPlayerId,
    val player2Id: BattleshipPlayerId,
) {
    constructor(
        player1Id: BattleshipPlayerId,
        player2Id: BattleshipPlayerId
    ) : this(BattleshipGameId(UUID.randomUUID().toString()), player1Id, player2Id)

    val log = BattleshipGameLog()

    lateinit var player1Field: BattleshipField
    lateinit var player2Field: BattleshipField

    var state: BattleshipState = BattleshipStateAwaitingField(player1Id)

    fun accept(action: BattleshipAction) {
        log.append(BattleshipGameLogActionEntry(action))

        when (action) {
            is BattleshipActionField -> {
                assertState<BattleshipStateAwaitingField> { state ->
                    require(state.playerId == action.playerId) {
                        "Unexpected player turn to provide field: expected ${state.playerId}, got: ${action.playerId}"
                    }
                }

                when (action.playerId) {
                    player1Id -> {
                        require(!::player1Field.isInitialized) {
                            "Field for player1[$player1Id] was already provided"
                        }
                        player1Field = action.field
                        state = BattleshipStateAwaitingField(player2Id)
                    }
                    player2Id -> {
                        require(!::player2Field.isInitialized) {
                            "Field for player2[$player2Id] was already provided"
                        }
                        player2Field = action.field
                        state = BattleshipStateAwaitingStrike(player1Id)
                    }
                }
            }
            is BattleshipActionStrike -> {
                assertState<BattleshipStateAwaitingStrike> { state ->
                    require(state.attackerId == action.attackerId) {
                        "Unexpected player turn to strike: expected ${state.attackerId}, got: ${action.attackerId}"
                    }
                }

                val victimId = otherPlayerId(action.attackerId)
                val victimField = playerField(victimId)

                victimField.strikeAt(action.position)

                if (victimField.hasAliveShips()) {
                    state = BattleshipStateAwaitingStrike(victimId)
                } else {
                    state = BattleshipStateGameOver(action.attackerId)
                    accept(BattleshipActionGameOver(action.attackerId))
                }
            }
            is BattleshipActionGameOver -> {
                assertState<BattleshipStateGameOver> { state ->
                    require(state.winnerId == action.winnerId) {
                        "Unexpected winner: expected ${state.winnerId}, got: ${action.winnerId}"
                    }
                }
            }
        }
    }

    fun otherPlayerId(playerId: BattleshipPlayerId) =
        if (playerId == player1Id) player2Id else player1Id

    fun playerField(playerId: BattleshipPlayerId) =
        if (playerId == player1Id) player1Field else player2Field

    private inline fun <reified T : BattleshipState> assertState(assertPlayerId: (T) -> Unit) {
        require(state is T) { "Expected current game state to be ${typeOf<T>()}, got $state" }
        assertPlayerId(state as T)
    }
}
