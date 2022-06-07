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
			if parseIntRes > 1000000000000 {
				parseIntRes /= 1000
			}
			startTime = time.Unix(parseIntRes, 0).Format("2006-01-02 15:04:05")
		}
	} else {
		startTime = time.Now().Format("2006-01-02 15:04:05") // 不传默认为当前服务器时间
	}
	userId := UsersLoginInfo[token].ID
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
	videoList := make([]Video, 0, lenFeedVideoList)
	for _, video := range feedVideoList {
		videoList = append(videoList, Video{
			Id: int64(video.ID),
			Author: User{
				Id:            int64(video.User.ID),
				Name:          video.User.UserName,
				FollowCount:   int64(video.User.FollowerCount),
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
	// 这里需要处理用户登陆的逻辑 (登陆了优先推荐他关注的人发布的视频)
	if userId != 0 {
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
