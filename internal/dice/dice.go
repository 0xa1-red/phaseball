package dice

import (
	"math/rand"
	"time"
)

func Roll(sides, times uint8, mod int) int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	total := 0
	var i uint8
	for i = 0; i < times; i++ {
		total += rnd.Intn(int(sides)) + 1
	}
	return total + mod
}
