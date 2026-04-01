package domain

// Competition is the master record of a national competition (one row per competition).
type Competition struct {
	ID                int64  `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Level             string `json:"level" db:"level"`
	OfficialURL       string `json:"official_url" db:"official_url"`
	NameKey           string `json:"-" db:"name_key"`
	TypicalTimeWindow string `json:"typical_time_window" db:"typical_time_window"`
	TimelineHint      string `json:"timeline_hint" db:"timeline_hint"`
	Note              string `json:"note" db:"note"`
}

