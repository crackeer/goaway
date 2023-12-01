package model

// Log
type Log struct {
	ID       int64  `json:"id"`
	DataID   int64  `json:"data_id"`
	Table    string `json:"table"`
	Data     string `json:"data"`
	Action   string `json:"action"`
	CreateAt int64  `json:"create_at"`
	ModifyAt int64  `json:"modify_at"`
}

func (Log) TableName() string {
	return "log"
}
