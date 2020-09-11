/**
* @program: go-pack
*
* @description:
*
* @author: lemo
*
* @create: 2020-09-11 13:01
**/

package pack

import (
	"os"
	"time"
)

type Pack map[string][]byte

func Unpack(pack Pack) *FileTree {
	var root = NewFileTree("", nil)

	for path := range pack {
		var mode = os.ModeDir
		if pack[path] != nil {
			mode = 0
		}
		root.AddChild(NewFileTree(path, NewFile(path, mode, time.Now(), pack[path], root)))
	}

	return root
}
