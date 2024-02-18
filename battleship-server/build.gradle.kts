plugins {
    id("org.springframework.boot") version "3.2.2"
    id("io.spring.dependency-management") version "1.1.4"
    id("battleship.kotlin-conventions")
    kotlin("plugin.spring") version "1.9.22"
    id("com.ncorti.ktfmt.gradle") version "0.17.0"
}

repositories {
    mavenCentral()
}

dependencies {
    implementation(project(":battleship-core"))

    implementation("org.springframework.boot:spring-boot-starter")
    implementation("org.springframework.boot:spring-boot-starter-webflux")
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("io.projectreactor:reactor-test")
}

ktfmt {
    kotlinLangStyle()
}
