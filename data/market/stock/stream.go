package stock

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"sync"
)

type Stream struct {
	mut          sync.Mutex
	trades       chan stream.Trade
	quotes       chan stream.Quote
	bars         chan stream.Bar
	stocksClient *stream.StocksClient
}

func NewStream(ctx context.Context, trades chan stream.Trade, quotes chan stream.Quote, bars chan stream.Bar) (*Stream, error) {
	stocksClient := stream.NewStocksClient(marketdata.SIP)

	err := stocksClient.Connect(ctx)

	if err != nil {
		return nil, err
	}

	return &Stream{stocksClient: stocksClient, trades: trades, quotes: quotes, bars: bars}, nil
}

func (s *Stream) SubscribeTo(symbols ...string) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	err := s.stocksClient.SubscribeToTrades(func(trade stream.Trade) {
		s.trades <- trade
	}, symbols...)

	if err != nil {
		return err
	}

	err = s.stocksClient.SubscribeToQuotes(func(quote stream.Quote) {
		s.quotes <- quote
	}, symbols...)

	if err != nil {
		return err
	}

	return s.stocksClient.SubscribeToBars(func(bar stream.Bar) {
		s.bars <- bar
	}, symbols...)
}

func (s *Stream) UnsubscribeFrom(symbols ...string) error {
	err := s.stocksClient.UnsubscribeFromTrades(symbols...)

	if err != nil {
		return err
	}

	err = s.stocksClient.UnsubscribeFromQuotes(symbols...)

	if err != nil {
		return err
	}

	return s.stocksClient.UnsubscribeFromBars(symbols...)
}
