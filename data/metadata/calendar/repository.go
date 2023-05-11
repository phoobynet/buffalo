package calendar

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository struct {
	alpacaClient *alpaca.Client
	db           *gorm.DB
	populated    bool
}

func NewRepository(db *gorm.DB, alpacaClient *alpaca.Client) (*Repository, error) {
	err := db.AutoMigrate(&Calendar{})

	if err != nil {
		return nil, err
	}

	var count int64

	db.Model(&Calendar{}).Count(&count)

	if count == 0 {
		calendarDays, err := alpacaClient.GetCalendar(alpaca.GetCalendarRequest{
			Start: carbon.Now().SubYears(5).ToStdTime(),
			End:   carbon.Now().AddYears(5).ToStdTime(),
		})

		if err != nil {
			return nil, err
		}

		log.Printf("Populating calendar with %d records", len(calendarDays))

		calendars := make([]*Calendar, len(calendarDays))
		var calendar *Calendar

		for i, calendarDay := range calendarDays {
			calendar, err = ToCalendarFromDay(calendarDay)

			if err != nil {
				return nil, err
			}

			calendars[i] = calendar
		}

		err = db.CreateInBatches(calendars, 100).Error

		if err != nil {
			return nil, err
		}

		log.Printf("Populating calendar with %d records...COMPLETED", len(calendarDays))
	}

	return &Repository{
		alpacaClient: alpacaClient,
		db:           db,
	}, nil
}

func (r *Repository) CurrentCalendar() (*Calendar, error) {
	date := time.Now().Format("2006-01-02")

	var calendar = Calendar{
		Date: date,
	}

	result := r.db.First(&calendar)

	if result.Error != nil {
		return nil, result.Error
	}

	return &calendar, nil
}

func (r *Repository) PreviousCalendar() (*Calendar, error) {
	date := time.Now().Format("2006-01-02")

	var calendar Calendar

	result := r.db.
		Model(&Calendar{}).
		Where("date < ?", date).
		Order("date desc").
		First(&calendar)

	if result.Error != nil {
		return nil, result.Error
	}

	return &calendar, nil
}
