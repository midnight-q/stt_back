package webapp

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"stt_back/services/file_storage"
)

func StaticFileLoader(w http.ResponseWriter, httpRequest *http.Request) {
	vars := mux.Vars(httpRequest)
	name := vars["name"]
	folder := vars["folder"]
	data, _ := file_storage.LoadFile(name, folder)
	file_ext := filepath.Ext(name)
	if file_ext == ".wav"{
		w.Header().Set("Content-Type", "audio/wav")
	}
	if file_ext == ".pdf"{
		w.Header().Set("Content-Type", "application/pdf")
	}
	if file_ext == ".docx"{
		w.Header().Set("Content-Type", "application/docx")
	}
	if file_ext == ".txt"{
		w.Header().Set("Content-Type", "text")
	}
	if file_ext == ".html"{
		w.Header().Set("Content-Type", "text/html")
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	if _, err := w.Write(data); err != nil {
		log.Println("unable to write file")
	}
}

