// Package service
// @author ufec https://github.com/ufec
// @date 2022/5/22
package service

import (
	"errors"
	"github.com/ufec/douyin_be/config"
	"github.com/ufec/douyin_be/model"
	"gorm.io/gorm"
)

type RelationService struct{}

// FollowUser
//  @Description: 执行用户关注操作
//  @receiver s *RelationService
//  @param fromUserId uint 关注者
//  @param toUserId uint 被关注者
//  @return bool 关注操作是否成功
//  @return error 是否出错
func (s *RelationService) FollowUser(fromUserId, toUserId uint) (bool, error) {
	var relation model.Relation
	if !errors.Is(config.DB.Where("from_user_id = ? AND to_user_id = ?", fromUserId, toUserId).First(&relation).Error, gorm.ErrRecordNotFound) {
		return false, errors.New("已关注该用户")
	}
	// 如果 toUserId  之前关注了 fromUserId  那就设置 互相关注即可
	if !errors.Is(config.DB.Where("from_user_id = ? AND to_user_id = ?", toUserId, fromUserId).First(&relation).Error, gorm.ErrRecordNotFound) {
		relation.IsMutual = 1
		if err := config.DB.Save(&relation).Error; err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	// fromUserId 没有关注 toUserId， 反过来也没有 那就新创建一条记录
	relation.FromUserId, relation.ToUserId = fromUserId, toUserId
	if err := config.DB.Create(&relation).Error; err != nil {
		return false, err
	}
	return true, nil
}

// UnFollowUser
//  @Description: 取消关注
//  @receiver s *RelationService
//  @param fromUserId uint 发起取消关注的用户ID
//  @param toUserId uint 被取消关注的用户ID
//  @return bool 操作是否成功
//  @return error 操作错误信息
func (s *RelationService) UnFollowUser(fromUserId, toUserId uint) (bool, error) {
	var relation model.Relation
	// fromUserId 关注了 toUserId 或者 toUserId 关注了 fromUserId
	result := config.DB.Where("from_user_id=? AND to_user_id=?", fromUserId, toUserId).Or("from_user_id=? AND to_user_id=?", toUserId, fromUserId).First(&relation)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 二者互相关注了
		if relation.IsMutual == 1 {
			relation.IsMutual = 0 // 有一方取消关注 则要解除互相关注标志
		}
		if relation.FromUserId == fromUserId && relation.ToUserId == toUserId {
			relation.FromUserId, relation.ToUserId = relation.ToUserId, relation.FromUserId // 通过交换二者位置，来视作 fromUserId 取消关注 toUserId
		}
		// 更新数据库
		if err := config.DB.Save(&relation).Error; err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("双方尚未建立关系")
}

// GetUserFollowList
//  @Description: 获取用户的关注列表用户ID
//  @receiver s *RelationService
//  @param userId uint
//  @return []uint
//  @return error
func (s *RelationService) GetUserFollowList(userId uint) ([]uint, error) {
	var relations []model.Relation
	config.DB.Where("from_user_id=?", userId).Or("to_user_id=? and is_mutual=1", userId).Find(&relations)
	var followUserIds []uint
	for _, relation := range relations {
		if relation.ToUserId == userId && relation.IsMutual == 1 {
			followUserIds = append(followUserIds, relation.FromUserId)
		} else {
			followUserIds = append(followUserIds, relation.ToUserId)
		}
	}
	return followUserIds, nil
}

// GetUserFollowerList
//  @Description: 获取用户的粉丝列表用户ID
//  @receiver s *RelationService
//  @param userId uint
//  @return []uint
//  @return error
func (s *RelationService) GetUserFollowerList(userId uint) ([]uint, error) {
	var relations []model.Relation
	config.DB.Where("from_user_id=? and is_mutual=1", userId).Or("to_user_id=?", userId).Find(&relations)
	var followerUserIds []uint
	for _, relation := range relations {
		if relation.FromUserId == userId && relation.IsMutual == 1 {
			followerUserIds = append(followerUserIds, relation.ToUserId)
		} else {
			followerUserIds = append(followerUserIds, relation.FromUserId)
		}
	}
	return followerUserIds, nil
}

// IsFollow
//  @Description: 判断是否关注
//  @receiver s *RelationService
//  @param fromUserId uint 关注者
//  @param toUserId uint 被关注者
//  @return bool
func (s *RelationService) IsFollow(fromUserId, toUserId uint) bool {
	err := config.DB.Where("from_user_id=? and to_user_id=?", fromUserId, toUserId).Or("from_user_id=? and to_user_id=? and is_mutual=1", toUserId, fromUserId).First(&model.Relation{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
