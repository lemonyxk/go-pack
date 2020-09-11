/**
* @program: go-pack
*
* @description:
*
* @author: lemo
*
* @create: 2020-09-11 13:48
**/

package pack

import (
	"bytes"
	"net/http"
	"os"
	"path"
)

type FileSystem struct {
	FileTree *FileTree
	Dir      string
	Debug    bool
}

func (f *FileSystem) Open(p string) (http.File, error) {
	p = path.Clean(p)

	if f.Dir == "" {
		f.Dir = "/"
	}

	if f.Debug {
		return http.Dir(f.Dir).Open(p)
	}

	var r = f.FileTree.FindChild(path.Join(f.Dir, p))

	if r == nil {
		return nil, os.ErrNotExist
	}

	if !r.file.IsDir() {
		// Make a copy for reading
		r.file.buf = bytes.NewReader(r.file.data)
		return r.file, nil
	}

	return r.file, nil
}
