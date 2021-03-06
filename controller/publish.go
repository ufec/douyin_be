package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ufec/douyin_be/config"
	"github.com/ufec/douyin_be/utils"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	userId := c.GetUint("userID")
	// 默认值为0 主键ID不为0 则说明用户不存在
	if userId == 0 {
		Failed(c, "不存在该用户")
		return
	}
	title := c.PostForm("title")
	file, getUploadFileErr := c.FormFile("data")
	if getUploadFileErr != nil {

		return
	}
	pwd, getPwdErr := os.Getwd()
	if getPwdErr != nil {
		Failed(c, getPwdErr.Error())
		return
	}
	nowTime, absoluteSaveDir := time.Now(), ""
	if runtime.GOOS == "windows" {
		absoluteSaveDir += fmt.Sprintf(".\\public\\%d\\%d\\%d\\", nowTime.Year(), nowTime.Month(), nowTime.Day())
	} else {
		absoluteSaveDir += fmt.Sprintf("./public/%d/%d/%d/", nowTime.Year(), nowTime.Month(), nowTime.Day())
	}
	// 不存在该目录则自动创建
	if !utils.PathExists(absoluteSaveDir) {
		if mkDirErr := utils.MakeDir(absoluteSaveDir); mkDirErr != nil {
			fmt.Println(mkDirErr.Error())
			Failed(c, "创建目录失败")
			return
		}
	}
	// 用户id_文件名_文件大小 拼接文件名 用于后续制作视频封面 提取文件后缀 后续对文件名进一步处理
	fileName, fileExt := fmt.Sprintf("%d_%s_%d", userId, file.Filename, file.Size), filepath.Ext(file.Filename)
	// 用时间戳 对拼接后的文件名进行 hmac_sha256 散列 输出bas464格式
	saveFileName := utils.HmacSha256(fileName, strconv.FormatInt(nowTime.Unix(), 10), "hex")
	// 最终保存的目录+文件名组成为最终该文件被存储的路径
	absoluteSaveVideoFile := filepath.Join(absoluteSaveDir, saveFileName+fileExt)

	// 保存上传的视频文件
	if err := c.SaveUploadedFile(file, filepath.Join(pwd, absoluteSaveVideoFile)); err != nil {
		fmt.Println(err.Error())
		Failed(c, "保存视频文件失败")
		return
	}
	// 视频保存成功后 制作视频封面
	absoluteSaveThumbnailFile := filepath.Join(absoluteSaveDir, saveFileName+"_thumbnail.png")
	if err := utils.BuildThumbnailWithVideo(filepath.Join(pwd, absoluteSaveVideoFile), filepath.Join(pwd, absoluteSaveThumbnailFile)); err != nil {
		fmt.Println(err.Error())
		Failed(c, "封面图生成失败")
		return
	}
	// 保存到数据库
	if _, createVideoErr := videoService.Create(absoluteSaveVideoFile, absoluteSaveThumbnailFile, title, userId); createVideoErr != nil {
		fmt.Println(createVideoErr.Error())
		Failed(c, "数据保存失败")
		return
	}
	Success(c, "视频发布成功")
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "不存在该用户",
			},
			VideoList: []Video{},
		})
		return
	}
	userPublishList := videoService.UserPublishList(userId)
	if len(userPublishList) == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: []Video{},
		})
		return
	}
	videoList := make([]Video, 0, len(userPublishList))
	for i := 0; i < len(userPublishList); i++ {
		videoList = append(videoList, Video{
			Id: int64(userPublishList[i].ID),
			Author: User{
				Id:            int64(userPublishList[i].User.ID),
				Name:          userPublishList[i].User.UserName,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
			},
			PlayUrl:       config.ServerDomain + userPublishList[i].PlayUrl,
			CoverUrl:      config.ServerDomain + userPublishList[i].CoverUrl,
			FavoriteCount: userPublishList[i].FavoriteCount,
			CommentCount:  userPublishList[i].CommentCount,
			IsFavorite:    false,
			Title:         userPublishList[i].Description,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
	})
}
