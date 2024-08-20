package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func streamVideo(c *gin.Context) {
	videoPath := "src/media/test.mp4"

	file, err := os.Open(videoPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening video file")
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting file stats")
		return
	}

	fileSize := fileStat.Size()

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		c.Header("Content-Type", "video/mp4")
		c.File(videoPath)
		return
	}

	ranges := strings.Split(rangeHeader, "=")
	if len(ranges) < 2 {
		c.String(http.StatusBadRequest, "Invalid Range header")
		return
	}

	rangeParts := strings.Split(ranges[1], "-")
	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid start range")
		return
	}

	end := fileSize - 1
	if len(rangeParts) > 1 && rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid end range")
			return
		}
	}

	if start > end || start >= fileSize || end >= fileSize {
		c.String(http.StatusRequestedRangeNotSatisfiable, "Range not satisfiable")
		return
	}

	c.Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(fileSize, 10))
	c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
	c.Header("Content-Type", "video/mp4")
	c.Status(http.StatusPartialContent)

	bufferSize := 1024 * 1024 // 1 MB buffer
	buffer := make([]byte, bufferSize)

	file.Seek(start, 0)

	for {
		if start > end {
			break
		}

		bytesRead, err := file.Read(buffer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file")
			return
		}

		if int64(bytesRead) > end-start+1 {
			bytesRead = int(end-start) + 1
		}

		c.Writer.Write(buffer[:bytesRead])
		start += int64(bytesRead)
	}
}

func main() {
	r := gin.Default()
	r.GET("/video", streamVideo)
	r.Run(":8080")
}
