package arc

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormAdapter struct {
	*gorm.DB
}

func NewGormAdapter() *GormAdapter {
	dbCfg := BeanFactory.Get(IConfig(nil)).(IConfig).GetDB()

	switch dbCfg.GetDriverName() {
	case "mysql":
		return createMysql(dbCfg)
	case "sqlite3":
		return createSqlite(dbCfg)
	default:
		return nil
	}
	// return &GormAdapter{}
}

func (g *GormAdapter) Name() string {
	return "GormAdapter"
}

func (g *GormAdapter) IMigrateTo(des ...interface{}) {
	if err := g.AutoMigrate(des...); err != nil {
		zap.L().Error("DB AutoMigrate Failed.", zap.Error(err))
	}
}

func createMysql(config IDBConfig) *GormAdapter {

	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		zap.L().Error("Create MySQL DB file Failed.", zap.Error(err))
		return nil
	}
	d, _ := db.DB()
	if d != nil {
		connSetting := config.GetConnSetting()
		d.SetMaxIdleConns(connSetting[0])
		d.SetMaxOpenConns(connSetting[1])
		d.SetConnMaxLifetime(time.Second * 30)
	}
	if config.GetMode() != "prod" {
		db = db.Debug()
	}
	return &GormAdapter{DB: db}

}

func createSqlite(config IDBConfig) *GormAdapter {
	dbDir := filepath.Dir(filepath.Clean(config.GetDSN()))
	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		zap.L().Error("Create Sqlite DB file Failed.", zap.Error(err))
		return nil
	}
	db, err := gorm.Open(sqlite.Open(config.GetDSN()), &gorm.Config{SkipDefaultTransaction: true, PrepareStmt: true})
	if err != nil {
		return nil
	}
	d, err := db.DB()
	if err != nil {
		return nil
	}
	d.SetMaxOpenConns(1)
	if config.GetMode() != "prod" {
		db = db.Debug()
	}
	return &GormAdapter{DB: db}
}
