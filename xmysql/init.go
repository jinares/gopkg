package xmysql

import (
	"errors"
	//"fmt"

	"github.com/jinares/gopkg/xlog"
	//"github.com/jinares/gopkg/xtools"
	"github.com/siddontang/go-mysql/canal"
)

var (
	ulog = xlog.GetLog()
)

const (
	DeleteAction string = canal.DeleteAction
	InsertAction string = canal.InsertAction
	UpdateAction string = canal.UpdateAction
)

//IsDML IsDML
func IsDML(action string) int {
	data := map[string]int{

		canal.DeleteAction: 1,
		canal.InsertAction: 2,
		canal.UpdateAction: 3,
	}
	if val, isok := data[action]; isok {
		return val
	}
	return -1
}

//GetRowMap GetRowMap
func getRowMap(event *canal.RowsEvent) (RowData, error) {
	coll := event.Table.Columns
	rows := event.Rows
	dbname := event.Table.Schema
	tablename := event.Table.Name
	data := make([]map[string]interface{}, 0)
	if len(coll) < 1 || len(rows) < 1 {
		return RowData{}, errors.New("empty")
	}
	for _, row := range rows {
		item := map[string]interface{}{}
		for i, v := range row {
			if len(coll) > i {
				val := coll[i]
				if val.RawType == "text" {
					item[val.Name] = string(v.([]byte))
				} else {
					item[val.Name] = v
				}

			}

		}
		data = append(data, item)
	}

	return RowData{
		DBName:    dbname,
		TableName: tablename,
		Data:      data,
		Action:    event.Action,
	}, nil
}
