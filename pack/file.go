/**
* @program: go-pack
*
* @description:
*
* @author: lemo
*
* @create: 2020-09-10 20:05
**/

package pack

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func NewFile(path string, fileMode os.FileMode, mtime time.Time, data []byte, tree *FileTree) *File {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &File{
		path:     absPath,
		fileMode: fileMode,
		mtime:    mtime,
		data:     data,
		tree:     tree,
	}
}

// An asset file.
type File struct {
	// The full asset file path
	path string

	// The asset file mode
	fileMode os.FileMode

	// The asset modification time
	mtime time.Time

	// The asset data. Note that this data might be in gzip compressed form.
	data []byte

	buf      *bytes.Reader
	tree     *FileTree
	index    int
	children *FileTree
}

func (f *File) Path() string {
	return f.path
}

// Implementation of os.FileInfo
func (f *File) Name() string {
	return path.Base(f.path)
}

func (f *File) Mode() os.FileMode {
	return f.fileMode
}

func (f *File) ModTime() time.Time {
	return f.mtime
}

func (f *File) IsDir() bool {
	return f.fileMode.IsDir()
}

func (f *File) Size() int64 {
	return int64(len(f.data))
}

func (f *File) Sys() interface{} {
	return nil
}

// Implementation of http.File
func (f *File) Close() error {
	f.buf = nil
	f.index = 0
	f.children = nil
	return nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if f.IsDir() {
		var res []os.FileInfo
		if f.children == nil {
			f.children = f.tree.FindChild(f.path)
		}

		if count < 0 {
			if f.children == nil {
				return nil, nil
			}

			for i := 0; i < len(f.children.children); i++ {
				res = append(res, f.children.children[i].file)
			}
			return res, nil
		}

		if f.children == nil {
			return nil, io.EOF
		}

		var index = f.index + count

		if index > len(f.children.children) {
			index = len(f.children.children)
		}

		for i := f.index; i < index; i++ {
			res = append(res, f.children.children[i].file)
		}

		f.index += count

		if f.index == len(f.children.children) {
			return res, io.EOF
		}

		return res, nil

	} else {
		return nil, os.ErrInvalid
	}
}

func (f *File) Read(data []byte) (int, error) {
	if f.buf == nil {
		f.buf = bytes.NewReader(f.data)
	}
	return f.buf.Read(data)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.buf == nil {
		f.buf = bytes.NewReader(f.data)
	}
	return f.buf.Seek(offset, whence)
}
