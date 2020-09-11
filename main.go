/**
* @program: go-pack
*
* @description:
*
* @author: lemo
*
* @create: 2020-09-10 19:22
**/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%v", err)
		}
	}()

	src := flag.String("src", "", "pack source")
	packageName := flag.String("package", "main", "package name")
	output := flag.String("output", "", "output go file")

	flag.Parse()

	if *src == "" {
		panic("src is empty")
	}

	if *output == "" {
		panic("output is empty")
	}

	create(*src, *packageName, *output)
}

func create(src, packageName, output string) {

	absSrc, err := filepath.Abs(src)
	if err != nil {
		panic(err)
	}

	absOutput, err := filepath.Abs(output)
	if err != nil {
		panic(err)
	}

	var temp = `package ` + packageName + "\n"

	temp += `
func Load() map[string][]byte {
	return filesMap
}

var filesMap = map[string][]byte{
`

	var tempBts = bytes.NewBufferString(temp)

	_ = filepath.Walk(absSrc, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			tempBts.Write([]byte("\t" + `"` + path + `"` + ":" + "nil" + ",\n"))
		} else {
			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}

			var s = bytes.NewBuffer([]byte("[]byte{"))
			for i := 0; i < len(b); i++ {
				if i == len(b)-1 {
					s.Write([]byte(strconv.Itoa(int(b[i]))))
				} else {
					s.Write([]byte(strconv.Itoa(int(b[i])) + ","))
				}
			}
			s.Write([]byte("\t}"))

			tempBts.Write([]byte("\t" + `"` + path + `"` + ":" + s.String() + ",\n"))
		}

		return nil
	})

	tempBts.Write([]byte("}\n"))

	write(absOutput, tempBts)
}

func write(path string, buf *bytes.Buffer) {

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer func() { _ = f.Close() }()

	_, err = f.WriteString(buf.String())
	if err != nil {
		panic(err)
	}

}
