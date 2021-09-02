package filestore

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

// filestore.go
// a file store that stores files in a directory
// and provides restful api to access them

// FileStore is a file store that stores files in a directory
// and provides restful api to access them
type FileStore struct {
	Path string
}

// NewFileStore creates a new file store
func NewFileStore(path string) (*FileStore, error) {
	if path == "" {
		log.Panic("path is reqired")
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	return &FileStore{path}, nil
}

// ListenAndServe starts the file store server
func (fs *FileStore) Router() *mux.Router {
	r := mux.NewRouter()
	log.Println("register router")
	r.PathPrefix("/").Methods("GET").HandlerFunc(fs.GetFile)
	r.PathPrefix("/").Methods("PUT", "POST").HandlerFunc(fs.PutFile)
	return r
}

// GetFile returns the file
func (fs *FileStore) GetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFile")
	path := path.Join(fs.Path, r.URL.Path)
	_, fname := filepath.Split(path)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	io.Copy(w, f)
	// w.Write([]byte(path))
}

// PutFile puts the file
func (fs *FileStore) PutFile(w http.ResponseWriter, r *http.Request) {
	log.Println("PutFile")
	path := path.Join(fs.Path, r.URL.Path)
	// _, fname := filepath.Split(path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
	w.WriteHeader(http.StatusOK)
}
