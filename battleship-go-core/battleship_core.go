package core

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mtratsiuk/adventofcode/gotils"
	pbcore "github.com/mtratsiuk/battleship/gen/proto/go/core/v1"
	pbserver "github.com/mtratsiuk/battleship/gen/proto/go/server/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const BattleshipFieldSize = 10

type BattleshipKind rune

const (
	BattleshipKindPatrolBoat BattleshipKind = 'P'
	BattleshipKindSubmarine  BattleshipKind = 'S'
	BattleshipKindDestroyer  BattleshipKind = 'D'
	BattleshipKindBattleship BattleshipKind = 'B'
	BattleshipKindCarrier    BattleshipKind = 'C'
)

var BattleshipKinds = []BattleshipKind{
	BattleshipKindPatrolBoat,
	BattleshipKindSubmarine,
	BattleshipKindDestroyer,
	BattleshipKindBattleship,
	BattleshipKindCarrier,
}

var BattleshipKindSizes = map[BattleshipKind]int{
	BattleshipKindPatrolBoat: 2,
	BattleshipKindSubmarine:  3,
	BattleshipKindDestroyer:  3,
	BattleshipKindBattleship: 4,
	BattleshipKindCarrier:    5,
}

func (b BattleshipKind) IsBattleshipKind() bool {
	return slices.Contains(BattleshipKinds, b)
}

func (b BattleshipKind) Size() int {
	return BattleshipKindSizes[b]
}

type BattleshipTileKind int

const (
	BattleshipTileKindEmpty BattleshipTileKind = iota
	BattleshipTileKindShip
)

type BattleshipTile struct {
	Kind BattleshipTileKind
	Ship BattleshipKind
}

func (b BattleshipTile) IsEmpty() bool {
	return b.Kind == BattleshipTileKindEmpty
}

func NewEmptyBattleshipTile() BattleshipTile {
	return BattleshipTile{BattleshipTileKindEmpty, -1}
}

func NewBattleshipTile(ship BattleshipKind) BattleshipTile {
	return BattleshipTile{BattleshipTileKindShip, ship}
}

type BattleshipPos struct {
	X, Y int
}

func NewBattleshipPosFromProto(p *pbcore.BattleshipPosProto) BattleshipPos {
	pos := BattleshipPos{}

	pos.X = int(p.X)
	pos.Y = int(p.Y)

	return pos
}

type BattleshipField struct {
	Field  [BattleshipFieldSize][BattleshipFieldSize]BattleshipTile
	Hits   gotils.Set[BattleshipPos]
	Misses gotils.Set[BattleshipPos]
}

func NewBattleshipField() BattleshipField {
	bf := BattleshipField{}
	bf.Hits = gotils.NewSet[BattleshipPos]()
	bf.Misses = gotils.NewSet[BattleshipPos]()

	return bf
}

func (b *BattleshipField) ToProto() *pbcore.BattleshipFieldProto {
	bp := pbcore.BattleshipFieldProto{}
	bp.Hits = make([]*pbcore.BattleshipPosProto, 0)
	bp.Misses = make([]*pbcore.BattleshipPosProto, 0)

	var sb strings.Builder

	for _, l := range b.Field {
		for _, t := range l {
			if t.Kind == BattleshipTileKindEmpty {
				sb.WriteByte('.')
			} else {
				sb.WriteRune(rune(t.Ship))
			}
		}
		sb.WriteByte('\n')
	}

	bp.Field = strings.TrimSpace(sb.String())

	return &bp
}

func NewBattleshipFieldFromProto(p *pbcore.BattleshipFieldProto) (BattleshipField, error) {
	bf := NewBattleshipField()

	for _, h := range p.Hits {
		bf.Hits.Add(NewBattleshipPosFromProto(h))
	}

	for _, m := range p.Misses {
		bf.Misses.Add(NewBattleshipPosFromProto(m))
	}

	for y, l := range strings.Split(p.Field, "\n") {
		for x, c := range l {
			if c == '.' {
				bf.Field[y][x] = NewEmptyBattleshipTile()
			} else {
				ship := BattleshipKind(c)

				if !ship.IsBattleshipKind() {
					return bf, fmt.Errorf("unexpected battleship char code %v", c)
				}

				bf.Field[y][x] = NewBattleshipTile(ship)
			}
		}
	}

	return bf, nil
}

func (b *BattleshipField) Strike(pos BattleshipPos) {
	if b.Field[pos.Y][pos.X].Kind == BattleshipTileKindEmpty {
		b.Misses.Add(pos)
	} else {
		b.Hits.Add(pos)
	}
}

func NewBattleshipServerServiceClient() (*pbserver.BattleshipServerServiceClient, func(), error) {
	grpcHost := EnvOr("BATTLESHIP_SERVER_GRPC_HOST", "localhost")
	grpcPort := EnvOr("BATTLESHIP_SERVER_GRPC_PORT", "6969")
	serverGrpcUrl := fmt.Sprintf("%v:%v", grpcHost, grpcPort)

	conn, err := grpc.Dial(serverGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pbserver.NewBattleshipServerServiceClient(conn)

	return &client, func() { conn.Close() }, nil
}

func EnvOr(name, fallback string) string {
	val, ok := os.LookupEnv(name)

	if !ok {
		return fallback
	}

	return val
}
