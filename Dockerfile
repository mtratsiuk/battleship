FROM fedora:39

COPY --from=golang:1.22.1 /usr/local/go/ /usr/local/go/
COPY --from=bufbuild/buf /usr/local/bin/buf /usr/local/go/bin/buf
COPY --from=amazoncorretto:21-al2023-headless /usr/lib/jvm/java-21-amazon-corretto /usr/lib/jvm/java-21-amazon-corretto

ENV PATH="/usr/local/go/bin:/usr/lib/jvm/java-21-amazon-corretto/bin:${PATH}"
ENV JAVA_HOME=/usr/lib/jvm/java-21-amazon-corretto
ENV GRADLE_USER_HOME=/gradle

COPY gradle /battleship/gradle
COPY gradlew /battleship/gradlew

RUN /battleship/gradlew --version

WORKDIR /battleship

ENTRYPOINT [ "/bin/bash", "-c" ]

