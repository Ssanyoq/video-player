package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env_example")
	if err != nil {
		println("Make a .env please")
		panic(err)
	}

	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	port := os.Getenv("PORT")
	fmt.Printf("Website example is running at http://localhost:%s\n", port) // localhost and not 127.0.0.1 since youtube doesn't like direct IPs idk why
	r.Run(":" + port)                                                       // Gin can tell that this is a port if : are present
}
