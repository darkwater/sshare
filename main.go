package main

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

var config struct {
	rootURL          string
	storagePath      string
	storagePerm      int
	randomCharacters string
	userLength       int
	filenameLength   int
}

var store struct {
	users map[string]string
}

type user string

type response struct {
	URL string `json:"url"`
}

func getUser(key string) (string, error) {
	if key == "" {
		return "", errors.New("invalid sshare key")
	}

	user := store.users[key]
	if key == "" {
		return "", errors.New("unknown sshare key")
	}

	return user, nil
}

func handleUpload(c echo.Context) error {
	// auth user
	user, err := getUser(c.Request().Header["Key"][0])
	if err != nil {
		return err
	}

	userpath := config.storagePath + "/" + user + "/"
	os.MkdirAll(userpath, 0755)

	// open input file
	file, err := c.FormFile("contents")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// create output file
	filename := generateFilename()
	if i := strings.IndexRune(file.Filename, '.'); i > 0 {
		// add original extension
		// intentionally don't match a dot at the start of the filename
		filename += file.Filename[i:]
	}
	dst, err := os.Create(userpath + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// write to file
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	res := response{
		URL: config.rootURL + user + "/" + filename,
	}

	return c.JSON(http.StatusOK, res)
}

func generateFilename() string {
	// TODO: check for collisions
	name := make([]byte, config.filenameLength)
	for i := range name {
		name[i] = config.randomCharacters[rand.Intn(len(config.randomCharacters))]
	}

	return string(name)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	e := echo.New()

	config.rootURL = "http://172.24.0.2:3636/"
	config.storagePath = "storage"
	config.storagePerm = 0755
	config.randomCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	config.filenameLength = 5

	store.users = map[string]string{
		"deadbeef": "xf",
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// serve web interface
	e.Static("/", "public")

	// upload endpoint
	e.POST("/upload", handleUpload)

	// API endpoint
	e.POST("/api", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["foo"] = "bar"
		sess.Save(c.Request(), c.Response())
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start("172.24.0.2:3636"))
}
