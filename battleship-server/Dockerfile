FROM battleship-builder as builder

COPY . /battleship

RUN /battleship/bin/build.sh

FROM amazoncorretto:21-al2023-headless

ENV PATH="/usr/lib/jvm/java-21-amazon-corretto/bin:${PATH}"
ENV JAVA_HOME=/usr/lib/jvm/java-21-amazon-corretto

COPY --from=builder /battleship/battleship-server/build/libs/battleship-server.jar /battleship/app.jar

WORKDIR /battleship

ENTRYPOINT [ "java", "-jar", "app.jar" ]
