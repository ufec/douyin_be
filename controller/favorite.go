package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ufec/douyin_be/config"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		Failed(c, "用户不存在")
		return
	}
	// 这里不应该从Query中取出用户ID
	videoId, actionType := c.Query("video_id"), c.Query("action_type")
	uintVideoId, parseUintErr := strconv.ParseUint(videoId, 10, 64)
	if parseUintErr != nil {
		Failed(c, "视频ID非法")
		return
	}
	if actionType != "1" && actionType != "2" {
		Failed(c, "非法操作")
		return
	}
	ok, actionLikeErr := favoriteService.Action(userId, uint(uintVideoId), actionType)
	if !ok || actionLikeErr != nil {
		Failed(c, "操作失败")
		return
	}
	// 同步 FavoriteCount 字段值
	if actionType == "1" {
		if _, err := videoService.UpdateNumberField(uint(uintVideoId), 1, "FavoriteCount"); err != nil {
			Failed(c, err.Error())
			return
		}
	} else {
		if _, err := videoService.UpdateNumberField(uint(uintVideoId), -1, "FavoriteCount"); err != nil {
			Failed(c, err.Error())
			return
		}
	}
	Success(c, "操作成功")
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		Failed(c, "用户不存在")
		return
	}
	// 根据用户ID 取出该用户点赞的所有视频ID
	videoIds, getFLVErr := favoriteService.GetFavoriteListVid(userId)
	if getFLVErr != nil {
		Failed(c, getFLVErr.Error())
		return
	}
	// 根据点赞过的视频ID 取出所有对应的视频信息
	videoInfoList, getVIBIErr := videoService.GetVideoInfoByIds(videoIds)
	if getVIBIErr != nil {
		Failed(c, getVIBIErr.Error())
		return
	}
	favoriteVideoList := make([]Video, 0, len(videoInfoList))
	for _, video := range videoInfoList {
		favoriteVideoList = append(favoriteVideoList, Video{
			Id: int64(video.ID),
			Author: User{
				Id:            int64(video.User.ID),
				Name:          video.User.UserName,
				FollowCount:   int64(video.User.FollowCount),
				FollowerCount: int64(video.User.FollowerCount),
				IsFollow:      IsFollow(userId, video.User.ID),
			},
			PlayUrl:       config.ServerDomain + video.PlayUrl,
			CoverUrl:      config.ServerDomain + video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    IsFavorite(userId, video.ID),
			Title:         video.Description,
		})
	}
	c.JSON(http.StatusOK, ResponseVideoList{
		Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		favoriteVideoList,
	})
}
