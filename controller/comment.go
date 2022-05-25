package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	actionTypeStr, videoIdStr := c.Query("action_type"), c.Query("video_id")
	videoId, parseVideoId := strconv.ParseUint(videoIdStr, 10, 64)
	if parseVideoId != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: parseVideoId.Error()})
		return
	}
	actionType, parseActionType := strconv.ParseUint(actionTypeStr, 10, 64)
	if parseActionType != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: parseActionType.Error()})
		return
	}
	// 发布评论
	if actionType == 1 {
		//TODO:考虑注入问题
		commentText := c.Query("comment_text")
		if err := commentService.PostComment(userId, 0, uint(videoId), commentText); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		if _, err := videoService.UpdateNumberField(uint(videoId), 1, "CommentCount"); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
		return
	}
	// 删除评论
	commentIdStr := c.Query("comment_id")
	commentId, parseCommentIdErr := strconv.ParseUint(commentIdStr, 10, 64)
	if parseCommentIdErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: parseCommentIdErr.Error()})
		return
	}
	if err := commentService.DeleteComment(userId, 0, uint(videoId), uint(commentId)); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if _, err := videoService.UpdateNumberField(uint(videoId), -1, "CommentCount"); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
	return
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	userId := UsersLoginInfo[token].ID
	videoIdStr := c.Query("video_id")
	videoId, parseVideoId := strconv.ParseUint(videoIdStr, 10, 64)
	if parseVideoId != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: parseVideoId.Error()})
		return
	}
	comments, getCommentsErr := commentService.GetCommentListByVideoId(uint(videoId))
	if getCommentsErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: getCommentsErr.Error()})
		return
	}
	commentList := make([]Comment, 0, len(comments))
	for _, comment := range comments {
		commentList = append(commentList, Comment{
			Id: int64(comment.ID),
			User: User{
				Id:            int64(comment.User.ID),
				Name:          comment.User.UserName,
				FollowCount:   int64(comment.User.FollowCount),
				FollowerCount: int64(comment.User.FollowerCount),
				IsFollow:      IsFollow(userId, comment.User.ID),
			},
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commentList,
	})
}
