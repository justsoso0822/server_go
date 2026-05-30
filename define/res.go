package define

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
