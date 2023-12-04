package console

import (
	ginHelper "github.com/crackeer/gopkg/gin"
	api "github.com/crackeer/simple_http"
	"github.com/gin-gonic/gin"
)

func Create(ctx *gin.Context) {
	db := container.GetModelDB()
	var (
		table string = getTable(ctx)
		value interface{}
	)
	value, _ = model.NewModel(table)
	if err := ctx.ShouldBindJSON(value); err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	result := db.Create(value)
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
		return
	}
	dataID := extractID(value)
	ctx.Set("data_id", dataID)
	user := getCurrentUser(ctx)
	db.Table(getTable(ctx)).Where(map[string]interface{}{"id": dataID}).Updates(map[string]interface{}{
		"create_at": time.Now().Unix(),
		"modify_at": time.Now().Unix(),
		"user_id" : user.ID,
	})
	ginHelper.Success(ctx, value)
}