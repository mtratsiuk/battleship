package dev.spris.battleship.server.grpc

import dev.spris.battleship.core.*
import dev.spris.battleship.proto.core.v1.*
import dev.spris.battleship.proto.server.v1.*
import dev.spris.battleship.proto.server.v1.gameLogEntryProto
import dev.spris.battleship.proto.server.v1.playerProto
import dev.spris.battleship.server.repository.GameState
import dev.spris.battleship.server.repository.Player

fun GameState.toProto() =
    when (this) {
        GameState.IDLE -> GameStateProto.IDLE
        GameState.RUNNING -> GameStateProto.RUNNING
        GameState.FINISHED -> GameStateProto.FINISHED
    }

fun Player.toProto() = playerProto {
    id = this@toProto.id.id
    name = this@toProto.name
}

fun BattleshipGameLogEntry.toProto() = gameLogEntryProto {
    when (val entry = this@toProto) {
        is BattleshipGameLogActionEntry ->
            when (val action = entry.action) {
                is BattleshipActionField -> {
                    field = battleshipActionFieldProto {
                        playerId = action.playerId.id
                        field = battleshipFieldProto { field = action.field.toString() }
                    }
                }
                is BattleshipActionGameOver -> {
                    gameOver = battleshipActionGameOverProto { winnerId = action.winnerId.id }
                }
                is BattleshipActionStrike -> {
                    strike = battleshipActionStrikeProto {
                        attackerId = action.attackerId.id
                        position = action.position.toProto()
                    }
                }
            }
        is BattleshipGameLogErrorEntry -> {
            error = entry.error.message ?: ""
        }
    }
}

fun BattleshipField.toProto() = battleshipFieldProto {
    field = this@toProto.toString()
    hits.addAll(this@toProto.hits.map { it.toProto() })
    misses.addAll(this@toProto.misses.map { it.toProto() })
}

fun BattleshipPos.toProto() = battleshipPosProto {
    x = this@toProto.x
    y = this@toProto.y
}
