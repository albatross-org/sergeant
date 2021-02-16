package server

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	fakeCard = gin.H{
		"path":        "further-maths/core-pure-1/chapter-4-roots-of-polynomials/ex4a",
		"questionImg": "https://i.imgur.com/YawL568.png",
		"answerImg":   "https://i.imgur.com/TP2SsRW.png",
		"id":          "0NiDQqGdzxTSipJa",
	}

	fakeSetList = []gin.H{
		{
			"displayName": "All",
			"name":        "all",
			"description": "This set contains all cards added to the program.",
			"background":  "linear-gradient(315deg, #BD4F6C 0%, #D7816A 74%)",
			"color":       "#D7816A",
		},
		{
			"displayName": "Further Maths",
			"name":        "further-maths",
			"description": "This set contains all Further Maths cards.",
			"background":  "linear-gradient(315deg, #3c1053, #ad5389)",
			"color":       "#ad5389",
		},
		{
			"displayName": "Maths",
			"name":        "maths",
			"description": "This set contains all Maths cards.",
			"background":  "linear-gradient(315deg, #0083b0, #00b4db)",
			"color":       "#00b4db",
		},
		{
			"displayName": "Physics",
			"name":        "physics",
			"description": "This set contains all Physics cards.",
			"background":  "linear-gradient(315deg, #CAC531, #F3F9A7)",
			"color":       "#F3F9A7",
		},
		{
			"displayName": "Computing",
			"name":        "computing",
			"description": "This set contains all of the Computing cards.",
			"background":  "linear-gradient(315deg, #11998e, #38ef7d)",
			"color":       "#38ef7d",
		},
	}
)

func fakeHeatmapData() []gin.H {
	days := []gin.H{}

	for i := 0; i < 100; i++ {
		if rand.Intn(3) == 1 {
			day := time.Now().Add(time.Hour * time.Duration(24*i))
			days = append(days, gin.H{
				"day":   day.Format("2006-01-02"),
				"value": rand.Intn(100),
			})
		}
	}

	return days
}
