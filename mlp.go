// TODO
// rename artwork folder
// custom audio/images through arguments
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


var inputDir string
var coverName string

var visitedDirCounter = 0
var renamedFilesCounter = 0
var audioFormats = []string{".m4a", ".flac"}
var imageFormats = []string{".jpg", ".jpeg", ".png"}

// errors handling
// -1 for non-critical, 0 panic, 1 fatal
func isError(e error, status int) bool {
	if e != nil {
		if status == 1 {
			log.Fatal(e)
		} else if status == 0 {
			log.Panic(e)
		} else if status == -1 {
			log.Println(e)
		}
	}
	return e != nil
}

// check if the slice contains the item
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// find albums and artwork
func visit(path string, file os.FileInfo, err error) (e error) {
	if file.IsDir() {
		isAlbum, _ := IsAlbumFolder(path)
		if isAlbum {
			visitedDirCounter++
		}
		if !IsValidToTraverse(path) {
			//fmt.Println("SKIPPED " + path)
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

	newName := filepath.Join(dir, coverName+filepath.Ext(path))
	err = os.Rename(path, newName)
	if isError(err, -1) {
		return
	}
	renamedFilesCounter++

	fmt.Printf("Renamed \"%s\" >> \"%s\"\nin \"%s\"\n", filepath.Base(path), coverName+filepath.Ext(path), filepath.Base(filepath.Dir(path)))
	return
}

// determine if we need to traverse into the directory
func IsValidToTraverse(path string) bool {
	isAlbum, _ := IsAlbumFolder(path)

	if isAlbum {
		return true
	}
	if !isAlbum {
		dir, err := os.Open(path)
		isError(err, 0)
		defer dir.Close()

		files, err := dir.Readdir(-1)
		isError(err, 0)

		for _, file := range files {
			if file.Mode().IsDir() {
				return true
			}
		}
	}
	return false
}

// determine if folder is an album, e.g. contains music files
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

	return true, path + " looks like album"
}

func init() {
	const (
		defaultCoverName = "folder"
		coverNameArgumentUsage = "desirable artwork name"
		inputDirArgumentUsage = "input directory from which it starts traverse"
	)
	
	// define input arguments (flags)
	flag.StringVar(&coverName, "coverName", defaultCoverName, coverNameArgumentUsage)
	flag.StringVar(&coverName, "n", defaultCoverName, coverNameArgumentUsage)
	flag.StringVar(&inputDir, "inputDir", "", "input directory for traverse")
	flag.StringVar(&inputDir, "i", "", "input directory for traverse")
}

func main() {
	flag.Parse()
	
	fmt.Println("New artwork name: " + coverName + ".*")
	fmt.Println("Input directory: " + inputDir + "\n")

	if inputDir == "" {
		log.Fatal("At least one argument (input directory) is compulsory")
	}
	
	// check if the input directory exists
	src, err := os.Stat(inputDir)
	isError(err, 1)

	// check if input is a folder
	if !src.IsDir() {
		log.Fatal(inputDir + " is not a directory")
	}

	// start traverse throught the root
	filepath.Walk(inputDir, visit)
	// print some statistics
	fmt.Printf("\nfinished: %d album\u0028s\u0029 found, %d artwork\u0028s\u0029 renamed", visitedDirCounter, renamedFilesCounter)
}
