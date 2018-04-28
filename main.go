package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type handler struct {
	storagePath      string
	storagePerm      int
	users            map[string]string
	randomCharacters string
	filenameLength   int
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		fmt.Fprintf(w, "expected a multipart/form-data POST request")
		return
	}

	key := r.Header["Key"][0]
	if key == "" {
		fmt.Fprintf(w, "please send a 'Key' header with your sshare key")
		return
	}

	user := h.users[key]
	if key == "" {
		fmt.Fprintf(w, "unknown sshare key")
		return
	}

	userPath := h.storagePath + "/" + user + "/"
	os.MkdirAll(userPath, 0755)

	fileWritten := false

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(w, "unexpected error while reading multipart form data")
			return
		}

		switch part.FormName() {
		case "contents":
			if fileWritten {
				continue
			}

			filename := h.generateFilename()

			if i := strings.IndexRune(part.FileName(), '.'); i > 0 {
				// intentionally don't match a dot at the start of the filename
				filename += part.FileName()[i:]
			}

			// TODO: handle collisions
			filepath := userPath + filename
			file, err := os.Create(filepath)
			if err != nil {
				println(err)
				return
			}
			defer file.Close()

			written, err := io.Copy(file, part)
			if err != nil {
				println(err)
				return
			}
			_ = written

			fileWritten = true
		}
	}
}

func (h handler) generateFilename() string {
	for {
		name := make([]byte, h.filenameLength)
		for i := range name {
			name[i] = h.randomCharacters[rand.Intn(len(h.randomCharacters))]
		}

		return string(name)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	h := handler{
		storagePath: "/tmp/sshare",
		storagePerm: 0755,
		users: map[string]string{
			"deadbeef": "xf",
		},
		randomCharacters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		filenameLength:   5,
	}

	http.Handle("/upload", h)

	log.Fatal(http.ListenAndServe(":3636", nil))
}
