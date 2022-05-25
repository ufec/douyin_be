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

type VideoService struct {
}

// Create
//  @Description: 新增一个视频
//  @receiver videoService *VideoService
//  @param playUrl string 视频播放地址
//  @param coverUrl string 视频封面地址
//  @param desc string	视频描述
//  @param userId uint	发布视频的用户ID
//  @return model.Video
//  @return error
func (videoService *VideoService) Create(playUrl, coverUrl, desc string, userId uint) (model.Video, error) {
	video := model.Video{
		UserId:      userId,
		PlayUrl:     playUrl,
		CoverUrl:    coverUrl,
		Description: desc,
	}
	return video, config.DB.Create(&video).Error
}

// Feed
//  @Description: 获取视频Feed
//  @receiver videoService *VideoService
//  @param startTime string	起始的时间
//  @return *[]model.Video
func (videoService *VideoService) Feed(startTime string) *[]model.Video {
	var videoList *[]model.Video
	config.DB.Where("created_at <= ?", startTime).Preload("User").Order("created_at DESC").Limit(30).Find(&videoList)
	return videoList
}

// UserPublishList
//  @Description: 获取指定用户发布的视频列表
//  @receiver videoService *VideoService
//  @param userId uint	指定的用户ID
//  @return []model.Video
func (videoService *VideoService) UserPublishList(userId uint) []model.Video {
	var videoList []model.Video
	config.DB.Where("user_id = ?", userId).Preload("User").Find(&videoList)
	return videoList
}

// GetVideoInfoByIds
//  @Description: 批量获取视频信息
//  @receiver videoService *VideoService
//  @param videoIds []uint 视频ID
//  @return videoList []model.Video
//  @return err error
func (videoService *VideoService) GetVideoInfoByIds(videoIds []uint) (videoList []model.Video, err error) {
	err = config.DB.Preload("User").Find(&videoList, videoIds).Error
	return videoList, err
}

func (videoService *VideoService) GetVideoInfoById(videoId uint) {

}

// UpdateNumberField
//  @Description: 更新Video表中数值字段
//  @receiver videoService *VideoService
//  @param videoId uint 视频ID
//  @param diff int 更新的差值 增加传正值反之负值
//  @param filed string	字段可选值 FavoriteCount | CommentCount
//  @return bool 更新是否成功
//  @return error 更新是否出错
func (videoService *VideoService) UpdateNumberField(videoId uint, diff int, filed string) (bool, error) {
	var video model.Video
	if errors.Is(config.DB.Where("id=?", videoId).First(&video).Error, gorm.ErrRecordNotFound) {
		return false, errors.New("不存在该视频")
	}
	switch filed {
	case "FavoriteCount":
		if video.FavoriteCount+int64(diff) < 0 {
			return false, errors.New("非法数值")
		}
		video.FavoriteCount += int64(diff)
	case "CommentCount":
		if video.CommentCount+int64(diff) < 0 {
			return false, errors.New("非法数值")
		}
		video.CommentCount += int64(diff)
	default:
		return false, errors.New("不支持该字段")
	}
	err := config.DB.Save(&video).Error
	if err != nil {
		//  TODO:log记录错误
		return false, errors.New("系统错误")
	}
	return true, nil
}
