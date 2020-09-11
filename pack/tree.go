/**
* @program: go-pack
*
* @description:
*
* @author: lemo
*
* @create: 2020-09-10 19:29
**/

package pack

import (
	"os"
	path2 "path"
	"path/filepath"
	"strings"
	"time"
)

type FileTree struct {
	name     string
	path     string
	file     *File
	children []*FileTree
	parent   *FileTree
}

func (f *FileTree) Parent() *FileTree {
	return f.parent
}

func (f *FileTree) Children() []*FileTree {
	return f.children
}

func (f *FileTree) File() *File {
	return f.file
}

func (f *FileTree) Path() string {
	return f.path
}

func (f *FileTree) Name() string {
	return f.name
}

func NewFileTree(path string, file *File) *FileTree {

	if path == "" {
		return &FileTree{path: "", file: nil, name: ""}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return &FileTree{path: absPath, file: file, name: path2.Base(path)}
}

func (f *FileTree) FindChild(path string) *FileTree {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	var pathSlice = getPathSlice(f.path, absPath)

	if pathSlice == nil {
		return f
	}

	var r = f

	for i := 0; i < len(pathSlice); i++ {

		var child = getChildByPath(pathSlice[i], r)

		if child == nil {
			return nil
		}

		if i == len(pathSlice)-1 {
			return child
		} else {
			r = child
		}
	}

	return nil
}

func (f *FileTree) RemoveChild(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	var self = f.FindChild(absPath)

	if self == nil {
		return
	}

	var name = path2.Base(absPath)

	var index = -1
	for i := 0; i < len(self.parent.children); i++ {
		if self.parent.children[i].name == name {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	self.parent.children = append(self.parent.children[0:index], self.parent.children[index+1:]...)
}

func (f *FileTree) AddChild(fileTree *FileTree) {
	if f == nil {
		panic("root is empty")
	}

	var pathSlice = getPathSlice(f.path, fileTree.path)

	var r = f

	for i := 0; i < len(pathSlice); i++ {
		var p = path2.Join(pathSlice[:i+1]...)
		var child = getChildByPath(pathSlice[i], r)

		if child == nil {
			// last
			if i == len(pathSlice)-1 {
				r.children = append(r.children, fileTree)
				fileTree.parent = r
				// reset
				r = fileTree
			} else {
				var newFileTree = NewFileTree(p, NewFile(p, os.ModeDir, time.Now(), []byte(p), f))
				newFileTree.parent = r
				r.children = append(r.children, newFileTree)
				// reset
				r = newFileTree
			}
		}

		if child != nil {
			if i == len(pathSlice)-1 {
				fileTree.parent = r
				child = fileTree
			} else {
				r = child
			}
		}
	}
}

func getPathSlice(rootPath, path string) []string {
	var pathSlice = strings.Split(path, "/")

	if rootPath == "" {
		// need find from /
		pathSlice[0] = "/"
	} else {
		// means itself
		if path == "/" {
			return nil
		}
		// find without /
		pathSlice = pathSlice[1:]
	}

	if pathSlice[len(pathSlice)-1] == "" {
		pathSlice = pathSlice[:len(pathSlice)-1]
	}

	return pathSlice
}

func getChildByPath(name string, r *FileTree) *FileTree {
	for i := 0; i < len(r.children); i++ {
		if r.children[i].name == name {
			return r.children[i]
		}
	}
	return nil
}

func (f *FileTree) Walk(fn func(path string, name string, file *File)) {

	if f == nil {
		return
	}

	fn(f.path, f.name, f.file)

	for i := 0; i < len(f.children); i++ {
		if len(f.children[i].children) != 0 {
			f.children[i].Walk(fn)
		} else {
			fn(f.children[i].path, f.children[i].name, f.children[i].file)
		}
	}

}
