package webapp

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"stt_back/services/file_storage"
)

func StaticFileLoader(w http.ResponseWriter, httpRequest *http.Request) {
	vars := mux.Vars(httpRequest)
	name := vars["name"]
	folder := vars["folder"]
	data, _ := file_storage.LoadFile(name, folder)
	w.Header().Set("Content-Type", "audio/wav")

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	if _, err := w.Write(data); err != nil {
		log.Println("unable to write file")
	}
}

