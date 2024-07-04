package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tokens []string

var folderPath = "./files"

func main() {
	r := gin.Default()
	r.POST("/login", gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	}), func(c *gin.Context) {
		token, _ := randomHex(20)
		tokens = append(tokens, token)

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	r.GET("/resource", func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(bearerToken, " ")
		if len(splitToken) < 2 {
			// handle the error, for example:
			fmt.Println("Invalid bearer token")
			return
		}
		reqToken := splitToken[1]
		for _, token := range tokens {
			if token == reqToken {
				files, err := readAllFilesInDir(folderPath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "internal server error",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"files": files,
				})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	})

	r.Run() // Listen and serve on 0.0.0.0:8080 (for Windows "localhost:8080")
}

type FileMetaData struct {
	FileName   string
	FileSize   int64
	FileType   string
	UploadTime int64
}

// Function to read all files in a directory
func readAllFilesInDir(dirPath string) ([]FileMetaData, error) {
	// Walk through the directoryx
	var files []FileMetaData
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// Check for errors
		if err != nil {
			return err
		}

		// Check if it's a regular file (not a directory)
		if !info.IsDir() {
			// Open the file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(file)

			files = append(files, FileMetaData{
				FileName:   info.Name(),
				FileSize:   info.Size(),
				FileType:   filepath.Ext(info.Name()),
				UploadTime: info.ModTime().Unix(),
			})
		}

		return nil
	})
	return files, err
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
