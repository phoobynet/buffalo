package main

import (
	"context"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/buffalo/data/configuration"
	"github.com/phoobynet/buffalo/data/market/clock"
	"github.com/phoobynet/buffalo/data/market/stock"
	"github.com/phoobynet/buffalo/data/market/stock/bar"
	"github.com/phoobynet/buffalo/data/metadata/asset"
	"github.com/phoobynet/buffalo/data/metadata/calendar"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

// App struct
type App struct {
	mut                        sync.Mutex
	ctx                        context.Context
	streamCtx                  context.Context
	cancelStream               context.CancelFunc
	db                         *gorm.DB
	assetRepository            *asset.Repository
	alpacaClient               *alpaca.Client
	marketDataClient           *marketdata.Client
	stockStream                *stock.Stream
	ready                      bool
	trades                     chan stream.Trade
	quotes                     chan stream.Quote
	bars                       chan stream.Bar
	currentSymbol              string
	calendarRepository         *calendar.Repository
	appConfigurationRepository *configuration.Repository
	status                     chan clock.Status
	statusClock                *clock.Clock
	barRepository              *bar.Repository
}

// NewApp creates a new App application struct
func NewApp() *App {
	db, err := gorm.Open(sqlite.Open("buffalo.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}

	trades := make(chan stream.Trade, 10_000)
	quotes := make(chan stream.Quote, 10_000)
	bars := make(chan stream.Bar, 100)

	status := make(chan clock.Status, 1)
	updateTicker := time.NewTicker(100 * time.Millisecond)
	snapshotTicker := time.NewTicker(1 * time.Second)

	app := &App{
		db:     db,
		trades: trades,
		quotes: quotes,
		bars:   bars,
		status: status,
	}

	go func(app *App) {
		var lastTrade stream.Trade
		var lastTradeEmitted stream.Trade
		var lastQuote stream.Quote
		var lastQuoteEmitted stream.Quote
		var lastBar stream.Bar
		var lastBarEmitted stream.Bar

		for {
			select {
			case lastTrade = <-app.trades:
			case lastQuote = <-app.quotes:
			case lastBar = <-app.bars:
			case currentStatus := <-app.status:
				app.Emit(currentStatus)
			case <-updateTicker.C:
				if lastTradeEmitted.ID != lastTrade.ID {
					app.Emit(lastTrade)
					lastTradeEmitted = lastTrade
				}

				if !lastQuoteEmitted.Timestamp.Equal(lastQuote.Timestamp) {
					app.Emit(lastQuote)
					lastQuoteEmitted = lastQuote
				}

				if !lastBarEmitted.Timestamp.Equal(lastBar.Timestamp) {
					app.Emit(lastBar)
					lastBarEmitted = lastBar
				}
			case t := <-snapshotTicker.C:
				if t.Second() == 0 && app.marketDataClient != nil && app.currentSymbol != "" {
					log.Println("Getting snapshot")
					snapshot := app.GetSnapshot(app.currentSymbol)
					app.Emit(snapshot)
				}
			}
		}
	}(app)

	return app
}

func (a *App) Emit(data any) {
	if a.ctx == nil {
		return
	}

	var eventName string

	switch data.(type) {
	case stream.Trade:
		eventName = "trade"
	case stream.Quote:
		eventName = "quote"
	case stream.Bar:
		eventName = "bar"
	case *marketdata.Snapshot:
		eventName = "snapshot"
	case *alpaca.Asset:
		eventName = "asset"
	case *clock.Status:
		eventName = "clock-status"
	default:
		panic(fmt.Sprintf("Unknown type: %T", data))
	}

	if eventName == "trade" {
		streamTrade := data.(stream.Trade)

		runtime.EventsEmit(a.ctx, eventName, map[string]interface{}{
			"S": streamTrade.Symbol,
			"p": streamTrade.Price,
			"x": streamTrade.Exchange,
			"s": streamTrade.Size,
			"c": streamTrade.Conditions,
			"t": streamTrade.Timestamp,
			"i": streamTrade.ID,
			"z": streamTrade.Tape,
		})
	} else if eventName == "quote" {
		streamQuote := data.(stream.Quote)

		runtime.EventsEmit(a.ctx, eventName, map[string]interface{}{
			"S":  streamQuote.Symbol,
			"bp": streamQuote.BidPrice,
			"bs": streamQuote.BidSize,
			"ap": streamQuote.AskPrice,
			"as": streamQuote.AskSize,
			"t":  streamQuote.Timestamp,
			"c":  streamQuote.Conditions,
			"z":  streamQuote.Tape,
		})
	} else if eventName == "bar" {
		streamBar := data.(stream.Bar)

		runtime.EventsEmit(a.ctx, eventName, map[string]interface{}{
			"S":  streamBar.Symbol,
			"o":  streamBar.Open,
			"h":  streamBar.High,
			"l":  streamBar.Low,
			"c":  streamBar.Close,
			"v":  streamBar.Volume,
			"t":  streamBar.Timestamp,
			"vw": streamBar.VWAP,
			"n":  streamBar.TradeCount,
		})
	} else {
		runtime.EventsEmit(a.ctx, eventName, data)
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.alpacaClient = alpaca.NewClient(alpaca.ClientOpts{})
	a.marketDataClient = marketdata.NewClient(marketdata.ClientOpts{})

	streamCtx, cancel := context.WithCancel(a.ctx)
	a.streamCtx = streamCtx
	a.cancelStream = cancel

	assetRepository, err := asset.NewRepository(a.db, a.alpacaClient)
	fatal(err)
	a.assetRepository = assetRepository

	calendarRepository, err := calendar.NewRepository(a.db, a.alpacaClient)
	fatal(err)
	a.calendarRepository = calendarRepository

	barRepository, err := bar.NewRepository(a.marketDataClient, a.calendarRepository)
	fatal(err)
	a.barRepository = barRepository

	stockStream, err := stock.NewStream(a.streamCtx, a.trades, a.quotes, a.bars)
	fatal(err)
	a.stockStream = stockStream

	statusClock, err := clock.NewClock(a.ctx, a.status, a.calendarRepository)
	fatal(err)
	a.statusClock = statusClock

	appConfigurationRepository, err := configuration.NewRepository(a.db)
	fatal(err)
	a.appConfigurationRepository = appConfigurationRepository

	isEmpty, err := a.appConfigurationRepository.IsEmpty()
	fatal(err)

	if isEmpty {
		x, y := runtime.WindowGetPosition(ctx)
		width, height := runtime.WindowGetSize(ctx)
		err := a.appConfigurationRepository.UpdateWindow(x, y, width, height)
		fatal(err)
	} else {
		appConfiguration, err := a.appConfigurationRepository.Get()
		fatal(err)

		runtime.WindowSetPosition(ctx, appConfiguration.X, appConfiguration.Y)
		runtime.WindowSetSize(ctx, appConfiguration.Width, appConfiguration.Height)
		runtime.EventsEmit(a.ctx, "ready")
	}

	a.ready = true
}

func (a *App) GetIntradayBars(symbol string) ([]marketdata.Bar, error) {
	return a.barRepository.Intraday(symbol)
}

func (a *App) GetCurrentCalendar() (*calendar.Calendar, error) {
	return a.calendarRepository.CurrentCalendar()
}

func (a *App) GetPrevCalendar() (*calendar.Calendar, error) {
	return a.calendarRepository.PreviousCalendar()
}

func (a *App) GetAssets() ([][]string, error) {
	assets, err := a.assetRepository.GetAll()

	if err != nil {
		return nil, err
	}

	rows := make([][]string, len(assets))

	for i, x := range assets {
		rows[i] = []string{
			x.Symbol,
			x.Name,
			x.Exchange,
		}
	}

	return rows, nil
}

func (a *App) GetAsset(symbol string) (*alpaca.Asset, error) {
	symbolAsset, err := a.assetRepository.Get(symbol)

	if err != nil {
		return nil, err
	}

	return symbolAsset, nil
}

func (a *App) GetSnapshot(symbol string) *marketdata.Snapshot {
	if symbol == "" {
		symbol = a.currentSymbol
	}

	if symbol == "" {
		return nil
	}

	snapshot, err := a.marketDataClient.GetSnapshot(symbol, marketdata.GetSnapshotRequest{})
	fatal(err)

	return snapshot
}

func (a *App) IsReady() bool {
	return a.ready
}

func (a *App) Subscribe(symbol string) error {
	return a.stockStream.SubscribeTo(symbol)
}

func (a *App) Unsubscribe(symbol string) error {
	return a.stockStream.UnsubscribeFrom(symbol)
}

func (a *App) shutdown(ctx context.Context) {
	a.cancelStream()
	x, y := runtime.WindowGetPosition(ctx)
	width, height := runtime.WindowGetSize(ctx)
	err := a.appConfigurationRepository.UpdateWindow(x, y, width, height)
	fatal(err)
}
