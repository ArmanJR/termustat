package models

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestCustomTime_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantCt  CustomTime
		wantErr bool
	}{
		{
			name:  "scan nil",
			value: nil,
			wantCt: CustomTime{Time: time.Time{}, Valid: false},
			wantErr: false,
		},
		{
			name:  "scan time.Time",
			value: time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
			wantCt: CustomTime{Time: time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC), Valid: true},
			wantErr: false,
		},
		{
			name:  "scan string HH:MM:SS",
			value: "14:45:30",
			wantCt: CustomTime{Time: time.Date(0, 1, 1, 14, 45, 30, 0, time.UTC), Valid: true},
			wantErr: false,
		},
		{
			name:  "scan string HH:MM",
			value: "09:15",
			wantCt: CustomTime{Time: time.Date(0, 1, 1, 9, 15, 0, 0, time.UTC), Valid: true},
			wantErr: false,
		},
		{
			name:    "scan invalid string",
			value:   "invalid-time",
			wantCt:  CustomTime{Time: time.Time{}, Valid: false},
			wantErr: true,
		},
		{
			name:    "scan integer (unsupported)",
			value:   12345,
			wantCt:  CustomTime{Time: time.Time{}, Valid: false},
			wantErr: true,
		},
		{
			name:  "scan byte slice HH:MM:SS",
			value: []byte("16:30:00"),
			wantCt: CustomTime{Time: time.Date(0, 1, 1, 16, 30, 0, 0, time.UTC), Valid: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime
			err := ct.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomTime.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// For time comparison, ignore the location if not explicitly set in wantCt.Time for non-nil cases
			if tt.wantCt.Valid && !ct.Time.Equal(tt.wantCt.Time) {
                 // Re-parse to ensure consistent location for comparison if needed, or use specific parts
                 expectedParsed, _ := time.Parse("15:04:05", tt.wantCt.Time.Format("15:04:05"))
                 if !ct.Time.Equal(expectedParsed) {
                    t.Errorf("CustomTime.Scan() gotTime = %v, want %v", ct.Time, tt.wantCt.Time)
                 }
			}
            if ct.Valid != tt.wantCt.Valid {
                 t.Errorf("CustomTime.Scan() gotValid = %v, want %v", ct.Valid, tt.wantCt.Valid)
            }
		})
	}
}

func TestCustomTime_Value(t *testing.T) {
	tests := []struct {
		name    string
		ct      CustomTime
		wantVal driver.Value
		wantErr bool
	}{
		{
			name:    "value valid time",
			ct:      CustomTime{Time: time.Date(0, 1, 1, 10, 30, 15, 0, time.UTC), Valid: true},
			wantVal: "10:30:15",
			wantErr: false,
		},
		{
			name:    "value zero time but valid", // e.g. midnight
			ct:      CustomTime{Time: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			wantVal: "00:00:00",
			wantErr: false,
		},
		{
			name:    "value invalid time",
			ct:      CustomTime{Time: time.Time{}, Valid: false},
			wantVal: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.ct.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomTime.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("CustomTime.Value() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
