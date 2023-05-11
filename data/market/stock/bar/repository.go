package bar

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/buffalo/data/metadata/calendar"
)

// Repository provides access to market data bars information
type Repository struct {
	marketDataClient   *marketdata.Client
	calendarRepository *calendar.Repository
}

func NewRepository(marketDataClient *marketdata.Client, calendarRepository *calendar.Repository) (*Repository, error) {
	return &Repository{
		marketDataClient:   marketDataClient,
		calendarRepository: calendarRepository,
	}, nil
}

// Intraday returns the intraday bars for the given symbol for either the current or previous day
func (b *Repository) Intraday(symbol string) ([]marketdata.Bar, error) {
	currentCalendar, err := b.calendarRepository.CurrentCalendar()

	if err != nil {
		return nil, err
	}

	if currentCalendar == nil {
		currentCalendar, err = b.calendarRepository.PreviousCalendar()

		if err != nil {
			return nil, err
		}
	}

	bars, err := b.marketDataClient.GetBars(symbol, marketdata.GetBarsRequest{
		TimeFrame: marketdata.TimeFrame{
			N:    1,
			Unit: marketdata.Min,
		},
		Adjustment: marketdata.Split,
		Start:      currentCalendar.SessionOpen,
		End:        currentCalendar.SessionClose,
		Feed:       marketdata.SIP,
	})

	if err != nil {
		return nil, err
	}

	return bars, nil
}

// YTD returns the year-to-date bars for the given symbol with a daily interval
func (b *Repository) YTD(symbol string) ([]marketdata.Bar, error) {
	bars, err := b.marketDataClient.GetBars(symbol, marketdata.GetBarsRequest{
		TimeFrame: marketdata.TimeFrame{
			N:    1,
			Unit: marketdata.Day,
		},
		Adjustment: marketdata.Split,
		Start:      carbon.Now().SubYears(1).SubDays(1).ToStdTime(),
		End:        carbon.Now().StartOfDay().ToStdTime(),
		Feed:       marketdata.SIP,
	})

	if err != nil {
		return nil, err
	}

	return bars, nil
}
