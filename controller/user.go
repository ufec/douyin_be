package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ufec/douyin_be/model"
	"net/http"
	"strconv"
)

// UsersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var UsersLoginInfo = map[string]model.User{}

type UserLoginResponse struct {
	Response
	UserId int    `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User `json:"user"`
}

// Register
//  @Description: 用户注册接口
//  @param c *gin.Context
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	// 失败抛出异常 成功返回用户信息
	user, err := userService.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	UsersLoginInfo[token] = user
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		UserId: int(user.ID),
		Token:  token,
	})
}

// Login
//  @Description: 用户登录接口
//  @param c *gin.Context
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	user, err := userService.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	UsersLoginInfo[token] = user
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "登陆成功",
		},
		UserId: int(user.ID),
		Token:  token,
	})
}

func UserInfo(c *gin.Context) {
	loginUserId, queryUserIdStr := c.GetUint("userID"), c.Query("user_id")
	var queryUserId uint64
	if queryUserIdStr != "" {
		if _queryUserId, parseUintErr := strconv.ParseUint(queryUserIdStr, 10, 64); parseUintErr != nil {
			Failed(c, parseUintErr.Error())
			return
		} else {
			queryUserId = _queryUserId
		}
	}
	// 默认值为0 主键ID不为0 则说明用户不存在
	if loginUserId == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "不存在该用户"})
		return
	}
	isFollow, userId := false, loginUserId // 默认查当前登录用户信息
	if queryUserId != 0 {
		userId = uint(queryUserId) // 传入的查询用户ID不为空则查传入的用户id信息
		isFollow = IsFollow(loginUserId, uint(queryUserId))
	}
	userInfo, getUserInfoErr := userService.GetUserInfoById(userId)
	if getUserInfoErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: getUserInfoErr.Error()})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: User{
			Id:            int64(userInfo.ID),
			Name:          userInfo.UserName,
			FollowCount:   int64(userInfo.FollowCount),
			FollowerCount: int64(userInfo.FollowerCount),
			IsFollow:      isFollow,
		},
	})
}
