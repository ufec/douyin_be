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
	ServerDomain = "http://serverIp:port/" // 服务端地址，访问视频|图片资源
)
