package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"time"

	core "github.com/mtratsiuk/battleship/battleship-go-core"
	pbbot "github.com/mtratsiuk/battleship/gen/proto/go/bot/v1"
	pbcore "github.com/mtratsiuk/battleship/gen/proto/go/core/v1"
	pbserver "github.com/mtratsiuk/battleship/gen/proto/go/server/v1"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	botServer := NewBotServer()

	client, close, err := core.NewBattleshipServerServiceClient()
	if err != nil {
		botServer.logger.Error(fmt.Sprintf("failed to connect to the bot runner gRPC server: %v", err))
		os.Exit(1)
	}
	defer close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	request := &pbserver.JoinLobbyRequest{
		Addr: botServer.config.externalAddr,
		Name: botServer.config.botName,
	}

	botServer.logger.Info(fmt.Sprintf("joining lobby with a delay... %v", request))

	time.Sleep(time.Second * 5)
	resp, err := (*client).JoinLobby(ctx, request)
	if err != nil {
		botServer.logger.Error(fmt.Sprintf("failed to join lobby: %v", err))
		os.Exit(1)
	}
	botServer.logger.Info(fmt.Sprintf("joined lobby: %v", resp))

	gprcUrl := fmt.Sprintf("%v:%v", botServer.config.grpcServerHost, botServer.config.grpcServerPort)
	lis, err := net.Listen("tcp", gprcUrl)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	pbbot.RegisterBattleshipBotServiceServer(grpcServer, &botServer)
	reflection.Register(grpcServer)

	botServer.logger.Info(fmt.Sprintf("Starting gRPC server at: %v", gprcUrl))
	if err := grpcServer.Serve(lis); err != nil {
		botServer.logger.Error(fmt.Sprintf("failed to servce gRPC: %v", err))
		os.Exit(1)
	}
}

type CtxKey string

const (
	CtxKeyMethod = CtxKey("Method")
	CtxKeyGameId = CtxKey("GameId")
)

type Config struct {
	grpcServerHost string
	grpcServerPort string
	externalAddr   string
	botName        string
}

func NewConfig() Config {
	c := Config{}

	c.grpcServerHost = core.EnvOr("BATTLESHIP_BOT_GO_GRPC_HOST", "0.0.0.0")
	c.grpcServerPort = core.EnvOr("BATTLESHIP_BOT_GO_GRPC_PORT", "6968")
	c.externalAddr = core.EnvOr("BATTLESHIP_BOT_GO_EXTERNAL_ADDR", "0.0.0.0:6968")
	c.botName = core.EnvOr("BATTLESHIP_BOT_GO_NAME", "Go Bot")

	return c
}

type BotServer struct {
	pbbot.UnimplementedBattleshipBotServiceServer

	config Config
	logger *slog.Logger
}

func NewBotServer() BotServer {
	b := BotServer{}
	b.config = NewConfig()
	b.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return b
}

func (b *BotServer) GetField(ctx context.Context, request *pbbot.GetFieldRequest) (*pbbot.GetFieldResponse, error) {
	ctx = context.WithValue(ctx, CtxKeyMethod, "GetField")
	ctx = context.WithValue(ctx, CtxKeyGameId, request.GameId)
	b.logger.InfoContext(ctx, "Received GetField request")

	f := core.NewBattleshipField()

	for _, ship := range core.BattleshipKinds {
		for {
			if ok := tryToAddShip(&f, ship); ok {
				break
			}
		}
	}

	resp := pbbot.GetFieldResponse{Field: f.ToProto().Field}

	return &resp, nil
}

func (b *BotServer) GetStrike(ctx context.Context, request *pbbot.GetStrikeRequest) (*pbbot.GetStrikeResponse, error) {
	ctx = context.WithValue(ctx, CtxKeyMethod, "GetStrike")
	ctx = context.WithValue(ctx, CtxKeyGameId, request.GameId)
	b.logger.InfoContext(ctx, "Received GetStrike request")

	x, y, err := findStrikePos(request)

	if err != nil {
		b.logger.WarnContext(ctx, err.Error())
		return nil, err
	}

	resp := &pbbot.GetStrikeResponse{Pos: &pbcore.BattleshipPosProto{X: int32(x), Y: int32(y)}}

	return resp, nil
}

func tryToAddShip(f *core.BattleshipField, ship core.BattleshipKind) bool {
	shipSize := ship.Size()
	horizontal := true
	if rand.Intn(2) == 1 {
		horizontal = false
	}

	curShip := make([]core.BattleshipPos, 0)
	x := rand.Intn(core.BattleshipFieldSize)
	y := rand.Intn(core.BattleshipFieldSize)

	if horizontal {
		x = min(x, core.BattleshipFieldSize-shipSize)
	} else {
		y = min(y, core.BattleshipFieldSize-shipSize)
	}

	for d := 0; d < shipSize; d += 1 {
		cx := x
		cy := y

		if horizontal {
			cx += d
		} else {
			cy += d
		}

		if f.Field[cy][cx].Kind != core.BattleshipTileKindEmpty {
			return false
		}

		curShip = append(curShip, core.BattleshipPos{X: cx, Y: cy})
	}

	for _, pos := range curShip {
		f.Field[pos.Y][pos.X] = core.NewBattleshipTile(ship)
	}

	return true
}

func findStrikePos(request *pbbot.GetStrikeRequest) (int, int, error) {
	ps := make(map[int]struct{}, 0)

	for y := 0; y < core.BattleshipFieldSize; y += 1 {
		for x := 0; x < core.BattleshipFieldSize; x += 1 {
			key := x*10 + y
			ps[key] = struct{}{}
		}
	}

	for _, h := range request.OtherField.Hits {
		delete(ps, int(h.X*10+h.Y))
	}

	for _, m := range request.OtherField.Misses {
		delete(ps, int(m.X*10+m.Y))
	}

	toStrike := maps.Keys(ps)

	if len(toStrike) == 0 {
		return 0, 0, fmt.Errorf("nowhere left to strike")
	}

	strike := toStrike[rand.Intn(len(toStrike))]
	return strike / 10, strike % 10, nil
}
