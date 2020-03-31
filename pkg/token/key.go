package token

import (
	"time"
)

type Key interface {
	ID() string
	Read() ([]byte, error)
	ExpireAt() *time.Time
}
