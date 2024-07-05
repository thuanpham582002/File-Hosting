package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var folderPath = "./files"

const userkey = "user"

var secret = []byte("secret")

func main() {
	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
		return
	}
}

func engine() *gin.Engine {
	r := gin.New()
	r.LoadHTMLGlob("templates/*")
	// Setup the cookie store for session management
	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))
	// Login and logout routes
	r.POST("/login", login)
	r.GET("/login", loginForm)
	r.GET("/logout", logout)

	// Private group, require authentication to access
	private := r.Group("/private")
	private.Use(AuthRequired)
	{
		private.GET("/dashboard", dashBoard)
		private.GET("/rename", rename)
		private.POST("/upload", uploadFile)
		private.GET("/me", me)
		private.GET("/status", status)
	}
	return r
}

func uploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filePath := filepath.Join(folderPath, file.Filename)
	if err := context.SaveUploadedFile(file, filePath); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func loginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login_page.html", gin.H{})
}

// AuthRequired is a simple middleware to check the session.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.Redirect(http.StatusFound, "/login")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// login is a handler that parses a form and checks for specific data.
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	if username != "admin" || password != "123456aA@" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	// navigate to the dashboard
	c.Redirect(http.StatusFound, "/private/dashboard")
}

// logout is the handler called for the user to log out.
func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func rename(c *gin.Context) {
	fmt.Println(c.Request.Body)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// me is the handler that will return the user information stored in the
// session.
func me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// status is the handler that will tell the user whether it is logged in or not.
func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

func dashBoard(c *gin.Context) {
	files, err := readAllFilesInDir(folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// convert files to json
	filesJson, err := json.Marshal(files)
	fmt.Println(string(filesJson))
	c.HTML(http.StatusOK, "dash_board.html", gin.H{"files": string(filesJson)})
}

type FileMetaData struct {
	FileName   string
	FileSize   int64
	FileType   string
	UploadTime int64
}

func readAllFilesInDir(dirPath string) ([]FileMetaData, error) {
	var files []FileMetaData
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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
