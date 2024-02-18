package dev.spris.battleship.server

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication class BattleshipServerApplication

fun main(args: Array<String>) {
    runApplication<BattleshipServerApplication>(*args)
}
