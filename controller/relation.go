package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	fromUserId := c.GetUint("userID")
	if fromUserId == 0 {
		Failed(c, "用户不存在")
		return
	}
	toUserIdStr, actionTypeStr := c.Query("to_user_id"), c.Query("action_type")
	toUserId, parseUintErr := strconv.ParseUint(toUserIdStr, 10, 64)
	if parseUintErr != nil {
		Failed(c, "非法字段 to_user_id")
		return
	}
	// 自己不能关注/取消关注自己
	if toUserId == uint64(fromUserId) {
		Failed(c, "不能关注自己")
		return
	}
	// actionType 1执行关注操作 2执行取消关注操作
	actionType, parseIntErr := strconv.ParseInt(actionTypeStr, 10, 64)
	if parseIntErr != nil {
		Failed(c, "非法字段 action_type")
		return
	}
	// 被操作人是否存在
	if _, getUserInfoErr := userService.GetUserInfoById(uint(toUserId)); getUserInfoErr != nil {
		Failed(c, getUserInfoErr.Error())
		return
	}
	var diff int
	if actionType == 1 {
		if _, err := relationService.FollowUser(fromUserId, uint(toUserId)); err != nil {
			Failed(c, err.Error())
			return
		}
		diff = 1
	} else {
		if _, err := relationService.UnFollowUser(fromUserId, uint(toUserId)); err != nil {
			Failed(c, err.Error())
			return
		}
		diff = -1
	}
	// 增加 / 减少 关注取消关注 都是同步的
	if _, err := userService.UpdateFollowCountOrFollowerCountById(fromUserId, diff, "FollowCount"); err != nil {
		Failed(c, err.Error())
		return
	}
	if _, err := userService.UpdateFollowCountOrFollowerCountById(uint(toUserId), diff, "FollowerCount"); err != nil {
		Failed(c, err.Error())
		return
	}
	Success(c, "操作成功")
	return
}

// FollowList all users have same followed list
func FollowList(c *gin.Context) {
	loginUserId := c.GetUint("userID")
	toUserIdStr := c.Query("user_id")
	toUserId, parseUintErr := strconv.ParseUint(toUserIdStr, 10, 64)
	if parseUintErr != nil {
		Failed(c, parseUintErr.Error())
		return
	}
	if loginUserId == 0 {
		Failed(c, "用户不存在")
		return
	}
	userId := loginUserId
	if toUserId != 0 {
		userId = uint(toUserId)
	}
	if followList, getFollowErr := getFollowListOrFansList(userId, "follow"); getFollowErr != nil {
		Failed(c, getFollowErr.Error())
	} else {
		c.JSON(http.StatusOK, struct {
			Response
			UserList `json:"user_list"`
		}{
			Response{StatusCode: 0, StatusMsg: "success"},
			*followList,
		})
	}
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	loginUserId := c.GetUint("userID")
	toUserIdStr := c.Query("user_id")
	toUserId, parseUintErr := strconv.ParseUint(toUserIdStr, 10, 64)
	if parseUintErr != nil {
		Failed(c, parseUintErr.Error())
		return
	}
	if loginUserId == 0 {
		Failed(c, "用户不存在")
		return
	}
	userId := loginUserId
	if toUserId != 0 {
		userId = uint(toUserId)
	}
	if followerList, getFollowErr := getFollowListOrFansList(userId, "follower"); getFollowErr != nil {
		Failed(c, getFollowErr.Error())
	} else {
		c.JSON(http.StatusOK, struct {
			Response
			UserList `json:"user_list"`
		}{
			Response{StatusCode: 0, StatusMsg: "success"},
			*followerList,
		})
	}
}

func getFollowListOrFansList(userId uint, flag string) (*[]User, error) {
	var (
		userIds      []uint
		getUserIdErr error
		followIds    []uint
	)
	switch flag {
	case "follow":
		if userIds, getUserIdErr = relationService.GetUserFollowList(userId); getUserIdErr != nil {
			return nil, getUserIdErr
		}
	case "follower":
		if userIds, getUserIdErr = relationService.GetUserFollowerList(userId); getUserIdErr != nil {
			return nil, getUserIdErr
		}
		if followIds, getUserIdErr = relationService.GetUserFollowList(userId); getUserIdErr != nil {
			return nil, getUserIdErr
		}
		sort.Slice(followIds, func(i, j int) bool { return followIds[i] < followIds[j] }) // 排序
	default:
		return nil, errors.New("不支持获取该数据")
	}
	userInfoList, getUsersInfoErr := userService.GetUserInfoByIds(userIds)
	if getUsersInfoErr != nil {
		return nil, getUsersInfoErr
	}
	userList := make([]User, 0, len(userInfoList))
	for _, user := range userInfoList {
		userList = append(userList, User{
			Id:            int64(user.ID),
			Name:          user.UserName,
			FollowCount:   int64(user.FollowCount),
			FollowerCount: int64(user.FollowerCount),
			IsFollow: func(ids *[]uint) bool {
				if flag == "follow" {
					return true
				}
				if len(followIds) == 0 {
					return false
				}
				return searchWithSortUintSlice(ids, user.ID) // 判断当前用户id是否在 followIds 中 存在为true，不存在为false
			}(&followIds),
		})
	}
	return &userList, nil
}

func searchWithSortUintSlice(slice *[]uint, key uint) bool {
	if len(*slice) == 0 {
		return false
	}
	l, r := uint(0), uint(len(*slice)-1)
	for l < r {
		mid := l + (r-l)>>1
		if (*slice)[mid] == key {
			return true
		} else if (*slice)[mid] > key {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return false
}
