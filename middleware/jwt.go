// Package middleware
// @author ufec https://github.com/ufec
// @date 2022/5/11
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ufec/douyin_be/controller"
	"net/http"
)

// JWTAuth
//  @Description: 类似于JWT中间件
//  @param where string 由于请求token位置不固定，通过指定位置来获取token
//  @return gin.HandlerFunc
func JWTAuth(where string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		// where可选值: query/form/header/body json
		switch where {
		case "query":
			token = c.Query("token")
		case "form":
			token = c.PostForm("token")
		default:
			token = c.Query("token")
		}
		// 不存在该用户token则直接抛出用户不存在错误信息
		if _, exists := controller.UsersLoginInfo[token]; !exists {
			c.JSON(http.StatusOK, controller.Response{StatusCode: 1, StatusMsg: "token鉴权失败, 非法操作"})
			c.Abort()
			return
		}
		// 统一获取user的位置 后续流程直接从上下文取 user 即可
		// 此方法不可取 gin 不支持设置指定数据类型的数据，需要通过json序列化 反序列化来完成类型转换 损失性能
		//c.Set("user", controller.UsersLoginInfo[token])

		// 设置Token也是不必要的 从整体来看 我们仅仅只需要用户ID便能唯一确定用户, gin 也支持获取基础数据类型 恰好符合要求
		//c.Set("token", token)
		c.Set("userID", controller.UsersLoginInfo[token].ID)
		c.Next()
	}
}
