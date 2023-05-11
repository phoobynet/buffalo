package clock

import (
	"context"
	"github.com/phoobynet/buffalo/data/metadata/calendar"
	"sync"
	"time"
)

type Status struct {
	CurrentTime  time.Time          `json:"currentTime"`
	Calendar     *calendar.Calendar `json:"calendar"`
	IsOpen       bool               `json:"isOpen"`
	IsClosed     bool               `json:"isClosed"`
	IsPreMarket  bool               `json:"isPreMarket"`
	IsPostMarket bool               `json:"isPostMarket"`
	IsTradingDay bool               `json:"isTradingDay"`
}

type Clock struct {
	mut                sync.RWMutex
	calendarRepository *calendar.Repository
	ticker             *time.Ticker
	currentTime        time.Time
	currentDate        string
	currentCalendar    *calendar.Calendar
	currentStatus      *Status
	status             chan Status
}

func NewClock(ctx context.Context, status chan Status, calendarRepository *calendar.Repository) (*Clock, error) {
	ticker := time.NewTicker(time.Second)
	now := time.Now()

	currentCalendar, err := calendarRepository.CurrentCalendar()

	if err != nil {
		return nil, err
	}

	c := &Clock{
		calendarRepository: calendarRepository,
		ticker:             ticker,
		currentTime:        now,
		currentDate:        now.Format("2006-01-02"),
		currentCalendar:    currentCalendar,
		status:             status,
	}

	go func(c *Clock) {
		var nextDate string
		var t time.Time
		for {
			select {
			case t = <-ticker.C:
				c.mut.Lock()
				c.currentTime = t
				nextDate = t.Format("2006-01-02")
				if c.currentDate != nextDate {
					c.currentDate = nextDate
					c.currentCalendar, err = calendarRepository.CurrentCalendar()

					if err != nil {
						panic(err)
					}
				}

				currentStatus, err := c.CurrentStatus()

				if err != nil {
					panic(err)
				}

				c.status <- *currentStatus
				c.currentStatus = currentStatus

				c.mut.Unlock()
			case <-ctx.Done():
				ticker.Stop()
			}
		}
	}(c)

	return c, nil
}

func (c *Clock) CurrentStatus() (*Status, error) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	status := &Status{
		CurrentTime:  c.currentTime,
		Calendar:     c.currentCalendar,
		IsOpen:       c.IsOpen(),
		IsClosed:     c.IsClosed(),
		IsPreMarket:  c.IsPreMarket(),
		IsPostMarket: c.IsPostMarket(),
		IsTradingDay: c.IsTradingDay(),
	}

	return status, nil
}

// IsOpen returns true if the market is open (inc. pre- and post-market), false otherwise.
func (c *Clock) IsOpen() bool {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if c.currentCalendar == nil {
		return false
	}

	if c.currentTime.Equal(c.currentCalendar.SessionClose) || c.currentTime.After(c.currentCalendar.SessionClose) {
		return false
	}

	if c.currentTime.Before(c.currentCalendar.SessionOpen) {
		return false
	}

	return true
}

// IsClosed returns true if the market is closed for the day, not opened yet, or not a trading day.
func (c *Clock) IsClosed() bool {
	return !c.IsOpen()
}

// IsTradingDay returns true if today is a trading day, false otherwise.
func (c *Clock) IsTradingDay() bool {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.currentCalendar != nil
}

func (c *Clock) IsPreMarket() bool {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if c.currentCalendar == nil {
		return false
	}

	return (c.currentTime.Equal(c.currentCalendar.SessionOpen) || c.currentTime.After(c.currentCalendar.SessionOpen)) &&
		c.currentTime.Before(c.currentCalendar.Open)
}

func (c *Clock) IsPostMarket() bool {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if c.currentCalendar == nil {
		return false
	}

	return c.currentTime.Equal(c.currentCalendar.SessionClose) || c.currentTime.After(c.currentCalendar.SessionClose)
}
