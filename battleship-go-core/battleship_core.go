package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mtratsiuk/adventofcode/gotils"
	pbcore "github.com/mtratsiuk/battleship/gen/proto/go/core/v1"
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

func (b BattleshipKind) IsBattleshipKind() bool {
	return slices.Contains(BattleshipKinds, b)
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
