package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func main() {

	subDirToSkip := ".git"

	err := filepath.Walk("..", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		if info.IsDir() {
			fmt.Printf("Visited Directory: %s\n", strings.Split(path, "/")[1:])
		} else {
			fmt.Printf("Visited File: %s\n", strings.Split(path, "/")[1:])
		}

		return nil
	})
	
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
	}
}
