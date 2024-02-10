plugins {
    id("battleship.kotlin-conventions")
    id("application")
}

group = "dev.spris"
version = "1.0-SNAPSHOT"

dependencies {
    implementation(project(":battleship-core"))
}

application {
    mainClass = "dev.spris.battleship.cli.MainKt"
}
