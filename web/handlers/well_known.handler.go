package handlers

import "github.com/gin-gonic/gin"

func AssetLinks(c *gin.Context) {
	c.File("./web/.well-known/assetlinks.json")
}
