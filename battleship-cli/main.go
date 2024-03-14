package main

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	pbserver "github.com/mtratsiuk/battleship/gen/proto/go/server/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PADDING = 1

	VIEW_GAMES_LIST = "games-list-view"
	VIEW_PLAYERS    = "players-view"
	VIEW_GAME       = "game-view"
	VIEW_ERRORS     = "errors-view"
	VIEW_FIELD_1    = "field-1-view"
	VIEW_FIELD_2    = "field-2-view"
	VIEW_HOT_KEYS   = "hot-keys-view"
)

type AppConfig struct {
	ServerGrpcUrl string
}

type App struct {
	cfg    AppConfig
	client pbserver.BattleshipServerServiceClient
	close  func()

	curView      string
	curGameIdx   int
	curPlayerIdx int
	err          string
	games        []*pbserver.GetGamesResponseEntry
	players      []AppPlayer
	game         *pbserver.GameProto
}

type AppPlayer struct {
	id   string
	name string
	wins int
}

func NewApp() App {
	cfg := AppConfig{}
	cfg.ServerGrpcUrl = "localhost:6969"

	app := App{}
	app.cfg = cfg
	app.games = make([]*pbserver.GetGamesResponseEntry, 0)
	app.players = make([]AppPlayer, 0)
	app.err = ""
	app.curView = VIEW_GAMES_LIST
	app.curGameIdx = 0
	app.curPlayerIdx = 0

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

	if err := g.SetKeybinding("", 'q', gocui.ModNone, app.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'r', gocui.ModNone, app.RefreshGames); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'a', gocui.ModNone, app.AddRandomBot); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, app.HandleArrowDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, app.HandleArrowUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", '1', gocui.ModNone, app.CreateFocusViewHandler(VIEW_GAMES_LIST)); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", '2', gocui.ModNone, app.CreateFocusViewHandler(VIEW_PLAYERS)); err != nil {
		log.Panicln(err)
	}

	go app.Poll(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (app *App) Layout(g *gocui.Gui) error {
	if err := app.GamesListView(g); err != nil {
		return err
	}

	if err := app.PlayersView(g); err != nil {
		return err
	}

	if err := app.GameView(g); err != nil {
		return err
	}

	if err := app.ErrorsView(g); err != nil {
		return err
	}

	if err := app.HotKeysView(g); err != nil {
		return err
	}

	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	if _, err := g.SetCurrentView(app.curView); err != nil {
		return err
	}

	return nil
}

func (app *App) GamesListView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_GAMES_LIST, PADDING, PADDING, maxX/2-PADDING, maxY/2-PADDING); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

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

func (app *App) GameView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_GAME, PADDING, maxY/2, maxX-PADDING, maxY-PADDING*3*2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		app.RenderGame(v)
	}
	return nil
}

func (app *App) ErrorsView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(VIEW_ERRORS, PADDING, maxY-PADDING*3*2, maxX-PADDING, maxY-PADDING*3-1); err != nil {
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
	v.Title = "1. Games"

	v.Highlight = app.curView == VIEW_GAMES_LIST
	v.SelFgColor = gocui.ColorGreen

	_, h := v.Size()
	curPage := app.curGameIdx / h

	v.SetOrigin(0, curPage*h)
	v.SetCursor(0, app.curGameIdx%h)

	for i, game := range app.games {
		idxL := strconv.Itoa(i)
		if i == app.curGameIdx && v.Highlight {
			idxL = ">"
		}
		idxR := strconv.Itoa(i)
		if i == app.curGameIdx && v.Highlight {
			idxR = "<"
		}

		p1 := game.Player_1.Name
		if game.WinnerId != nil && *game.WinnerId == game.Player_1.Id {
			p1 = "[W] " + p1
		}

		p2 := game.Player_2.Name
		if game.WinnerId != nil && *game.WinnerId == game.Player_2.Id {
			p2 += " [W]"
		}

		fmt.Fprintf(v, "%4v) %12v %32v vs %-32v (%v\n", idxL, game.State, p1, p2, idxR)
	}

	return nil
}

func (app *App) RenderPlayers(v *gocui.View) error {
	v.Clear()
	v.Title = "2. Leaderboard"

	v.Highlight = app.curView == VIEW_PLAYERS
	v.SelFgColor = gocui.ColorGreen

	_, h := v.Size()
	curPage := app.curPlayerIdx / h

	v.SetOrigin(0, curPage*h)
	v.SetCursor(0, app.curPlayerIdx%h)

	for i, p := range app.players {
		idxL := strconv.Itoa(i)
		if i == app.curPlayerIdx && v.Highlight {
			idxL = ">"
		}
		idxR := strconv.Itoa(i)
		if i == app.curPlayerIdx && v.Highlight {
			idxR = "<"
		}

		fmt.Fprintf(v, "%3v) %-32v %4v (%v\n", idxL, p.name, p.wins, idxR)
	}

	return nil
}

func (app *App) RenderGame(v *gocui.View) error {
	v.Clear()
	v.Title = "Game"

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
		app.SetGames(games)
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

func (app *App) HandleArrowDown(g *gocui.Gui, v *gocui.View) error {
	if app.curView == VIEW_GAMES_LIST {
		app.curGameIdx = min(len(app.games)-1, app.curGameIdx+1)
	}

	if app.curView == VIEW_PLAYERS {
		app.curPlayerIdx = min(len(app.players)-1, app.curPlayerIdx+1)
	}

	return app.ReRender(g)
}

func (app *App) HandleArrowUp(g *gocui.Gui, v *gocui.View) error {
	if app.curView == VIEW_GAMES_LIST {
		app.curGameIdx = max(0, app.curGameIdx-1)
	}

	if app.curView == VIEW_PLAYERS {
		app.curPlayerIdx = max(0, app.curPlayerIdx-1)
	}

	return app.ReRender(g)
}

func (app *App) CreateFocusViewHandler(name string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		app.curView = name

		if _, err := g.SetCurrentView(name); err != nil {
			return err
		}

		return app.ReRender(g)
	}
}

func (app *App) Poll(g *gocui.Gui) {
	for {
		time.Sleep(time.Second)
		g.Update(func(g *gocui.Gui) error {
			return app.RefreshGames(g, g.CurrentView())
		})
	}
}

func (app *App) SetGames(gamesResponse *pbserver.GetGamesResponse) {
	app.games = gamesResponse.Games

	players := make(map[string]*AppPlayer, 0)

	for _, game := range app.games {
		players[game.Player_1.Id] = &AppPlayer{game.Player_1.Id, game.Player_1.Name, 0}
		players[game.Player_2.Id] = &AppPlayer{game.Player_2.Id, game.Player_2.Name, 0}
	}

	for _, game := range app.games {
		p1 := game.Player_1
		p2 := game.Player_2

		if game.WinnerId != nil && p1.Id == *game.WinnerId {
			players[p1.Id].wins += 1
		}

		if game.WinnerId != nil && p2.Id == *game.WinnerId {
			players[p2.Id].wins += 1
		}
	}

	playersList := make([]AppPlayer, 0)

	for _, p := range players {
		playersList = append(playersList, *p)
	}

	slices.SortFunc(playersList, func(a, b AppPlayer) int { return strings.Compare(a.name, b.name) })
	slices.SortFunc(playersList, func(a, b AppPlayer) int { return b.wins - a.wins })

	app.players = playersList
}
