package setting

import (
	"context"
	"errors"
	"fmt"
	"golang-todo-app-atp/util"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	getEnv = util.GoDotEnvVariable
	db     *gorm.DB
)

type SlqLogger struct {
	logger.Interface
}

func (l SlqLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n ========================================================\n", sql)
}

func InitTimeZone() error {
	// set timezone ให้กับระบบ ป้องกันการมีปัญหาหากไปใช้ container
	timeZone, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return errors.New("cannot set timezone !!")
	}

	time.Local = timeZone
	return nil
}

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		getEnv("DB_USERNAME", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", ""),
	)

	dial := mysql.Open(dsn)
	logrus.Infoln("Init database connection")
	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		// Logger: &SlqLogger{},
		DryRun: false,
	})
	if err != nil {
		return nil, err
	}

	tables := DatabaseSchema
	err = MigrateDatabase(tables)
	if err != nil {
		logrus.Errorln("cannot connection gorm")
		return nil, errors.New(fmt.Sprintf("err is : %v", err))
	}

	return db, nil
}

func MigrateDatabase(tables []interface{}) error {
	tx := db.Begin()
	for _, t := range tables {
		// fmt.Printf("%v\n", t)
		if err := tx.AutoMigrate(t); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
	// return nil
}
