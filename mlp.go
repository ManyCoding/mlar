// TODO
// / \
// determine if correct album folder
// custom names through arguments
// check if input dir exist
// extension in constants\array\enum
// if multiple artwork in album folder
// add statistics about removed files & scanned albums

package main

import ( 
		"os"
		"path/filepath"
		"strings"
		"log"
		"flag"
		"fmt"
)

const coverName = "folder"
const artworkDirName = "artwork"

// panic if any errors
func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// check each file or directory
func visit(path string, file os.FileInfo, err error) (e error) {
	if file.IsDir() && !IsAlbumFolder(path) {
		return filepath.SkipDir;
	}

	// format only artwork
	if IsArtwork(filepath.Ext(path)) && !strings.HasPrefix(file.Name(), coverName) {
		dir := filepath.Dir(path)
		newName := filepath.Join(dir, coverName + filepath.Ext(path))
		err := os.Rename(path, newName)
		check(err)
		fmt.Printf("Renamed \"%s\" >> \"%s\"\nin \"%s\"\n", filepath.Base(path), coverName + filepath.Ext(path), filepath.Base(filepath.Dir(path)))
	}
	return	
}

// determine if album folder
// is folder with music files and one arwork, or just folders
func IsAlbumFolder(path string) (i bool) {	
	files, _ := filepath.Glob(path + "/*.jp*g")
	files2, _ := filepath.Glob(path + "/*.png")

	files = append(files, files2...)
	
	if (len(files) > 1) {
		return false
	}
	files, _ = filepath.Glob(path + "/*.m4a")

	if files == nil {
		dir, err := os.Open(path)
		check(err)
		defer dir.Close()
		
		files, err := dir.Readdir(-1)
		check (err)
		
		for _, file := range files {
			if file.Mode().IsDir() {
				return true
			}
		}	
	}		
	return files != nil
}

func IsArtwork(ext string) (isArwork bool){
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func main() {
	flag.Parse()
	
	args := flag.Args()
	// check input argument
	if len(args) != 1 {
		log.Fatal("Only one argument (path to root) allowed")
	}
	// go traverse throught root, args[0] is path to root
	filepath.Walk(args[0], visit)
}