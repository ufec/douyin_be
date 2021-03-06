// Package initalize
// @author ufec https://github.com/ufec
// @date 2022/5/9
package initalize

import (
	"github.com/ufec/douyin_be/config"
	"github.com/ufec/douyin_be/model"
	"gorm.io/gorm"
)

// InitGorm
//  @Description: 初始化Gorm
//  @return *gorm.DB
func InitGorm() *gorm.DB {
	switch config.Config.Server.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

// CreateTable
//  @Description: 根据模型自动创建表
func CreateTable(db *gorm.DB) error {
	err := db.AutoMigrate(
		model.User{},
		model.Video{},
		model.Favorite{},
		model.Comment{},
		model.Relation{},
	)
	if err != nil {
		return err
	}
	return nil
}
