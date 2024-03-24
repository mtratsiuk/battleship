FROM battleship-builder as builder

COPY battleship-cli /battleship/battleship-cli
COPY battleship-go-core /battleship/battleship-go-core
COPY gen /battleship/gen
COPY go.mod /battleship/go.mod
COPY go.sum /battleship/go.sum

RUN go mod download
RUN go build -o /battleship/battleship-cli/dist ./battleship-cli/battleship_cli.go

FROM fedora:39

COPY --from=builder /battleship/battleship-cli/dist /battleship/app

ARG BATTLESHIP_SERVER_GRPC_HOST
ARG BATTLESHIP_SERVER_GRPC_PORT

ENV BATTLESHIP_SERVER_GRPC_HOST $BATTLESHIP_SERVER_GRPC_HOST
ENV BATTLESHIP_SERVER_GRPC_PORT $BATTLESHIP_SERVER_GRPC_PORT

WORKDIR /battleship

ENTRYPOINT [ "./app" ]
