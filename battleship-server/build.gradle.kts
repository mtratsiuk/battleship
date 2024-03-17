plugins {
    id("org.springframework.boot") version "3.2.2"
    id("io.spring.dependency-management") version "1.1.4"
    id("battleship.kotlin-conventions")
    kotlin("plugin.spring") version "1.9.22"
    id("com.ncorti.ktfmt.gradle") version "0.17.0"
    id("com.google.protobuf") version "0.9.4"
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

    implementation("io.grpc:grpc-netty-shaded:1.62.2")
    implementation("io.grpc:grpc-kotlin-stub:1.4.1")
    implementation("io.grpc:grpc-protobuf:1.62.2")
    implementation("io.grpc:grpc-services:1.62.2")
    implementation("com.google.protobuf:protobuf-kotlin:3.25.3")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("io.projectreactor:reactor-test")
}

ktfmt {
    kotlinLangStyle()
}

sourceSets.main {
    java.srcDirs(
        "../gen/proto/java",
        "../gen/proto/kotlin",
    )
}
