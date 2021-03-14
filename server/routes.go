package server

import (
	"github.com/gin-gonic/contrib/static"
)

func initRoutes() {
	// Serve the web UI at the base path.
	router.Use(
		static.Serve("/", static.LocalFile("./public/build", true)),
	)

	// Set up basic routes for the V1 api.
	api := router.Group("/api/v1")
	{
		cards := api.Group("/cards")
		{
			cards.PUT("/update", handlerCardUpdate)
		}

		sets := api.Group("/sets")
		{
			sets.GET("/get", handlerSetsGet)
			sets.GET("/list", handlerSetsList)
			sets.GET("/stats", handlerSetsStats)
		}
	}

}
