package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang-module/carbon/v2"
	"time"
)

type MarketStatus struct {
	CalendarDay *alpaca.CalendarDay
	Now         time.Time
	Status      string
}

func (a *App) marketStatus() MarketStatus {
	if a.calendar == nil {
		return MarketStatus{
			Now:    time.Now(),
			Status: "closed",
		}
	}

	now := carbon.Now()
}
