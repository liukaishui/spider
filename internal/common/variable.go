package common

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Func   *Function
	Config *viper.Viper
	DB     *gorm.DB
)
