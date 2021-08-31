package filestore

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path"
	"strings"

	ncore "github.com/blinkspark/prototypes/network_core"
	"github.com/libp2p/go-libp2p-core/network"
)

const (
	FILESTORE_PROTOCOL = "nealfree.cf/filestore/v0.1"
)

type FSFileStore struct {
	*ncore.NetworkCore
	StorePath string
	KeyPath   string
}

// handleFileStoreStream handles a stream for the filestore protocol
// first byte is the command
// 0x00 - get file
// 0x01 - put file
// 0x02 - delete file
// 0x03 - list files
// the next 4 bytes is filename length or path length
// if we got error return a message starts with 0xff
// if no error return a message starts with 0x00
func (fstore *FSFileStore) handleFileStoreStream(s network.Stream) {
	log.Println("handleFileStoreStream")
	defer s.Reset()
	buf := make([]byte, 1)
	_, err := s.Read(buf)
	if err != nil {
		return
	}
	switch buf[0] {
	case 0x00:
		fstore.handleGetFile(s)
	case 0x01:
		fstore.handlePutFile(s)
	case 0x02:
		fstore.handleDeleteFile(s)
	case 0x03:
		fstore.handleListFiles(s)
	default:
		s.Write([]byte{0xff})
	}
}

// handleGetFile handles a get file request
func (fstore *FSFileStore) handleGetFile(s network.Stream) {
	// TODO fix this
	buf := make([]byte, 4)
	_, err := s.Read(buf)
	if err != nil {
		return
	}
	filename := string(buf)
	file, err := fstore.GetFile(filename)
	if err != nil {
		s.Write([]byte{0xff})
		return
	}
	defer file.Close()

	s.Write([]byte{0x00})
	// make a buffer to read the file
	buf = make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		s.Write(buf[:n])
	}
}

// handlePutFile handles a put file request
func (fstore *FSFileStore) handlePutFile(s network.Stream) {
	log.Println("handlePutFile")
	// make a buffer to read stream
	buf := make([]byte, 1024)
	// read first 4 bytes for filename length
	_, err := s.Read(buf[:4])
	if err != nil {
		s.Write([]byte{0xff})
		return
	}
	filenameLength := binary.BigEndian.Uint32(buf[:4])
	// read filename
	_, err = s.Read(buf[:filenameLength])
	if err != nil {
		s.Write([]byte{0xff})
		return
	}
	filename := string(buf[:filenameLength])
	// write file
	fstore.PutFile(filename, s)
}

// PutFile puts a file to the store
func (fstore *FSFileStore) PutFile(filename string, s io.ReadWriter) {
	log.Println("PutFile")
	fpath := path.Join(fstore.StorePath, filename)
	file, err := os.Create(fpath)
	if err != nil {
		s.Write([]byte{0xff})
		return
	}
	defer file.Close()

	// make a buffer to read stream
	buf := make([]byte, 1024)
	for {
		n, err := s.Read(buf)
		if err != nil {
			break
		}
		file.Write(buf[:n])
	}
}

// GetFile gets a file from the store
func (fstore *FSFileStore) GetFile(filename string) (io.ReadCloser, error) {
	fpath := path.Join(fstore.StorePath, filename)
	return os.Open(fpath)
}

// handleListFiles handles a list files request
func (fstore *FSFileStore) handleListFiles(s network.Stream) {
	// read first 4 bytes for path name length
	buf := make([]byte, 4)
	_, err := s.Read(buf)
	if err != nil {
		return
	}
	pathname := string(buf)
	// invoke list files
	filena