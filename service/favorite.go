// Package service
// @author ufec https://github.com/ufec
// @date 2022/5/21
package service

import (
	"errors"
	"github.com/ufec/douyin_be/config"
	"github.com/ufec/douyin_be/model"
	"gorm.io/gorm"
)

type FavoriteService struct {
}

// Action
//  @Description: 点赞/取消点赞操作
//  @receiver s *FavoriteService
//  @param userId uint	用户ID
//  @param videoId uint	 视频ID
//  @param actionType string 点赞 "1" 取消点赞 "2"
//  @return bool 操作是否成功
//  @return error 操作中产生的错误
func (s *FavoriteService) Action(userId, videoId uint, actionType string) (bool, error) {
	favorite := model.Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	findRes := config.DB.Where("user_id = ? and video_id = ?", userId, videoId).First(&favorite)
	switch actionType {
	case "1":
		if errors.Is(findRes.Error, gorm.ErrRecordNotFound) {
			// 给一条未点赞过的视频点赞
			config.DB.Create(&favorite)
			if favorite.ID != 0 {
				return true, nil
			}
			return false, errors.New("点赞失败")
		}
		// 状态为 0 即为未点赞状态
		if favorite.Status == 0 {
			favorite.Status = 1
			config.DB.Save(&favorite)
			return true, nil
		}
		return false, errors.New("已经点过赞了")
	case "2":
		if errors.Is(findRes.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("没有点过赞无法取消点赞")
		}
		favorite.Status = 0
		config.DB.Save(&favorite)
		return true, nil
	default:
		return false, errors.New("非法操作")
	}
}

// GetFavoriteListVid
//  @Description: 获取用户点赞视频的ID
//  @receiver s *FavoriteService
//  @param userId uint 用户ID
//  @return []uint 点赞过的视频ID
//  @return error err
func (s *FavoriteService) GetFavoriteListVid(userId uint) ([]uint, error) {
	var favoriteList []model.Favorite
	err := config.DB.Select("video_id").Where("user_id = ? and status=1", userId).Find(&favoriteList).Error
	videoIds := make([]uint, 0, len(favoriteList))
	for _, favorite := range favoriteList {
		videoIds = append(videoIds, favorite.VideoId)
	}
	return videoIds, err
}

// IsFavorite
//  @Description: 判断用户是否点赞了视频
//  @receiver s *FavoriteService
//  @param userId uint 当前用户ID
//  @param videoId uint	当前视频ID
//  @return bool
func (s *FavoriteService) IsFavorite(userId, videoId uint) bool {
	if errors.Is(config.DB.Where("user_id=? and video_id=? and status=1", userId, videoId).First(&model.Favorite{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
