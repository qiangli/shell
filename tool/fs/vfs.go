package fs

import (
	"fmt"
	"io"
	"log"
	"strings"

	// "github.com/c2fo/vfs/v7"
	"github.com/c2fo/vfs/v7/vfssimple"
)

func VFS() {
	// Create local OS file from a URI
	osFile, err := vfssimple.NewFile("file:///tmp/example.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer osFile.Close()

	// Write to the file
	_, err = io.Copy(osFile, strings.NewReader("Hello from vfs!"))
	if err != nil {
		log.Fatal(err)
	}

	if err := osFile.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("File created and written:", osFile.URI())
}
