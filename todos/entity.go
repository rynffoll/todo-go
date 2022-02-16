package todos

import (
	"time"
)

type Todo struct {
	ID    int       `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
	Date  time.Time `json:"date" db:"date"`
	Done  bool      `json:"done" db:"done"`
}
