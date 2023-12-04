package model

// Log
type Log struct {
	ID       int64  `json:"id"`
	DataID   int64  `json:"data_id"`
	Table    string `json:"table"`
	Data     string `json:"data"`
	Action   string `json:"action"`
	UserID   int64  `json:"user_id"`
	CreateAt int64  `json:"create_at"`
	ModifyAt int64  `json:"modify_at"`
}

func (Log) TableName() string {
	return "log"
}

func init() {
	registerNewModelFunc("log", func() (interface{}, interface{}) {
		return &Log{}, []Log{}
	})
}
