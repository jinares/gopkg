package xmysql

import (
	"fmt"
	"testing"

	"github.com/jinares/gopkg/xtools"
)

func TestMysqlCanal(t *testing.T) {
	fmt.Println("start ")
	err := MysqlCanal(&SyncItemConfig{
		Name: "",
		Mysql: CanalConfig{
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "senseye3",
			ServerId: 201,
			Position: Position{Name: "mysql-bin.000001", Pos: 0},
		},
		Runing: 0,
		Rules:  nil,
	}, nil, func(row RowData) error {
		fmt.Println(row, xtools.JSONMarshal(row))
		return nil
	}, nil)
	fmt.Println(err)
}
