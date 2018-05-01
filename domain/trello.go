package domain

import (
	"time"
)

type TrelloCardDoingDuration struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Duration        int64     `json:"duration"`
	Date            time.Time `json:"date"`
	YoutrackSummary string    `json:"youtrack_summary"` // ovo staviti negdje drugdje
}
