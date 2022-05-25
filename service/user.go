// Package service
// @author ufec https://github.com/ufec
// @date 2022/5/10
package service

import (
	"errors"
	"github.com/ufec/douyin_be/config"
	"github.com/ufec/douyin_be/model"
	"gorm.io/gorm"
)

type UserService struct {
}

// Register
//  @Description: 用户注册服务
//  @receiver userService *UserService
//  @param u model.User
//  @return model.User
//  @return error
func (userService *UserService) Register(username, password string) (model.User, error) {
	var user model.User
	if len(password) <= 5 {
		return user, errors.New("密码长度不得小于5")
	}
	// 用户名不存在则直接注册
	if errors.Is(config.DB.Where("username=?", username).First(&user).Error, gorm.ErrRecordNotFound) {
		user.UserName = username
		user.PassWord = password
		return user, config.DB.Create(&user).Error
	}
	return user, errors.New("用户名已注册")
}

// Login
//  @Description: 用户登录服务
//  @receiver userService *UserService
//  @param u model.User
//  @return model.User
//  @return error
func (userService *UserService) Login(username, password string) (model.User, error) {
	var user model.User
	if errors.Is(config.DB.Where("username=? and password=?", username, password).First(&user).Error, gorm.ErrRecordNotFound) {
		return user, errors.New("用户名或密码错误")
	}
	return user, nil
}

// GetUserInfoById
//  @Description: 通过用户ID获取用户信息
//  @receiver userService *UserService
//  @param userId uint 用户ID
//  @return model.User 用户信息
//  @return error 错误信息
func (userService *UserService) GetUserInfoById(userId uint) (model.User, error) {
	var user model.User
	if errors.Is(config.DB.Where("id = ?", userId).First(&user).Error, gorm.ErrRecordNotFound) {
		return user, errors.New("用户不存在")
	}
	return user, nil
}

// GetUserInfoByIds
//  @Description: 通过ID批量获取用户信息
//  @receiver userService *UserService
//  @param userIds []uint
//  @return []model.User
//  @return error
func (userService *UserService) GetUserInfoByIds(userIds []uint) ([]model.User, error) {
	var users []model.User
	if len(userIds) == 0 {
		return nil, nil
	}
	config.DB.Find(&users, userIds)
	return users, nil
}

// UpdateFollowCountOrFollowerCountById
//  @Description: 更新用户关注数或粉丝数
//  @receiver userService *UserService
//  @param userId uint 操作的用户ID
//  @param diff int	更新值 增加传正值 减少传负值
//  @param field string	操作的字段 直接使用模型中的字段 FollowCount 关注数 FollowerCount 粉丝数
//  @return bool 操作成功
//  @return error 是否错误
func (userService *UserService) UpdateFollowCountOrFollowerCountById(userId uint, diff int, field string) (bool, error) {
	user, getUserInfoErr := userService.GetUserInfoById(userId)
	if getUserInfoErr != nil {
		return false, getUserInfoErr
	}
	switch field {
	case "FollowCount":
		// diff 可以为负值，但必须要检测
		if user.FollowCount+diff < 0 {
			return false, errors.New("更新数据非法")
		}
		user.FollowCount += diff
	case "FollowerCount":
		if user.FollowerCount+diff < 0 {
			return false, errors.New("更新数据非法")
		}
		user.FollowerCount += diff
	default:
		return false, errors.New("不支持此字段")
	}
	if err := config.DB.Save(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}
