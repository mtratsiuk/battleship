package main

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	core "github.com/mtratsiuk/battleship/battleship-go-core"
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
	gameRenderer *AppGameRenderer
}

type AppPlayer struct {
	id   string
	name string
	wins int
}

type AppGameRenderer struct {
	game        *pbserver.GameProto
	fields      map[string]*core.BattleshipField
	curEntryIdx int
	cancel      func()
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

	if err := g.SetKeybinding(VIEW_GAMES_LIST, gocui.KeyEnter, gocui.ModNone, app.FetchGame); err != nil {
		log.Panicln(err)
	}

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
		fmt.Fprintf(v, "enter: view game log, ")
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

	if app.gameRenderer == nil {
		v.Title = "Game"
		return nil
	}

	renderer := app.gameRenderer
	game := renderer.game

	v.Title = fmt.Sprintf("Game[%v]: %v vs %v", game.Id, game.Player_1.Name, game.Player_2.Name)

	renderField := func(f *core.BattleshipField) string {
		v := &strings.Builder{}

		cmiss := color.New(color.BgBlue).Add(color.BgWhite)
		cempty := color.New(color.BgBlue).Add(color.BgWhite)
		chit := color.New(color.BgRed).Add(color.BgWhite)
		cship := color.New(color.BgBlue).Add(color.BgWhite)

		fmt.Fprint(v, "  ")
		for y := range f.Field {
			fmt.Fprintf(v, "%v ", y)
		}
		fmt.Fprintln(v)

		for y, l := range f.Field {
			fmt.Fprintf(v, "%v ", y)

			for x, t := range l {
				pos := core.BattleshipPos{X: x, Y: y}

				if t.Kind == core.BattleshipTileKindEmpty {
					if f.Misses.Has(pos) {
						cmiss.Fprint(v, "o")
					} else {
						cempty.Fprint(v, ".")
					}
				} else {
					if f.Hits.Has(pos) {
						chit.Add(color.BgRed).Fprintf(v, string(t.Ship))
					} else {
						cship.Fprintf(v, string(t.Ship))
					}
				}

				cempty.Fprint(v, " ")
			}

			cempty.Fprintln(v)
		}

		return v.String()
	}

	f1 := strings.Split(renderField(renderer.fields[game.Player_1.Id]), "\n")
	f2 := strings.Split(renderField(renderer.fields[game.Player_2.Id]), "\n")

	pad := strings.Repeat(" ", 16)

	fmt.Fprintln(v)
	fmt.Fprintln(v)
	fmt.Fprintf(v, " %v", game.GetLog()[min(len(game.GetLog()) - 1, renderer.curEntryIdx)])

	fmt.Fprintln(v)
	fmt.Fprintln(v)

	for i := range f1 {
		fmt.Fprintf(v, " %v%v%v ", f1[i], pad, f2[i])
		fmt.Fprintln(v)
	}

	fmt.Fprintf(v, " %-38v", game.Player_1.Name)
	fmt.Fprintf(v, "%v", game.Player_2.Name)

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

	v, err = g.View(VIEW_GAME)
	if err != nil {
		return err
	}
	app.RenderGame(v)

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

func (app *App) FetchGame(g *gocui.Gui, v *gocui.View) error {
	if len(app.games) == 0 {
		app.err = "can't select a game"
		return app.ReRender(g)
	}

	game := app.games[app.curGameIdx]

	if game.State != pbserver.GameStateProto_FINISHED {
		app.err = "game is not finished yet"
		return app.ReRender(g)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := app.client.GetGame(ctx, &pbserver.GetGameRequest{Id: game.Id})
	if err != nil {
		app.err = err.Error()
	} else {
		if app.gameRenderer != nil {
			app.gameRenderer.cancel()
		}

		r, err := app.NewAppGameRenderer(g, response)

		if err != nil {
			return err
		}

		app.err = ""
		app.gameRenderer = r
	}

	return app.ReRender(g)
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

func (app *App) NewAppGameRenderer(g *gocui.Gui, p *pbserver.GetGameResponse) (*AppGameRenderer, error) {
	gr := &AppGameRenderer{}
	gr.game = p.Game
	gr.fields = make(map[string]*core.BattleshipField, 2)
	gr.curEntryIdx = 2

	p1id, field1, err := ExtractField(gr.game.Log[0])
	if err != nil {
		return gr, err
	}
	gr.fields[p1id] = &field1

	p2id, field2, err := ExtractField(gr.game.Log[1])
	if err != nil {
		return gr, err
	}
	gr.fields[p2id] = &field2

	ctx, cancel := context.WithCancel(context.Background())
	gr.cancel = cancel

	otherId := func(id string) string {
		if id == p1id {
			return p2id
		}
		return p1id
	}

	go func() {
		for gr.curEntryIdx < len(gr.game.Log) {
			select {
			case <-ctx.Done():
				return
			default:
				entry := gr.game.Log[gr.curEntryIdx]

				switch entry.GetAction().(type) {
				case *pbserver.GameLogEntryProto_Strike:
					strike := entry.GetStrike()
					gr.fields[otherId(strike.AttackerId)].Strike(core.NewBattleshipPosFromProto(strike.Position))
				default:
					// noop
				}

				gr.curEntryIdx += 1
				time.Sleep(time.Second / 2)

				g.Update(func(g *gocui.Gui) error {
					app.ReRender(g)
					return nil
				})
			}
		}
	}()

	return gr, nil
}

func ExtractField(p *pbserver.GameLogEntryProto) (string, core.BattleshipField, error) {
	field := p.GetField()

	if field == nil {
		return "", core.BattleshipField{}, fmt.Errorf("expected provided game log entry to be a field")
	}

	f, err := core.NewBattleshipFieldFromProto(field.Field)

	if err != nil {
		return "", core.BattleshipField{}, err
	}

	return field.PlayerId, f, nil
}
