package todos

import (
	"time"
)

type Todo struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Done  bool      `json:"done"`
}
