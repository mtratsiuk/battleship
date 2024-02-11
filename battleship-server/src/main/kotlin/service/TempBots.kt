package dev.spris.battleship.server.service

import jakarta.annotation.PostConstruct
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import org.springframework.stereotype.Service

@Service
class TempBots(
    val gameLobby: GameLobby,
) {
    @PostConstruct
    fun init() {
        GlobalScope.launch {
            gameLobby.join("inprocess://1", "first")
            gameLobby.join("inprocess://2", "second")

            delay(1000)
            gameLobby.join("inprocess://3", "third")
        }
    }
}
