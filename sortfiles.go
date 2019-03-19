package main

import (
    "fmt"
    "io/ioutil"
	"log"
	"strings"
	"os"
	"io"
	"time"
)

func directoryDoesntExists(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return true
	}
	return false
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return
    }
    err = out.Sync()
    return
}

func main() {
	start := time.Now()
	files_copied := 0
	directory := os.Args[1]
	destinationbase := os.Args[2] //"/mnt/f/Fotos/"
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println("Found ", len(files), "files")
	fmt.Println("Copying...")
	dirs := make(map[string]bool)
    for _, f := range files {
		date := strings.Split(f.Name(), " ")
		splitdate := strings.Split(date[0], "-")
		destination := destinationbase + splitdate[0] + "/" + splitdate[0] + splitdate[1]
		if !dirs[destination] && directoryDoesntExists(destination) {
			errDir := os.MkdirAll(destination, 0755)
			if errDir != nil {
				panic(err)
			}
			dirs[destination] = true
		}
		copyFileContents(directory + "/" + f.Name(), destination + "/" + f.Name())
		files_copied++
	}
	elapsed := time.Since(start)
	fmt.Println("Procedure took: ", elapsed)
	fmt.Println(files_copied, " files copied to ", len(dirs), " directories")
}