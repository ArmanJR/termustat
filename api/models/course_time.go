package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type CourseTime struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	DayOfWeek int        `gorm:"check:day_of_week BETWEEN 0 AND 6"`
	StartTime CustomTime `gorm:"type:time;not null"` // Changed from time.Time
	EndTime   CustomTime `gorm:"type:time;not null"` // Changed from time.Time
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

// CustomTime handles parsing time from string values from the database
type CustomTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface for CustomTime
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time, ct.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time, ct.Valid = v, true
		return nil
	case []byte: // Handle byte slices (often returned by database drivers for strings)
		s := string(v)
		// Try parsing in HH:MM:SS format first
		parsedTime, err := time.Parse("15:04:05", s)
		if err != nil {
			// If HH:MM:SS fails, try HH:MM
			parsedTime, err = time.Parse("15:04", s)
			if err != nil {
				ct.Valid = false
				return fmt.Errorf("CustomTime: unsupported type %T or format for value: %v", value, value)
			}
		}
		ct.Time, ct.Valid = parsedTime, true
		return nil
	case string:
		// Try parsing in HH:MM:SS format first
		parsedTime, err := time.Parse("15:04:05", v)
		if err != nil {
			// If HH:MM:SS fails, try HH:MM
			parsedTime, err = time.Parse("15:04", v)
			if err != nil {
				ct.Valid = false
				return fmt.Errorf("CustomTime: unsupported type %T or format for value: %v", value, value)
			}
		}
		ct.Time, ct.Valid = parsedTime, true
		return nil
	default:
		ct.Valid = false
		return fmt.Errorf("CustomTime: unsupported type %T for value: %v", value, value)
	}
}

// Value implements the driver Valuer interface for CustomTime
func (ct CustomTime) Value() (driver.Value, error) {
	if !ct.Valid {
		return nil, nil
	}
	// Return time in "HH:MM:SS" format, which is standard for SQL TIME type
	return ct.Time.Format("15:04:05"), nil
}

func (CourseTime) TableName() string {
	return "course_times"
}
