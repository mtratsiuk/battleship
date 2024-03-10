package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	pbserver "github.com/mtratsiuk/battleship/gen/proto/go/server/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PADDING = 1

	VIEW_GAMES_LIST = "games-list-view"
)

type Config struct {
	ServerGrpcUrl string
}

type App struct {
	cfg    Config
	client pbserver.BattleshipServerServiceClient
	close  func()

	games []*pbserver.GetGamesResponseEntry
}

func NewApp() App {
	cfg := Config{}
	cfg.ServerGrpcUrl = "localhost:6969"

	app := App{}
	app.cfg = cfg
	app.games = make([]*pbserver.GetGamesResponseEntry, 0)

	conn, err := grpc.Dial(cfg.ServerGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicln(err)
	}
	app.close = func() { conn.Close() }

	app.client = pbserver.NewBattleshipServerServiceClient(conn)

	return app
}

func main() {
	app := NewApp()
	defer app.close()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(app.Layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, app.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'r', gocui.ModNone, app.RefreshGames); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (app *App) Layout(g *gocui.Gui) error {
	if err := app.GamesListView(g); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (app *App) GamesListView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_GAMES_LIST, PADDING, PADDING, maxX-PADDING, maxY/2-PADDING); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		app.RenderGamesList(v)
	}
	return nil
}

func (app *App) RenderGamesList(v *gocui.View) error {
	v.Clear()
	v.Title = "Games"

	for _, game := range app.games {
		p1 := game.Player_1.Name
		if game.WinnerId != nil && *game.WinnerId == game.Player_1.Id {
			p1 += " [W]"
		}

		p2 := game.Player_2.Name
		if game.WinnerId != nil && *game.WinnerId == game.Player_2.Id {
			p2 += " [W]"
		}

		fmt.Fprintf(v, "[%v]: %v %v vs %v\n", game.Id, game.State, p1, p2)
	}

	return nil
}

func (app *App) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (app *App) RefreshGames(g *gocui.Gui, _ *gocui.View) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	games, err := app.client.GetGames(ctx, &pbserver.GetGamesRequest{})
	if err != nil {
		return err
	}

	app.games = games.Games
	v, err := g.View(VIEW_GAMES_LIST)

	if err != nil {
		return err
	}

	app.RenderGamesList(v)

	return nil
}
