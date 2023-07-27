package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context) any

type Error struct {
	Code  int
	Error string
}

func myAdapter(h Handler) gin.HandlerFunc {

	return func(c *gin.Context) {

		res := h(c)

		switch v := res.(type) {
		case gin.H:
			// here v has type T
			c.JSON(
				http.StatusOK,
				v,
			)
		case Error:
			// here v has type S
			c.JSON(
				v.Code,
				gin.H{"error": v.Error},
			)
		default:
			c.JSON(
				http.StatusExpectationFailed,
				gin.H{
					"error": "Unknown handler response",
				},
			)
		}

	}

}

func main() {

	r := gin.Default()

	r.GET("/test", myAdapter(HelloHandler))
	r.GET("/test_error", myAdapter(ErrorHandler))
	r.GET("/test_not_found", myAdapter(ErrorNotFoundHandler))

	err := r.Run()

	if err != nil {
		panic(err)
	}
}

func HelloHandler(c *gin.Context) any {

	return gin.H{
		"data": "Hello World",
	}
}

func ErrorHandler(c *gin.Context) any {

	return Error{
		http.StatusInternalServerError,
		"Sample error",
	}
}

func ErrorNotFoundHandler(c *gin.Context) any {

	return Error{
		http.StatusNotFound,
		"Error not found",
	}
}
