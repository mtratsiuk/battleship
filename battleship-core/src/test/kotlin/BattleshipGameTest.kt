import dev.spris.battleship.core.*
import org.junit.jupiter.api.Test

import org.junit.jupiter.api.Assertions.*

class BattleshipGameTest {

    @Test
    fun `BattleshipGame ensures correct turns order and game state`() {
        val gameId = BattleshipGameId("gameId")
        val player1Id = BattleshipPlayerId("player1Id")
        val player2Id = BattleshipPlayerId("player2Id")
        val player1Field = BattleshipField.fromShips(*createValidBattleshipField().toTypedArray())
        val player2Field = BattleshipField.fromShips(*createValidBattleshipField().toTypedArray())

        val game = BattleshipGame(gameId, player1Id, player2Id)

        assertEquals(game.state, BattleshipStateAwaitingField(player1Id))
        game.accept(BattleshipActionField(player1Id, player1Field))

        assertEquals(game.state, BattleshipStateAwaitingField(player2Id))
        game.accept(BattleshipActionField(player2Id, player2Field))

        for (y in 0..<player2Field.field.size) {
            for (x in 0..<player2Field.field.size) {
                if (player2Field.field[y][x] is BattleshipShipTile) {
                    // Player 1 strikes ship tiles
                    assertEquals(game.state, BattleshipStateAwaitingStrike(player1Id))
                    game.accept(BattleshipActionStrike(player1Id, BattleshipPos(x, y)))

                    if (player2Field.hasAliveShips()) {
                        // Player 2 misses
                        assertEquals(game.state, BattleshipStateAwaitingStrike(player2Id))
                        game.accept(BattleshipActionStrike(player2Id, BattleshipPos(9, 9)))
                    } else {
                        assertEquals(game.state, BattleshipStateGameOver(player1Id))
                    }
                }
            }
        }
    }
}
