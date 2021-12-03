package model

type FilterLabel struct {
	Id        int64  `xorm:"pk autoincr id" json:"id"`
	LableName string `xorm:"varchar(50) name" json:"lable_name"`
}
