package dev.spris.battleship.server.service

import java.nio.ByteBuffer
import java.util.Base64
import java.util.UUID
import org.springframework.stereotype.Component

@Component
class IdGenerator {
    fun next(): String {
        val uuid = UUID.randomUUID()
        val bb = ByteBuffer.wrap(ByteArray(16))

        bb.putLong(uuid.mostSignificantBits)
        bb.putLong(uuid.leastSignificantBits)

        return Base64.getEncoder().encodeToString(bb.array()).trimEnd('=')
    }
}
