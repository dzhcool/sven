package models

import (
	"encoding/gob"
	"fmt"
	"log"
	"time"
	"website/model/base"

	"github.com/dzhcool/sven/setting"
	"github.com/jinzhu/gorm"
)

// 数据库链接
var db *gorm.DB

// Ping 测试数据库连接
func Ping() {
	db.DB().Ping()
}

func InitDB() {
	var err error
	config := base.LoadConfigs("api")
	if db, err = base.NewEngine(config); err != nil {
		log.Fatalf("[orm] error: %v\n", err)
	}
	db.DB().SetMaxIdleConns(10)
	// 开启调试
	if setting.AppDebug {
		db.LogMode(true)
	}

	// 同步MySQL结构
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET="+config.Charset).AutoMigrate()

	// 注册gob类型，for 缓存
	gob.Register(time.Time{})
}

func errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf(format)
}

//数据库结构基础字段
type DBaseTime struct {
	Addtime int64  `json:"addtime" gorm:"type:int(10) unsigned NOT NULL DEFAULT '0';"`
	Addate  string `json:"addate" gorm:"-"`
	Uptime  int64  `json:"uptime" gorm:"type:int(10) unsigned NOT NULL DEFAULT '0';"`
	Update  string `json:"update" gorm:"-"`
}

//创建数据自动插入添加时间
func (p *DBaseTime) BeforeCreate(scope *gorm.Scope) (err error) {
	p.Addtime = time.Now().Unix()
	p.Uptime = p.Addtime
	return nil
}

//更新数据自动插入更新时间
func (p *DBaseTime) BeforeSave(scope *gorm.Scope) (err error) {
	uptime := time.Now().Unix()
	scope.SetColumn("uptime", uptime)
	return nil
}

func (p *DBaseTime) BeforeUpdate(scope *gorm.Scope) (err error) {
	uptime := time.Now().Unix()
	scope.SetColumn("uptime", uptime)
	return nil
}

type DBase struct {
	Id int `json:"id" gorm:"primary_key; type:int(11) unsigned NOT NULL AUTO_INCREMENT;"`
	DBaseTime
}

// 新结构
type GBaseTime struct {
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAtHuman string    `json:"created_at_human" gorm:"-"`
	UpdatedAtHuman string    `json:"updated_at_human" gorm:"-"`
}

//创建数据自动插入添加时间
func (p *GBaseTime) BeforeCreate(scope *gorm.Scope) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	return nil
}

//更新数据自动插入更新时间
func (p *GBaseTime) BeforeSave(scope *gorm.Scope) (err error) {
	uptime := time.Now()
	scope.SetColumn("updated_at", uptime)
	return nil
}

func (p *GBaseTime) BeforeUpdate(scope *gorm.Scope) (err error) {
	uptime := time.Now()
	scope.SetColumn("updated_at", uptime)
	return nil
}

func (p *GBaseTime) FormatTime() {
	p.CreatedAtHuman = p.CreatedAt.Format("2006-01-02 15:04:05")
	p.UpdatedAtHuman = p.UpdatedAt.Format("2006-01-02 15:04:05")
}

type GBase struct {
	ID int64 `json:"id" gorm:"column:id" form:"id"`
	GBaseTime
}

//模型基础配置
type model struct {
}

//分页参数处理
func pagecut(page, pagesize int) (int, int) {
	if pagesize <= 0 {
		return 0, -1
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pagesize
	return offset, pagesize
}

// 通用搜索条件
type SearchBaseForm struct {
	ID        int64  `form:"id" binding:"-"`
	Keyword   string `form:"keyword" binding:"-"`
	CompanyId int64  `form:"company_id" binding:"-"`
}

const (
	ImageSepChar = ","
)
