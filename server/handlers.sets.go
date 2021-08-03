package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/albatross-org/sergeant"
	"github.com/gin-gonic/gin"
)

func handlerSetsGet(c *gin.Context) {
	setConfig, err := setConfigFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	viewName, exists := c.GetQuery("viewName")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify a viewName query parameter",
		})
		return
	}

	view := sergeant.DefaultViews[viewName]
	if view == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the view %q doesn't exist", viewName),
		})
		return
	}

	set, _, err := store.SetFromConfig(setConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("error loading set %q: %s", setConfig.Name, err),
		})
		return
	}

	if len(set.Cards) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("there's no cards left in the %q set", setConfig.Name),
		})
		return
	}

	var user string

	auth := c.Request.Header["Authorization"]

	if len(auth) == 1 {
		if strings.HasPrefix(auth[0], "Basic ") {
			b64 := strings.TrimPrefix(auth[0], "Basic ")
			authBytes, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				user = ""
			} else {
				user = strings.Split(string(authBytes), ":")[0]
			}
		}
	}

	card := view.Next(set, user)
	if card == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't get a card from the %q view", viewName),
		})
		return
	}

	cardJSON, err := cardToJSON(card)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("couldn't turn card into JSON: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, cardJSON)
}

func handlerSetsList(c *gin.Context) {
	c.JSON(http.StatusOK, getSetListJSON())
}

func handlerSetsStats(c *gin.Context) {
	setConfig, err := setConfigFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user string

	auth := c.Request.Header["Authorization"]

	if len(auth) == 1 {
		if strings.HasPrefix(auth[0], "Basic ") {
			b64 := strings.TrimPrefix(auth[0], "Basic ")
			authBytes, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				user = ""
			} else {
				user = strings.Split(string(authBytes), ":")[0]
			}
		}
	}

	set, _, err := store.SetFromConfig(setConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("error loading set %q: %s", setConfig.Name, err),
		})
		return
	}

	c.JSON(http.StatusOK, getSetHeatmapJSON(set, user))
}
