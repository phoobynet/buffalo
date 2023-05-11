package calendar

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang-module/carbon/v2"
	"strings"
	"time"
)

type Calendar struct {
	Date         string    `json:"date" gorm:"primaryKey"`
	Open         time.Time `json:"open"`
	Close        time.Time `json:"close"`
	SessionOpen  time.Time `json:"sessionOpen"`
	SessionClose time.Time `json:"sessionClose"`
}

// parseTime - parses a time string returned from Alpaca CalendarDay query which may contain a colon
func parseTime(t string) string {
	if strings.Contains(t, ":") {
		return t
	}

	return fmt.Sprintf("%s:%s", t[:2], t[2:])
}

func toTime(date, time string) (*time.Time, error) {
	result := carbon.Parse(fmt.Sprintf("%s %s", date, parseTime(time)), carbon.NewYork)

	if result.Error != nil {
		return nil, result.Error
	}

	t := result.ToStdTime()

	return &t, nil
}

func ToCalendarFromDay(calendarDay alpaca.CalendarDay) (*Calendar, error) {
	calendar := Calendar{}

	if openingTime, err := toTime(calendarDay.Date, calendarDay.Open); err == nil {
		calendar.Open = *openingTime
	} else {
		return nil, err
	}

	if closingTime, err := toTime(calendarDay.Date, calendarDay.Close); err == nil {
		calendar.Close = *closingTime
	} else {
		return nil, err
	}

	if sessionOpeningTime, err := toTime(calendarDay.Date, "04:00"); err == nil {
		calendar.SessionOpen = *sessionOpeningTime
	} else {
		return nil, err
	}

	if sessionClosingTime, err := toTime(calendarDay.Date, "20:00"); err == nil {
		calendar.SessionClose = *sessionClosingTime
	} else {
		return nil, err
	}

	return &calendar, nil
}
