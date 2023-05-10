package main

import (
	"context"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Trade struct {
	Symbol    string  `json:"S"`
	Price     float64 `json:"p"`
	Size      float64 `json:"s"`
	Timestamp int64   `json:"t"`
}

type Quote struct {
	Symbol    string  `json:"S"`
	BidPrice  float64 `json:"bp"`
	BidSize   float64 `json:"bs"`
	AskPrice  float64 `json:"ap"`
	AskSize   float64 `json:"as"`
	Timestamp int64   `json:"t"`
}

// App struct
type App struct {
	mut                sync.Mutex
	ctx                context.Context
	streamCtx          context.Context
	cancelStream       context.CancelFunc
	db                 *gorm.DB
	alpacaClient       *alpaca.Client
	marketDataClient   *marketdata.Client
	stocksStreamClient *stream.StocksClient
	ready              bool
	trade              Trade
	tradeChannel       chan stream.Trade
	quoteChannel       chan stream.Quote
	currentSymbol      string
}

// NewApp creates a new App application struct
func NewApp() *App {
	db, err := gorm.Open(sqlite.Open("buffalo.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&AppConfiguration{})

	if err != nil {
		log.Fatal(err)
	}

	tradeChannel := make(chan stream.Trade, 10_000)
	quoteChannel := make(chan stream.Quote, 10_000)
	updateTicker := time.NewTicker(100 * time.Millisecond)

	snapshotTicker := time.NewTicker(1 * time.Second)

	app := &App{
		db:           db,
		tradeChannel: tradeChannel,
		quoteChannel: quoteChannel,
	}

	go func(app *App) {
		var lastTrade stream.Trade
		var lastQuote stream.Quote
		for {
			select {
			case lastTrade = <-app.tradeChannel:
			case lastQuote = <-app.quoteChannel:
			case <-updateTicker.C:
				app.Emit(lastTrade)
				app.Emit(lastQuote)
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

	switch data.(type) {
	case stream.Trade:
		trade := data.(stream.Trade)

		runtime.EventsEmit(a.ctx, "trade", Trade{
			Symbol:    trade.Symbol,
			Price:     trade.Price,
			Size:      float64(trade.Size),
			Timestamp: trade.Timestamp.UnixMilli(),
		})
	case stream.Quote:
		quote := data.(stream.Quote)

		runtime.EventsEmit(a.ctx, "quote", Quote{
			Symbol:    quote.Symbol,
			BidPrice:  quote.BidPrice,
			BidSize:   float64(quote.BidSize),
			AskPrice:  quote.AskPrice,
			AskSize:   float64(quote.AskSize),
			Timestamp: quote.Timestamp.UnixMilli(),
		})
	case *marketdata.Snapshot:
		snapshot := data.(*marketdata.Snapshot)
		runtime.EventsEmit(a.ctx, "snapshot", snapshot)
	case *alpaca.Asset:
		asset := data.(*alpaca.Asset)
		runtime.EventsEmit(a.ctx, "asset", asset)
	default:
		panic(fmt.Sprintf("Unknown type: %T", data))
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.alpacaClient = alpaca.NewClient(alpaca.ClientOpts{})
	a.marketDataClient = marketdata.NewClient(marketdata.ClientOpts{})
	a.stocksStreamClient = stream.NewStocksClient(marketdata.SIP)
	streamCtx, cancel := context.WithCancel(a.ctx)
	a.streamCtx = streamCtx
	a.cancelStream = cancel

	err := a.stocksStreamClient.Connect(a.streamCtx)
	fatal(err)

	a.ready = true

	// load last size configuration
	result := a.db.Model(&AppConfiguration{}).FirstOrCreate(&AppConfiguration{
		Key:    "justme",
		X:      0,
		Y:      0,
		Width:  1024,
		Height: 768,
	})
	fatal(result.Error)

	var appConfiguration AppConfiguration

	result.Scan(&appConfiguration)

	runtime.WindowSetPosition(ctx, appConfiguration.X, appConfiguration.Y)
	runtime.WindowSetSize(ctx, appConfiguration.Width, appConfiguration.Height)
	runtime.EventsEmit(a.ctx, "ready")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetAsset(symbol string) *alpaca.Asset {
	if symbol == "" {
		symbol = a.currentSymbol
	}

	if symbol == "" {
		return nil
	}

	asset, err := a.alpacaClient.GetAsset(symbol)
	fatal(err)

	return asset
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

func (a *App) Subscribe(symbol string) bool {
	if !a.ready {
		return false
	}

	if a.currentSymbol == symbol {
		return false
	}

	a.mut.Lock()
	defer a.mut.Unlock()

	var err error

	if len(a.currentSymbol) > 0 {
		err = a.stocksStreamClient.UnsubscribeFromTrades(a.currentSymbol)
		fatal(err)

		err = a.stocksStreamClient.UnsubscribeFromQuotes(a.currentSymbol)
		fatal(err)
	}

	err = a.stocksStreamClient.SubscribeToTrades(func(trade stream.Trade) {
		a.tradeChannel <- trade
	}, symbol)
	fatal(err)

	err = a.stocksStreamClient.SubscribeToQuotes(func(quote stream.Quote) {
		a.quoteChannel <- quote
	}, symbol)

	a.currentSymbol = symbol

	snapshot := a.GetSnapshot(symbol)
	a.Emit(snapshot)

	asset := a.GetAsset(symbol)
	a.Emit(asset)

	return true
}

func (a *App) shutdown(ctx context.Context) {
	a.cancelStream()
	x, y := runtime.WindowGetPosition(ctx)
	width, height := runtime.WindowGetSize(ctx)

	a.db.Model(&AppConfiguration{})

	result := a.db.Model(&AppConfiguration{}).Where(AppConfiguration{
		Key: "justme",
	}).UpdateColumns(&AppConfiguration{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	})

	if result.Error != nil {
		log.Fatal(result.Error)
	}

	log.Println("position saved")
}
