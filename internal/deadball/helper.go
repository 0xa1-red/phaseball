package deadball

import (
	"fmt"
	"strconv"
)

func LastDigit(num int) int {
	s := fmt.Sprintf("%d", num)
	digit := s[len(s)-1:]

	if i, err := strconv.Atoi(digit); err != nil {
		return -1
	} else {
		return i
	}
}
