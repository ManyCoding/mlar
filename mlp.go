// TODO
// custom audio/images through arguments
// display folders relatively to input
// find an album and then do renaming
// log to file
// skip showWarnings

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
var artworkFolderName string

var visitedDirCounter = 0
var renamedCoversCounter = 0
var renamedArtworkFoldersCounter = 0

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
			fmt.Println(e)
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

// visit each file and directory in the input directory
func visit(path string, file os.FileInfo, err error) (e error) {
	// parent folder for current file/folder
	parentDir := filepath.Dir(path)
	if IsAlbumFolder(path) {
		visitedDirCounter++
	}
	if file.IsDir() && isArtworkFolder(path) {
		if file.Name() == artworkFolderName {
			return
		}
		newName := filepath.Join(parentDir, artworkFolderName)
		err = os.Rename(path, newName)
		// check if renamed
		if !isError(err, -1) {
			renamedArtworkFoldersCounter++
			fmt.Printf("Renamed \"%s\" >> \"%s\"\nin \"%s\"\n", file.Name(), filepath.Base(newName), filepath.Base(filepath.Dir(path)))
		}
	} else if IsAlbumFolder(parentDir) {
		// check if valid cover
		if !contains(imageFormats, filepath.Ext(path)) || strings.HasPrefix(file.Name(), coverName+".") {
			return
		}
		newName := filepath.Join(parentDir, coverName+filepath.Ext(path))
		err = os.Rename(path, newName)
		if !isError(err, -1) {
			renamedCoversCounter++
			fmt.Printf("Renamed \"%s\" >> \"%s\"\nin \"%s\"\n", filepath.Base(path), coverName+filepath.Ext(path), filepath.Base(filepath.Dir(path)))
		}
	}

	return
}

// determine if we need to traverse into the directory
func IsValidToTraverse(path string) bool {
	if IsAlbumFolder(path) {
		return true
	}
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
	return false
}

// determine if folder is an album, e.g. contains music files
func IsAlbumFolder(path string) bool {
	var files []string

	// get all audio files
	for _, audioFormat := range audioFormats {
		filesBuffer, _ := filepath.Glob(path + "/*" + audioFormat)
		files = append(files, filesBuffer...)
	}

	if files == nil {
		return false
	}

	return true
}

// determine if artwork folder, e.g. containts only images and parent is an album
func isArtworkFolder(path string) bool {
	var files []string

	// return if parent isn't an album
	if !IsAlbumFolder(filepath.Dir(path)) {
		return false
	}

	// get all images
	for _, imageFormat := range imageFormats {
		filesBuffer, _ := filepath.Glob(path + "/*" + imageFormat)
		files = append(files, filesBuffer...)
	}

	if files == nil {
		return false
	}

	return true
}

func init() {
	const (
		defaultCoverName            = "cover"
		defaultArtworkFolderName    = "artwork"
		inputDirArgumentUsage       = "input directory (MANDATORY)"
		coverNameArgumentUsage      = "cover name"
		artworkDirNameArgumentUsage = "artwork folder name"
	)

	// define input arguments (flags)
	flag.StringVar(&inputDir, "inputDir", "", inputDirArgumentUsage)
	flag.StringVar(&inputDir, "i", "", inputDirArgumentUsage)
	flag.StringVar(&coverName, "coverName", defaultCoverName, coverNameArgumentUsage)
	flag.StringVar(&coverName, "n", defaultCoverName, coverNameArgumentUsage)
	flag.StringVar(&artworkFolderName, "artworkDirName", defaultArtworkFolderName, artworkDirNameArgumentUsage)
	flag.StringVar(&artworkFolderName, "afn", defaultArtworkFolderName, artworkDirNameArgumentUsage)
}

func main() {
	flag.Parse()

	if inputDir == "" {
		log.Fatal("Input directory argument (-i=\"path_to_music\") is compulsory")
	}

	// check if the input directory exists
	src, err := os.Stat(inputDir)
	isError(err, 1)

	// check if input is a folder
	if !src.IsDir() {
		log.Fatal(inputDir + " is not a directory")
	}

	fmt.Println("New cover name: " + coverName + ".*")
	fmt.Println("New artwork folder name: " + artworkFolderName)
	fmt.Println("Input directory: " + inputDir + "\n")

	// start traverse throught the root
	filepath.Walk(inputDir, visit)
	// print some statistics
	fmt.Printf("\nfound %d album\u0028s\u0029\nrenamed:\n%d cover\u0028s\u0029\n%d artwork folder\u0028s\u0029", visitedDirCounter, renamedCoversCounter, renamedArtworkFoldersCounter)
}
