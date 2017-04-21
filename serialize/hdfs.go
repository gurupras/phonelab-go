package serialize

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gurupras/go-easyfiles"
	"github.com/gurupras/go-easyfiles/easyhdfs"
)

type HDFSSerializer struct {
	Addr string
}

func NewHDFSSerializer(addr string) *HDFSSerializer {
	return &HDFSSerializer{
		Addr: stripHDFSPrefix(addr),
	}
}

func stripHDFSPrefix(addr string) string {
	if strings.HasPrefix(addr, "hdfs://") {
		addr = addr[7:]
	}
	return addr
}

func (h *HDFSSerializer) Serialize(obj interface{}, filename string) error {
	// FIXME: We should use a pool of connections
	// This will blow up the number of connections if there are a large
	// number of goroutines.
	fs := easyhdfs.NewHDFSFileSystem(h.Addr)

	filename = stripHDFSPrefix(filename)
	fileType := easyfiles.GZ_FALSE
	if strings.HasSuffix(filename, ".gz") {
		fileType = easyfiles.GZ_TRUE
	}

	//Mkdirs
	outdir := path.Dir(filename)
	err := fs.Makedirs(outdir)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %v: %v", outdir, err)
	}

	file, err := fs.Open(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileType)
	if err != nil {
		return fmt.Errorf("Failed to open file: %v: %v", filename, err)
	}
	defer file.Close()

	writer, err := file.Writer(0)
	if err != nil {
		return fmt.Errorf("Failed to get writer to file: %v: %v", filename, err)
	}
	defer writer.Close()
	defer writer.Flush()

	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("Failed to marshal object to json: %v", err)
	}

	if _, err := writer.Write(b); err != nil {
		return fmt.Errorf("Failed to write json bytes to file: %v: %v", filename, err)
	}
	return nil
}
