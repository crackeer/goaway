package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

const (
	actionSelect   string = "select"
	actionCount    string = "count"
	actionDistinct string = "distinct"

	actionInsert string = "insert"
	actionUpdate string = "update"
	actionDelete string = "delete"

	actionDrop string = "drop"
	actionExec string = "exec"
)

// Request
type Request struct {
	Action string
	Table  string
	Body   []byte

	DB     *gorm.DB
	driver string
}

// SelectBody
type SelectBody struct {
	Query    map[string]interface{} `json:"query"`
	Fields   []string               `json:"fields"`
	OrderBy  string                 `json:"order_by"`
	Limit    int                    `json:"limit"`
	Offset   int                    `json:"offset"`
	Distinct []string               `json:"distinct"`
}

// UpdateBody
type UpdateBody struct {
	Query  map[string]interface{} `json:"query"`
	Update map[string]interface{} `json:"update"`
}

// ParseRequest
//
//	@param r
//	@param path
//	@return *Request
//	@return error
func ParseRequest(r *http.Request, path string) (*Request, error) {
	if len(path) < 1 {
		return nil, errors.New("nil path")
	}
	parts := strings.Split(strings.TrimLeft(path, "/"), "/")
	action := parts[0]
	bytes, _ := io.ReadAll(r.Body)

	return &Request{
		Action: action,
		Table:  parts[1],
		Body:   bytes,
	}, nil
}

// UseDB
//
//	@receiver req
//	@param db
//	@return *Request
func (req *Request) UseDB(db *gorm.DB, driver string) *Request {
	req.DB = db
	req.driver = driver
	return req
}

// IsSQLite
//
//	@receiver req
//	@return bool
func (req *Request) IsSQLite() bool {
	return req.driver == "sqlite"
}

// Handle
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Handle() (interface{}, error) {
	if req.Action == actionSelect {
		return req.Select()
	}

	if req.Action == actionCount {
		return req.Count()
	}

	if req.Action == actionDistinct {
		return req.Distinct()
	}

	if req.Action == actionExec {
		return req.Exec()
	}

	if req.Action == actionInsert {
		return req.Insert()
	}

	if req.Action == actionUpdate {
		return req.Update()
	}

	if req.Action == actionDelete {
		return req.Delete()
	}

	if req.Action == actionDrop {
		return req.Drop()
	}

	return nil, errors.New("no action match")
}

func (req *Request) decodeBody(dest interface{}) error {
	return json.Unmarshal(req.Body, dest)
}

// Select
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Select() (interface{}, error) {
	selectBody := &SelectBody{}
	if err := req.decodeBody(selectBody); err != nil {
		return nil, fmt.Errorf("decode select body error:%s", err.Error())
	}

	db := req.DB.Table(req.Table)
	if len(selectBody.Fields) > 0 {
		db = db.Select(selectBody.Fields)
	}

	if len(selectBody.Query) > 0 {
		sql, params := BuildQuery(selectBody.Query)
		db = db.Where(sql, params...)
	}

	if len(selectBody.OrderBy) > 0 {
		db = db.Order(selectBody.OrderBy)
	}

	if selectBody.Offset > 0 {
		db = db.Offset(selectBody.Offset)
	}
	if selectBody.Limit > 0 {
		db = db.Limit(selectBody.Limit)
	}

	list := []map[string]interface{}{}

	if err := db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("select error:%s", err.Error())
	}

	return list, nil
}

