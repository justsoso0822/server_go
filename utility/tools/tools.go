package tools

import (
	"math"
	"strconv"

	"server_go/internal/consts"
)

func ParseRes(items interface{}) []consts.ResItem {
	switch v := items.(type) {
	case []consts.ResItem:
		return v
	case string:
		return parseResString(v)
	default:
		return nil
	}
}

func parseResString(s string) []consts.ResItem {
	nums := PickNumbers(s)
	if len(nums) == 0 || len(nums)%3 != 0 {
		return nil
	}
	result := make([]consts.ResItem, 0, len(nums)/3)
	for i := 0; i < len(nums); i += 3 {
		result = append(result, consts.ResItem{Type: nums[i], Id: nums[i+1], Cnt: nums[i+2]})
	}
	return result
}

// PickNumbers 从字符串中提取所有整数。
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
