package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/albatross-org/sergeant"
	"github.com/gin-gonic/gin"
)

func handlerCardUpdate(c *gin.Context) {
	answer := &CardUpdateJSON{}

	err := c.BindJSON(answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't decode put request: %s", err),
		})
		return
	}

	if answer.ID == "" || answer.Duration == 0 || answer.Answer == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't decode put request: some fields are blank"),
		})
		return
	}

	if answer.Answer != "perfect" && answer.Answer != "minor" && answer.Answer != "major" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid answer field %q: please use 'perfect', 'minor' or 'major'", answer.Answer),
		})
		return
	}

	set, _, err := store.Set("all")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("couldn't get set: %s", err),
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

	var card *sergeant.Card
	for _, searchCard := range set.Cards {
		if searchCard.ID == answer.ID {
			card = searchCard
		}
	}

	err = store.AddCompletion(card.Path, answer.Answer, sergeant.Completion{
		Date:     time.Now(),
		Duration: time.Millisecond * time.Duration(answer.Duration),
		User:     user,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error adding %q completion to card %q: %s", answer.Answer, card.ID, err),
		})
		return
	}
}