// Select
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Count() (interface{}, error) {
	selectBody := &SelectBody{}
	if err := req.decodeBody(selectBody); err != nil {
		return nil, fmt.Errorf("decode select body error:%s", err.Error())
	}

	db := req.DB.Table(req.Table)
	if len(selectBody.Query) > 0 {
		sql, params := BuildQuery(selectBody.Query)
		db = db.Where(sql, params...)
	}

	var count int64
	if len(selectBody.Distinct) > 0 {
		db = db.Distinct(selectBody.Distinct[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("count error:%s", err.Error())
	}

	return count, nil
}

// Distinct
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Distinct() (interface{}, error) {
	selectBody := &SelectBody{}
	if err := req.decodeBody(selectBody); err != nil {
		return nil, fmt.Errorf("decode select body error:%s", err.Error())
	}

	if len(selectBody.Distinct) < 1 {
		return nil, fmt.Errorf("distinct colum nil")
	}

	db := req.DB.Table(req.Table)
	if len(selectBody.Query) > 0 {
		sql, params := BuildQuery(selectBody.Query)
		db = db.Where(sql, params...)
	}

	db = db.Distinct(selectBody.Distinct[0])

	list := []string{}

	if err := db.Pluck(selectBody.Distinct[0], &list).Error; err != nil {
		return nil, fmt.Errorf("count error:%s", err.Error())
	}

	return list, nil
}

// Desc
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Desc() (interface{}, error) {
	list := []map[string]interface{}{}
	if req.IsSQLite() {
		if err := req.DB.Raw(fmt.Sprintf("PRAGMA TABLE_INFO (%s)", req.Table)).Scan(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	}

	if err := req.DB.Raw("desc " + req.Table).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// Exec
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Exec() (interface{}, error) {
	if err := req.DB.Exec(string(req.Body)).Error; err != nil {
		return nil, err
	}
	return "ok", nil
}

// Insert
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Insert() (interface{}, error) {
	list := []map[string]interface{}{}
	if err := req.decodeBody(&list); err != nil {
		return nil, fmt.Errorf("decode insert body error:%s", err.Error())
	}

	if err := req.DB.Table(req.Table).Create(&list).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

// Update
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Update() (interface{}, error) {
	updateBody := &UpdateBody{}
	if err := req.decodeBody(updateBody); err != nil {
		return nil, fmt.Errorf("decode update body error:%s", err.Error())
	}

	sql, params := BuildQuery(updateBody.Query)

	db := req.DB.Table(req.Table).Where(sql, params...).Updates(updateBody.Update)
	if err := db.Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"affected_rows": db.RowsAffected,
	}, nil
}

func (req *Request) Delete() (interface{}, error) {
	updateBody := &UpdateBody{}
	if err := req.decodeBody(updateBody); err != nil {
		return nil, fmt.Errorf("decode update body error:%s", err.Error())
	}

	where, params := BuildQuery(updateBody.Query)
	sql := fmt.Sprintf("delete from %s where %s", req.Table, where)
	db := req.DB.Exec(sql, params...)
	if err := db.Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"affected_rows": db.RowsAffected,
	}, nil
}

// Drop
//
//	@receiver req
//	@return interface{}
//	@return error
func (req *Request) Drop() (interface{}, error) {
	db := req.DB.Exec("drop table " + req.Table)
	if err := db.Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"affected_rows": db.RowsAffected,
	}, nil
}

// BuildQuery
//
//	@param query
//	@return string
//	@return []interface{}
func BuildQuery(query map[string]interface{}) (string, []interface{}) {
	queryConditions := []string{}
	params := []interface{}{}
	for key, val := range query {

		if strings.HasPrefix(key, "like@") || strings.HasPrefix(key, "plike@") {
		} else {
			params = append(params, val)
		}

		if !strings.Contains(key, "@") {
			queryConditions = append(queryConditions, fmt.Sprintf("%s in (?)", key))
			continue

		}
		parts := strings.Split(key, "@")
		if len(parts) < 2 {
			queryConditions = append(queryConditions, fmt.Sprintf("%s in (?)", key))
			continue
		}

		switch parts[0] {
		case "gt":
			queryConditions = append(queryConditions, fmt.Sprintf("%s > ?", parts[1]))
		case "gte":
			queryConditions = append(queryConditions, fmt.Sprintf("%s >= ?", parts[1]))
		case "lt":
			queryConditions = append(queryConditions, fmt.Sprintf("%s < ?", parts[1]))
		case "lte":
			queryConditions = append(queryConditions, fmt.Sprintf("%s <= ?", parts[1]))
		case "like":
			queryConditions = append(queryConditions, fmt.Sprintf("%s like '%%%v%%'", parts[1], val))
		case "plike":
			queryConditions = append(queryConditions, fmt.Sprintf("%s like '%v%%'", parts[1], val))
		default:
			queryConditions = append(queryConditions, fmt.Sprintf("%s in (?)", key))
		}
	}
	return strings.Join(queryConditions, " and "), params
}
