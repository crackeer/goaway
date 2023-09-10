package container

import (
	"github.com/robfig/cron/v3"
)

var cronTab *cron.Cron

// StartSchedule
//
//	@param ctx
//	@param cfg
func Schedule() {
	cronTab = cron.New(cron.WithSeconds())
	cronTab.AddFunc("1 * * * * *", saveModel)
	cronTab.Start()
}
