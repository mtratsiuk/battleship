package main

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/jroimartin/gocui"
	pbserver "github.com/mtratsiuk/battleship/gen/proto/go/server/v1"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PADDING = 1

	VIEW_GAMES_LIST = "games-list-view"
	VIEW_PLAYERS    = "players-view"
	VIEW_ERRORS     = "errors-view"
	VIEW_FIELD_1    = "field-1-view"
	VIEW_FIELD_2    = "field-2-view"
	VIEW_HOT_KEYS   = "hot-keys-view"
)

type Config struct {
	ServerGrpcUrl string
}

type App struct {
	cfg    Config
	client pbserver.BattleshipServerServiceClient
	close  func()

	err   string
	games []*pbserver.GetGamesResponseEntry
}

func NewApp() App {
	cfg := Config{}
	cfg.ServerGrpcUrl = "localhost:6969"

	app := App{}
	app.cfg = cfg
	app.games = make([]*pbserver.GetGamesResponseEntry, 0)
	app.err = ""

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
	g.Mouse = true
	g.SetCurrentView(VIEW_GAMES_LIST)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, app.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, app.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'r', gocui.ModNone, app.RefreshGames); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'a', gocui.ModNone, app.AddRandomBot); err != nil {
		log.Panicln(err)
	}

	go app.Poll(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (app *App) Layout(g *gocui.Gui) error {
	if err := app.GamesListView(g); err != nil {
		log.Panicln(err)
	}

	if err := app.PlayersView(g); err != nil {
		log.Panicln(err)
	}

	if err := app.ErrorsView(g); err != nil {
		log.Panicln(err)
	}

	if err := app.HotKeysView(g); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (app *App) GamesListView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_GAMES_LIST, PADDING, PADDING, maxX/2-PADDING, maxY/2-PADDING); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = true
		app.RenderGamesList(v)
	}
	return nil
}

func (app *App) PlayersView(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_PLAYERS, maxX/2+PADDING, PADDING, maxX-PADDING, maxY/2-PADDING); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		app.RenderPlayers(v)
	}
	return nil
}

func (app *App) ErrorsView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_ERRORS, PADDING, maxY-PADDING*3*2, maxX-PADDING, maxY-PADDING*3 - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		return app.RenderErrors(v)
	}
	return nil
}

func (app *App) HotKeysView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_HOT_KEYS, PADDING, maxY-PADDING*3, maxX-PADDING, maxY-PADDING); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, "r: refresh games, ")
		fmt.Fprintf(v, "a: add random bot, ")
		fmt.Fprintf(v, "ctrl+c / q: exit")
	}

	return nil
}

func (app *App) RenderGamesList(v *gocui.View) error {
	v.Clear()
	v.Title = "Games"

	for i, game := range app.games {
		p1 := game.Player_1.Name
		if game.WinnerId != nil && *game.WinnerId == game.Player_1.Id {
			p1 = "[W] " + p1
		}

		p2 := game.Player_2.Name
		if game.WinnerId != nil && *game.WinnerId == game.Player_2.Id {
			p2 += " [W]"
		}

		fmt.Fprintf(v, "%4v) %12v %32v vs %-32v (%v\n", i, game.State, p1, p2, i)
	}

	return nil
}

func (app *App) RenderPlayers(v *gocui.View) error {
	v.Clear()
	v.Title = "Leaderboard"

	players := make(map[string]int, 0)

	for _, game := range app.games {
		p1 := game.Player_1
		p2 := game.Player_2

		if game.WinnerId != nil && p1.Id == *game.WinnerId {
			players[p1.Name] += 1
		}

		if game.WinnerId != nil && p2.Id == *game.WinnerId {
			players[p2.Name] += 1
		}
	}

	playersList := maps.Keys(players)
	slices.Sort(playersList)
	slices.SortFunc(playersList, func(a, b string) int { return players[b] - players[a] })

	for i, p := range playersList {
		fmt.Fprintf(v, "%3v) %-32v %v (%v\n", i, p, players[p], i)
	}

	return nil
}

func (app *App) RenderErrors(v *gocui.View) error {
	v.Clear()
	v.Title = "Errors"

	fmt.Fprintf(v, "%v", app.err)

	return nil
}

func (app *App) ReRender(g *gocui.Gui) error {
	v, err := g.View(VIEW_GAMES_LIST)
	if err != nil {
		return err
	}
	app.RenderGamesList(v)

	v, err = g.View(VIEW_PLAYERS)
	if err != nil {
		return err
	}
	app.RenderPlayers(v)

	v, err = g.View(VIEW_ERRORS)
	if err != nil {
		return err
	}
	app.RenderErrors(v)

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
		app.err = err.Error()
	} else {
		app.err = ""
		app.games = games.Games
	}

	return app.ReRender(g)
}

func (app *App) AddRandomBot(g *gocui.Gui, v *gocui.View) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := app.client.AddRandomBot(ctx, &pbserver.AddRandomBotRequest{})
	if err != nil {
		app.err = err.Error()
	} else {
		app.err = ""
	}

	return app.RefreshGames(g, v)
}

func (app *App) Poll(g *gocui.Gui) {
	for {
		time.Sleep(time.Second)
		g.Update(func(g *gocui.Gui) error {
			return app.RefreshGames(g, g.CurrentView())
		})
	}
}
