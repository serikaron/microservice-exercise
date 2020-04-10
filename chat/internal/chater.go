package internal

import (
	"fmt"
)

func Chat(chater string, msg string) string {
	return fmt.Sprintf("%s: %s", chater, msg)
}
