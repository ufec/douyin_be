package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ufec/douyin_be/config"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list"`
	NextTime  int64   `json:"next_time"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	token, lastTimestamp := c.Query("token"), c.Query("latest_time")
	startTime := ""
	if lastTimestamp != "" {
		if parseIntRes, parseIntErr := strconv.ParseInt(lastTimestamp, 10, 64); parseIntErr == nil {
			startTime = time.Unix(parseIntRes/1000, 0).Format("2006-01-02 15:04:05")
		}
	}
	userId := UsersLoginInfo[token].ID
	// 这里需要处理用户登陆的逻辑 (登陆了优先推荐他关注的人发布的视频)
	feedVideoList := *videoService.Feed(startTime)
	lenFeedVideoList := len(feedVideoList)
	if lenFeedVideoList == 0 {
		// 空数据处理
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: []Video{},
			NextTime:  time.Now().Unix(),
		})
		return
	}
	nextTime := feedVideoList[lenFeedVideoList-1].CreatedAt.Unix() // 防止后续的排序影响
	videoList := make([]Video, lenFeedVideoList)
	for i := 0; i < lenFeedVideoList; i++ {
		videoList[i] = Video{
			Id: int64(feedVideoList[i].ID),
			Author: User{
				Id:            int64(feedVideoList[i].User.ID),
				Name:          feedVideoList[i].User.UserName,
				FollowCount:   int64(feedVideoList[i].User.FollowerCount),
				FollowerCount: int64(feedVideoList[i].User.FollowerCount),
				IsFollow:      IsFollow(userId, feedVideoList[i].User.ID),
			},
			PlayUrl:       config.ServerDomain + feedVideoList[i].PlayUrl,
			CoverUrl:      config.ServerDomain + feedVideoList[i].CoverUrl,
			FavoriteCount: feedVideoList[i].FavoriteCount,
			CommentCount:  feedVideoList[i].CommentCount,
			IsFavorite:    IsFavorite(userId, feedVideoList[i].ID),
			Title:         feedVideoList[i].Description,
		}
	}
	if userId == 0 {
		sort.Slice(videoList, func(i, j int) bool { return videoList[i].Author.IsFollow || videoList[j].Author.IsFollow })
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}

// IsFavorite
//  @Description: 判断用户是否点赞当前视频
//  @param userId uint 当前登录用户ID
//  @param videoId uint	当前视频ID
//  @return bool 未登录则直接返回false
func IsFavorite(userId, videoId uint) bool {
	if userId == 0 {
		return false
	}
	return favoriteService.IsFavorite(userId, videoId)
}

// IsFollow
//  @Description: 判断登录用户是否关注了视频作者
//  @param fromUserId uint 登录用户
//  @param toUserId uint 视频作者
//  @return bool 未登录直接返回false
func IsFollow(fromUserId, toUserId uint) bool {
	if fromUserId == 0 {
		return false
	}
	return relationService.IsFollow(fromUserId, toUserId)
}