// Package config
// @author ufec https://github.com/ufec
// @date 2022/5/9
package config

import "gorm.io/gorm"

// Conf
//  @Description: 系统所有配置项
type Conf struct {
	Mysql  *MysqlConfig
	Server *ServerConfig
}

var (
	Config       *Conf
	GormConfig   *gorm.Config
	DB           *gorm.DB
	ServerDomain = "http://192.168.137.1:8080/"
)
