package main

import (
	"app/internal/common"
	"flag"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	c string
)

func newConfig(in string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(in)
	if err := v.ReadInConfig(); err != nil {
		return v, err
	}
	return v, nil
}

func newMySQL(DSN string, ConnMaxLifetime time.Duration, MaxOpenConn, MaxIdleConn int) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		return db, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	sqlDB.SetConnMaxLifetime(ConnMaxLifetime * time.Second)
	sqlDB.SetMaxOpenConns(MaxOpenConn)
	sqlDB.SetMaxIdleConns(MaxIdleConn)

	return db, nil
}

func init() {
	flag.StringVar(&c, "c", "./configs/config.yaml", "config file")
	flag.Parse()

	f := &common.Function{}

	v, err := newConfig(c)
	if err != nil {
		log.Fatal(err)
	}

	db, err := newMySQL(
		v.GetString("database.DSN"),
		v.GetDuration("database.ConnMaxLifetime"),
		v.GetInt("database.MaxOpenConn"),
		v.GetInt("database.MaxIdleConn"),
	)
	if err != nil {
		log.Fatal(err)
	}

	/*db.Logger = logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)*/

	common.Func = f
	common.Config = v
	common.DB = db
}
