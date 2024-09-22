package data

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/f-rambo/ship/internal/biz"
	"github.com/f-rambo/ship/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSystemRepo, NewCloudRepo)

type DBDriver string

const (
	DBDriverSQLite DBDriver = "sqlite"
)

func (d DBDriver) String() string {
	return string(d)
}

const (
	DatabaseName = "ship.db"
)

type Data struct {
	db            *gorm.DB
	dbLoggerLevel gormLogger.LogLevel
	log           *log.Helper
}

// NewData .
func NewData(logger log.Logger) (*Data, func(), error) {
	cleanup := func() {}
	data := &Data{
		log:           log.NewHelper(logger),
		dbLoggerLevel: gormLogger.Warn,
	}
	dbFilePath, err := utils.GetPackageStorePathByNames(DBDriverSQLite.String(), DatabaseName)
	if err != nil {
		return nil, cleanup, err
	}
	if dbFilePath != "" && !utils.IsFileExist(dbFilePath) {
		dir, _ := filepath.Split(dbFilePath)
		os.MkdirAll(dir, 0755)
		file, err := os.Create(dbFilePath)
		if err != nil {
			return nil, cleanup, err
		}
		file.Close()
	}
	client, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: data,
	})
	if err != nil {
		return nil, cleanup, err
	}
	// AutoMigrate
	err = client.AutoMigrate(
		&biz.System{},
		&biz.Cloud{},
	)
	if err != nil {
		return nil, cleanup, err
	}
	data.db = client
	return data, cleanup, nil
}

func (d *Data) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	d.dbLoggerLevel = level
	return d
}

func (d *Data) Info(ctx context.Context, msg string, args ...interface{}) {
	if d.dbLoggerLevel >= gormLogger.Info {
		d.log.WithContext(ctx).Infof(msg, args...)
	}
}

func (d *Data) Warn(ctx context.Context, msg string, args ...interface{}) {
	if d.dbLoggerLevel >= gormLogger.Warn {
		d.log.WithContext(ctx).Warnf(msg, args...)
	}
}

func (d *Data) Error(ctx context.Context, msg string, args ...interface{}) {
	if d.dbLoggerLevel >= gormLogger.Error {
		d.log.WithContext(ctx).Errorf(msg, args...)
	}
}

func (d *Data) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if d.dbLoggerLevel >= gormLogger.Info {
		sql, rows := fc()
		d.log.WithContext(ctx).Infof("begin: %s, sql: %s, rows: %d, err: %v", begin.Format("2006-01-02 15:04:05"), sql, rows, err)
	}
}
