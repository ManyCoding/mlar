// TODO
// simply divide into can rename or not
// custom names through arguments
// if multiple artwork in album folder
// / \
// display folders relatively to input

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	coverName      = "folder33333"
	artworkDirName = "artwork"
)

var visitedDirCounter = 0
var renamedFilesCounter = 0
var audioFormats = []string{".m4a", ".flac"}
var imageFormats = []string{".jpg", ".jpeg", ".png"}

// errors handling
func isError(e error, fatal bool) (isAnyError bool) {
	if e != nil {
		if fatal {
			log.Fatal(e)
		} else {
			log.Panic(e)
		}
	}
	return e != nil
}

// check if slice contains item
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// isError each file or directory
func visit(path string, file os.FileInfo, err error) (e error) {
	if file.IsDir() {
		isAlbum, message := IsAlbumFolder(path)
		fmt.Println(message)
		if isAlbum {
			visitedDirCounter++
		}
		if !IsValidToTraverse(path) {
			fmt.Println("SKIPPED " + path)
			return filepath.SkipDir
		}
	}

	// check if valid artwork
	if !contains(imageFormats, filepath.Ext(path)) || strings.HasPrefix(file.Name(), coverName+".") {
		return
	}

	dir := filepath.Dir(path)

	isAlbum, _ := IsAlbumFolder(dir)
	if !isAlbum {
		return
	}

	fmt.Println("I WAS HERE")

	newName := filepath.Join(dir, coverName+filepath.Ext(path))
	err = os.Rename(path, newName)
	if !isError(err, false) {
		renamedFilesCounter++
	}
	fmt.Printf("Renamed \"%s\" >> \"%s\"\nin \"%s\"\n", filepath.Base(path), coverName+filepath.Ext(path), filepath.Base(filepath.Dir(path)))
	return
}

// determine if we can traverse into directory
func IsValidToTraverse(path string) bool {
	isAlbum, _ := IsAlbumFolder(path)

	if (isAlbum) {
		return true
	}
	if !isAlbum {
		dir, err := os.Open(path)
		isError(err, false)
		defer dir.Close()

		files, err := dir.Readdir(-1)
		isError(err, false)

		for _, file := range files {
			if file.Mode().IsDir() {
				return true
			}
		}
	}
	return false
}

// determine if album folder
// is folder with music files and one artwork, or just folders
func IsAlbumFolder(path string) (bool, string) {
	var files []string

	// get all audio files
	for _, audioFormat := range audioFormats {
		filesBuff, _ := filepath.Glob(path + "/*" + audioFormat)
		files = append(files, filesBuff...)
	}

	if files == nil {
		return false, path + " is not album "
	}
	files = nil

	// get all files with extension from imageFormats
	for _, imageFormat := range imageFormats {
		filesBuff, _ := filepath.Glob(path + "/*" + imageFormat)
		files = append(files, filesBuff...)
	}

	if len(files) > 1 {
		return false, path + " looks like album but contains more than 1 artwork"
	}

	return true, path + " looks like album"
}

func main() {
	flag.Parse()

	args := flag.Args()
	// isError input argument
	if len(args) != 1 {
		log.Fatal("One argument (path to input directory) is compulsory")
	}
	fmt.Println("Input: " + args[0])

	// isError if dir exist
	src, err := os.Stat(args[0])
	isError(err, true)

	if !src.IsDir() {
		log.Fatal(args[0] + " is not a directory")
	}

	// go traverse throught root, args[0] is path to root
	filepath.Walk(args[0], visit)
	fmt.Printf("\nfinished: %d album\u0028s\u0029 found, %d artwork\u0028s\u0029 renamed", visitedDirCounter, renamedFilesCounter)
}
