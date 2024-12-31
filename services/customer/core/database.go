package core

import (
	"fmt"
	"log"
	"os"
	"time"

	// mysql driver

	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase returns a gorm.DB struct, gorm.DB.DB() returns a database handle
// see http://golang.org/pkg/database/sql/#DB
func NewDatabase(cfg *Config, dbName string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: gormLogger,
	}

	// Connection args
	args1 := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
	)

	db1, err := gorm.
		Open(
			mysql.New(
				mysql.Config{
					DSN:                       args1, // data source name
					DefaultStringSize:         256,   // default size for string fields
					DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
					DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
					DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
					SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
				},
			),
			gormConfig,
		)
	if err != nil {
		return nil, err
	}
	sqlDb1, err := db1.DB()
	defer sqlDb1.Close()

	db1.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "` DEFAULT CHARACTER SET ascii COLLATE ascii_general_ci;")
	db1.Exec("USE `" + dbName + "`;")

	args2 := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		dbName,
	)
	// fmt.Println(args2)

	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN:                       args2, // data source name
				DefaultStringSize:         256,   // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
			},
		),
		gormConfig,
	)
	if true == cfg.Debug {
		db = db.Debug()
	}
	if err != nil {
		return db, err
	}

	sqlDb, err := db.DB()

	// sqlDb.Set("gorm:table_options", "charset=ascii")
	sqlDb.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)
	// // Database logging
	// sqlDb.LogMode(cfg.Debug)

	return db, nil
}
