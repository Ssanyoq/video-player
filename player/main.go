package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.StaticFS("/", http.Dir("src/media"))
	router.Run(":8080")
}
