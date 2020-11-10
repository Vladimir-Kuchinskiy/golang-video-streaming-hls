package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexPage).Methods("GET")
	router.HandleFunc("/favicon.ico", http.NotFound).Methods("GET")
	router.HandleFunc("/media/{mid:[0-9]+}/stream/", streamHandler).Methods("GET")
	router.HandleFunc("/media/{mid:[0-9]+}/stream/{segName:index[0-9]+.ts}", streamHandler).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mid, err := strconv.Atoi(vars["mid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	segName, ok := vars["segName"]
	if !ok {
		mediaBase := getMediaBase(mid)
		m3u8Name := "index.m3u8"
		serveHLSm3u8(w, r, mediaBase, m3u8Name)

		return
	}

	mediaBase := getMediaBase(mid)
	serveHLSts(w, r, mediaBase, segName)
}

func getMediaBase(mid int) string {
	mediaRoot := "assets/media"
	return fmt.Sprintf("%s/%d", mediaRoot, mid)
}

func serveHLSm3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, m3u8Name)
	fmt.Println(mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func serveHLSts(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, segName)
	fmt.Println(mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}
