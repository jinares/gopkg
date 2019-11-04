package xmysql

import (
	"errors"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
)

//SyncItem SyncItem
type RuleItem struct {
	Type   int                 `yaml:"Type" json:"Type"` //0 #前缀匹配 ==1全等匹配
	DB     string              `yaml:"DB" json:"DB"`
	Table  string              `yaml:"Table" json:"Table"`
	Params map[string][]string `yaml:"Params" json:"Params"`
}

//UpdatePosFunc
type UpdatePosFunc func(pos Position) error

//CanalConfig CanalConfig
type CanalConfig struct {
	Addr     string   `yaml:"Addr" json:"Addr"`
	User     string   `yaml:"User" json:"User"`
	Password string   `yaml:"Password" json:"Password"`
	ServerId uint32   `yaml:"ServerId" json:"ServerId"`
	Position Position `yaml:"Position" json:"Position"`
}

type SyncItemConfig struct {
	Name   string      `yaml:"Name" json:"Name"`
	Mysql  CanalConfig `yaml:"Mysql"`
	Runing int         `yaml:"Runing"` //==0 runing ==1 stop
	Rules  []RuleItem  `yaml:"Rules"`
}

//Position Position
type Position struct {
	Name string `yaml:"Name" json:"Name"`
	Pos  uint32 `yaml:"Pos" json:"Pos"`
}

//HandlerFunc func(c *gin.Context) IRet
func newSyncEventHandler(conf *SyncItemConfig) *SyncEventHandler {

	return &SyncEventHandler{
		synced: true,
		conf:   conf,
	}
}

type SyncEventHandler struct {
	canal.DummyEventHandler

	synced    bool
	conf      *SyncItemConfig
	UpdataPos UpdatePosFunc
	SyncData  SyncDataHandler
}

func (h *SyncEventHandler) String() string { return "SyncEventHandler" }

func (h *SyncEventHandler) OnPosSynced(pos mysql.Position, gtid mysql.GTIDSet, ret bool) error {
	//fmt.Println("OnPosSynced",pos,ret)
	if h.synced {
		if h.UpdataPos != nil {
			return h.UpdataPos(Position{Name: pos.Name, Pos: pos.Pos})
		}
	}
	return nil
}

func (h *SyncEventHandler) OnRow(e *canal.RowsEvent) error {
	log := ulog.WithField("mysql", "canal")
	if h.conf == nil {
		log.Error("SyncEventHandler config err")
		return nil
	}
	if ftype := IsDML(e.Action); ftype > 0 {
		if h.SyncData != nil {
			data, err := getRowMap(e)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
			return h.SyncData(data)
		}
	}
	return nil
}

//MysqlCanal MysqlCanal
func MysqlCanal(
	opt *SyncItemConfig, updatePosHandler UpdatePosFunc,
	syncHandler SyncDataHandler, pos *Position,
) error {
	if opt == nil {
		return errors.New("conf err")
	}
	cfg := canal.NewDefaultConfig()
	cfg.Addr = opt.Mysql.Addr
	cfg.User = opt.Mysql.User
	cfg.Password = opt.Mysql.Password
	cfg.ServerID = opt.Mysql.ServerId

	cfg.Dump.TableDB = ""
	cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)

	if err != nil {
		return err

	}
	handler := newSyncEventHandler(opt)
	handler.UpdataPos = updatePosHandler
	handler.SyncData = syncHandler

	c.SetEventHandler(handler)

	if pos == nil || pos.Name == "" {
		pos = &opt.Mysql.Position
	}
	mpos := mysql.Position{
		Name: pos.Name,
		Pos:  pos.Pos,
	}
	if err := CheckPosition(c, &mpos); err != nil {
		return err
	}

	// Start
	defer c.Close()
	//c.Dump()
	return c.RunFrom(mpos)
}

//CheckPosition 检查日志是否存在　否则返回当前最新最新位置
func CheckPosition(c *canal.Canal, pos *mysql.Position) error {
	if pos.Name == "" {
		tmp, err := c.GetMasterPos()
		if err != nil {
			return nil
		}
		pos.Name = tmp.Name
		pos.Pos = tmp.Pos
		return nil
	}
	rr, err := c.Execute("SHOW MASTER LOGS;")
	if err != nil {
		return err
	}
	rownum := rr.RowNumber()

	for index := 0; index < rownum; index++ {
		val, err := rr.GetStringByName(index, "Log_name")
		if err != nil {
			continue
		}
		if val == pos.Name {
			return nil
		}
	}
	pos.Name, _ = rr.GetStringByName(0, "Log_name")
	pos.Pos = 0
	return nil
}
