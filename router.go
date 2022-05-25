package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ufec/douyin_be/controller"
	"github.com/ufec/douyin_be/middleware"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/public", "./public")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.Use(gin.Recovery())
	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.JWTAuth("query"), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/publish/list/", middleware.JWTAuth("query"), controller.PublishList)
	apiRouter.POST("/publish/action/", middleware.JWTAuth("form"), controller.Publish)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JWTAuth("query"), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.JWTAuth("query"), controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.JWTAuth("query"), controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JWTAuth("query"), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.JWTAuth("query"), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.JWTAuth("query"), controller.FollowerList)
}
