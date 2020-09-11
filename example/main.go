/**
* @program: go-pack
*
* @description:
*
* @author: lemo
*
* @create: 2020-09-11 13:06
**/

package main

import (
	"log"
	"net/http"

	"github.com/lemoyxk/go-pack/pack"
)

func main() {

	var root = pack.Unpack(Load())

	// root.AddChild(unpack.NewFileTree("/index.html", unpack.NewFile("/index.html", 0, time.Now(), []byte("hello"), root)))

	root.Walk(func(path string, name string, file *pack.File) {
		log.Println(path, name)
	})

	var fileSystem = &pack.FileSystem{
		// FileTree: root.FindChild("/Users/lemo/Downloads/lemo-1.0.0/console"),
		FileTree: root,
		// Dir:      "/Users/lemo/Downloads/lemo-1.0.0/console",
		Dir:   "/Users/lemo",
		Debug: false,
	}

	panic(http.ListenAndServe(":12345", http.FileServer(fileSystem)))
}
