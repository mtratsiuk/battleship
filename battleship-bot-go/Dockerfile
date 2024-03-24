FROM battleship-builder as builder

COPY battleship-bot-go /battleship/battleship-bot-go
COPY battleship-go-core /battleship/battleship-go-core
COPY gen /battleship/gen
COPY go.mod /battleship/go.mod
COPY go.sum /battleship/go.sum

RUN go mod download
RUN go build -o /battleship/battleship-bot-go/dist ./battleship-bot-go/battleship_bot.go

FROM fedora:39

COPY --from=builder /battleship/battleship-bot-go/dist /battleship/app

WORKDIR /battleship

ENTRYPOINT [ "./app" ]
