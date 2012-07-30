package util

import (
	"fmt"
)

func Address(hostname string, port int) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}
