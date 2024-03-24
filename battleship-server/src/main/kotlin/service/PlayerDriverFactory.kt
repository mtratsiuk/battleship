package dev.spris.battleship.server.service

import com.google.common.cache.CacheBuilder
import com.google.common.cache.CacheLoader
import dev.spris.battleship.core.*
import dev.spris.battleship.proto.bot.v1.BattleshipBotServiceGrpcKt
import dev.spris.battleship.proto.bot.v1.getFieldRequest
import dev.spris.battleship.proto.bot.v1.getStrikeRequest
import dev.spris.battleship.server.grpc.toDomain
import dev.spris.battleship.server.grpc.toOtherFieldProto
import dev.spris.battleship.server.grpc.toProto
import dev.spris.battleship.server.repository.Player
import io.grpc.ManagedChannelBuilder
import kotlin.random.Random
import kotlinx.coroutines.delay
import org.springframework.stereotype.Service

@Service
class PlayerDriverFactory {
    private val grpcPlayers =
        CacheBuilder.newBuilder()
            .build<Player, GrpcPlayerDriver>(
                CacheLoader.from { player -> GrpcPlayerDriver(player) }
            )

    suspend fun create(player: Player): PlayerDriver {
        if (player.addr.startsWith("inprocess")) {
            return InProcessRandomPlayerDriver(player)
        }

        return grpcPlayers.get(player)
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

class GrpcPlayerDriver(
    player: Player,
) : PlayerDriver {
    private val channel = ManagedChannelBuilder.forTarget(player.addr).usePlaintext().build()

    private val stub = BattleshipBotServiceGrpcKt.BattleshipBotServiceCoroutineStub(channel)

    override suspend fun requestField(gameId: BattleshipGameId): BattleshipField {
        val request = getFieldRequest { this.gameId = gameId.id }
        val response = stub.getField(request)

        return BattleshipField(
            field = BattleshipField.fieldArrayFromString(response.field),
            hits = mutableSetOf(),
            misses = mutableSetOf(),
        )
    }

    override suspend fun requestStrike(
        gameId: BattleshipGameId,
        ownField: BattleshipField,
        otherField: BattleshipField
    ): BattleshipPos {
        val request = getStrikeRequest {
            this.gameId = gameId.id
            this.ownField = ownField.toProto()
            this.otherField = otherField.toOtherFieldProto()
        }

        val response = stub.getStrike(request)

        return response.pos.toDomain()
    }
}

class InProcessRandomPlayerDriver(
    private val player: Player,
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
        delay(10)

        return BattleshipPos(
            Random.nextInt(otherField.field.size),
            Random.nextInt(otherField.field.size)
        )
    }
}
