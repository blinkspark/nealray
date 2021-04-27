package config

import (
	"io"
	"net/http"
	"os"
)

type DisableListingFileSystem struct {
	FS http.FileSystem
}

func (dfs *DisableListingFileSystem) Open(name string) (http.File, error) {
	f, err := dfs.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return DisableListingFileSystemFile{f, 2}, nil
}

type DisableListingFileSystemFile struct {
	http.File
	batchSize int
}

func (df DisableListingFileSystemFile) Stat() (os.FileInfo, error) {
	fi, err := df.File.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
	LOOP:
		for {
			flist, err := df.File.Readdir(df.batchSize)
			switch err {
			case io.EOF:
				break LOOP
			case nil:
				for _, f := range flist {
					if f.Name() == "index.html" {
						return fi, err
					}
				}
			default:
				return nil, err
			}
		}
		return nil, os.ErrNotExist
	}
	return fi, err
}
