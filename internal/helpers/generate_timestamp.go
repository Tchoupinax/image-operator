package helpers

import "time"

func GenerateTimestamp() float64 {
	return float64(time.Now().Unix())
}
