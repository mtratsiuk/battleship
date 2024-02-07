plugins {
    kotlin("jvm") version "1.9.22"
    id("application")
}

group = "dev.spris"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation(project(":battleship-core"))

    testImplementation("org.jetbrains.kotlin:kotlin-test")
}

tasks.test {
    useJUnitPlatform()
}

kotlin {
    jvmToolchain(21)
}

application {
    mainClass = "dev.spris.battleship.cli.MainKt"
}
