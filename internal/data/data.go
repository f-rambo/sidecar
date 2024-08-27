package data

import (
	"github.com/f-rambo/ship/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSystemRepo, NewClusterRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(logger log.Logger) (*Data, func(), error) {
	cleanup := func() {}
	dbFilePath, err := utils.GetPackageStorePathByNames("database", "ship.db")
	if err != nil {
		return nil, cleanup, err
	}
	if dbFilePath != "" && !utils.IsFileExist(dbFilePath) {
		path, filename := utils.GetFilePathAndName(dbFilePath)
		file, err := utils.NewFile(path, filename, true)
		if err != nil {
			return nil, cleanup, err
		}
		file.Close()
	}
	if dbFilePath == "" {
		dbFilePath = "file::memory:?cache=shared"
	}
	client, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, cleanup, err
	}
	return &Data{db: client}, cleanup, nil
}
