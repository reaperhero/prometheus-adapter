package model

type CaptrueMetric struct {
	Id       int64  `xorm:"pk autoincr id" json:"id"`
	CapName  string `xorm:"varchar(50) cap_name" json:"cap_name"`
	CapSql   string `xorm:"varchar(255) cap_sql" json:"cap_sql"`
	Status   bool   `xorm:"status" json:"status"`
	Instance string `xorm:"instance" json:"instance"`
}
