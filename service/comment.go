// Package service
// @author ufec https://github.com/ufec
// @date 2022/5/10
package service

import (
	"github.com/ufec/douyin_be/config"
	"github.com/ufec/douyin_be/model"
)

type CommentService struct {
}

// PostComment
//  @Description: 发布评论
//  @receiver s *CommentService
//  @param userId uint 发布评论的用户id
//  @param toUserId uint 回复用户的评论的用户ID
//  @param videoId uint 在哪个视频下评论的
//  @param content string 评论的内容
//  @return error 是否出错
func (s *CommentService) PostComment(userId, toUserId, videoId uint, content string) error {
	comment := model.Comment{
		Content:  content,
		VideoId:  videoId,
		ToUserId: toUserId,
		UserId:   userId,
	}
	return config.DB.Create(&comment).Error
}

// DeleteComment
//  @Description: 删除评论
//  @receiver s *CommentService
//  @param userId uint 发布评论用户的id
//  @param toUserId uint 这条评论恢复的用户的id
//  @param videoId uint	视频id
//  @param commentId uint 评论id
//  @return error
func (s *CommentService) DeleteComment(userId, toUserId, videoId, commentId uint) error {
	// 软删除  使用 Unscoped 永久删除
	return config.DB.Where("user_id=? and to_user_id=? and video_id=? and id=?", userId, toUserId, videoId, commentId).Delete(&model.Comment{}).Error
}

// GetCommentListByVideoId
//  @Description: 获取指定视频id的评论列表
//  @receiver s *CommentService
//  @param videoId uint	视频ID
//  @return []model.Comment	评论列表
//  @return error 错误信息
func (s *CommentService) GetCommentListByVideoId(videoId uint) ([]model.Comment, error) {
	var commentList []model.Comment
	if err := config.DB.Where("video_id=?", videoId).Preload("User").Find(&commentList).Error; err != nil {
		return nil, err
	}
	return commentList, nil
}
