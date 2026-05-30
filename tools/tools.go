package tools

import (
	"math"
	"strconv"
)

type ResItem struct {
	Type int `json:"type"`
	Id   int `json:"id"`
	Cnt  int `json:"cnt"`
}

const (
	ResTypeDiamond   = 1
	ResTypeGold      = 2
	ResTypeTili      = 3
	ResTypeExp       = 4
	ResTypeStar      = 5
	ResTypeItemOther = 6
)

func PickNumbers(s string) []int {
	var result []int
	current := ""
	for _, c := range s {
		if (c >= '0' && c <= '9') || c == '-' || c == '+' || c == '.' {
			current += string(c)
		} else if current != "" {
			appendNumber(&result, current)
			current = ""
		}
	}
	if current != "" {
		appendNumber(&result, current)
	}
	return result
}

func appendNumber(result *[]int, current string) {
	if n, err := strconv.Atoi(current); err == nil {
		*result = append(*result, n)
	} else if f, err := strconv.ParseFloat(current, 64); err == nil {
		*result = append(*result, int(math.Floor(f)))
	}
}
