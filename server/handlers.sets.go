package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlerSetsGet(c *gin.Context) {
	c.JSON(http.StatusOK, fakeCard)
}

func handlerSetsList(c *gin.Context) {
	c.JSON(http.StatusOK, fakeSetList)
}

func handlerSetsStats(c *gin.Context) {
	c.JSON(http.StatusOK, fakeHeatmapData())
}
