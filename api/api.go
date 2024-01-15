package api

import "github.com/gin-gonic/gin"

var Router *gin.Engine
var Files *gin.RouterGroup
var Api *gin.RouterGroup

func CreateRouter() {
	Router = gin.Default()
	Files = Router.Group("/files")
	Api = Router.Group("/api")
}
