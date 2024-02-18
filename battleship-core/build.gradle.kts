plugins {
    id("battleship.kotlin-conventions")
    id("com.ncorti.ktfmt.gradle") version "0.17.0"
}

ktfmt {
    kotlinLangStyle()
}
