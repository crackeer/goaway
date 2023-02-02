package container

import (
	"git.lianjia.com/vrlab_server/gopkg/guid"
	"github.com/sony/sonyflake"
)

var guidInstance *sonyflake.Sonyflake

// InitGUID ...
func InitGUID() {
	var err error
	guidInstance, err = guid.New(nil, "2022-01-05")

	if err != nil {
		panic("init guid failed with [" + err.Error() + "]")
	}

	if guidInstance == nil {
		panic("init guid failed with [empty guidInstance]")
	}
}

// NextID ...
func NextID() int64 {
	uint64, _ := guidInstance.NextID()
	return int64(uint64)
}
