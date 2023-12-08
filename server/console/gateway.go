package console

import (
	"fmt"
	"strconv"
	"time"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
)

func getTable(ctx *gin.Context) string {
	return ctx.Param("table")
}

func getDataID(ctx *gin.Context) int64 {
	id := ctx.Param("id")
	value, _ := strconv.Atoi(id)
	return int64(value)
}

// Create
//
//	@param ctx
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
	if err := model.CheckModelCreate(db, value); err != nil {
		ginHelper.Failure(ctx, -2, err.Error())
		return
	}
	result := db.Create(value)
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
		return
	}
	dataID := extractID(value)
	ctx.Set("data_id", dataID)
	user := GetCurrentUser(ctx)
	db.Table(getTable(ctx)).Where(map[string]interface{}{"id": dataID}).Updates(map[string]interface{}{
		"create_at": time.Now().Unix(),
		"modify_at": time.Now().Unix(),
		"user_id":   user.ID,
	})
	ginHelper.Success(ctx, value)
}

// Modify
//
//	@param ctx
func Modify(ctx *gin.Context) {
	if dataID := getDataID(ctx); dataID < 1 {
		ginHelper.Failure(ctx, -1, "data id = 0")
		return
	}
	db := container.GetModelDB()
	updateData := ginHelper.AllPostParams(ctx)
	if newUpdateData, err := model.MakeModifyData(getTable(ctx), updateData); err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	} else {
		updateData = newUpdateData
	}
	updateData["modify_at"] = time.Now().Unix()
	result := db.Table(getTable(ctx)).Where(map[string]interface{}{"id": getDataID(ctx)}).Updates(updateData)
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
	} else {
		ginHelper.Success(ctx, map[string]interface{}{
			"affected": result.RowsAffected,
		})
	}
}

// Delete
//
//	@param ctx
func Delete(ctx *gin.Context) {
	if dataID := getDataID(ctx); dataID < 1 {
		ginHelper.Failure(ctx, -1, "data id = 0")
		return
	}
	db := container.GetModelDB()
	result := db.Exec(fmt.Sprintf("DELETE FROM %s where id = %d", getTable(ctx), getDataID(ctx)))
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
	} else {
		ginHelper.Success(ctx, map[string]interface{}{
			"affected": result.RowsAffected,
		})
	}
}

// Query
//
//	@param ctx
func Query(ctx *gin.Context) {
	var (
		list []map[string]interface{}
	)
	db := container.GetModelDB()
	query := ginHelper.AllGetParams(ctx)

	db.Table(getTable(ctx)).Where(query).Order("id desc").Limit(1000).Find(&list)
	ginHelper.Success(ctx, list)
}
