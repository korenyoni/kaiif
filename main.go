package main

import (
	"strings"
	"net/http"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if strings.HasSuffix(path,"/") {
		path = strings.TrimSuffix(path,"/")
	}
	_, err := AssetInfo(path)
	pathIsFile := false 
	pathIsDir := false
	var dirContents []string

	if err == nil {
		pathIsFile = true
	}
	if !pathIsFile {
		dirContents, err = AssetDir(path)
		if err == nil {
			pathIsDir = true
		}
	}

	if pathIsFile {
		log.WithFields(log.Fields{
			"path": r.URL.Path,
		}).Info("got file request")
		w.Write(MustAsset(path))
	} else if pathIsDir {
		log.WithFields(log.Fields{
			"path": r.URL.Path,
		}).Info("got path request")
		w.Write([]byte(fmt.Sprintf("%s",dirContents)))
	} else {
		log.WithFields(log.Fields{
			"path": r.URL.Path,
		}).Info("404")
		w.WriteHeader(404)
	}
}

func serveHTTP(c chan error) {
	http.HandleFunc("/", handler)
	c <- http.ListenAndServe("0.0.0.0:8080", nil)
}

func main() {
	c := make(chan error)
	go serveHTTP(c)
	err := <-c
	if err != nil {
		log.Fatal(err)
	}
}
