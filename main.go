package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

const PERMS fs.FileMode = 0755

const HTML_START string = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>`

const HTML_END string = `</body></html>`

// Given `input` bytes from a md file, generate the corresponding html output
// and write these contents to `fName`
func mdOutput(input []byte, fName string) {
	bfOut := blackfriday.Run(input)
	outStr := fmt.Sprintf("%s%s%s", HTML_START, bfOut, HTML_END)

	err := os.WriteFile(fName, []byte(outStr), PERMS)
	if err != nil {
		log.Fatalf("Error writing to file %s Error: %s\n", fName, err.Error())
	} else {
		log.Printf("Wrote file: %s", fName)
	}
}

// create the output directory with the name `outdir` and permissions `perms`
func mkOutDir(outdir string, perms fs.FileMode) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		err := os.Mkdir(outdir, perms)
		if err != nil {
			log.Fatalf("Error creating directory: '%s'\n", outdir)
		}

		log.Printf("Created output directory: %s\n", outdir)
	} else {
		log.Printf("Output directory '%s' exists\n", outdir)
	}

}

// return the entries of the provided `dir`
func dirents(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %s, Error: %s", dir, err.Error())
	}
	return entries
}

// Create and store generated html files in `outDir` from md files in `inDir`.
// the output directory only contains the generated html files;
// any non-md files are ignored
func createOutFiles(inDir *string, outDir *string) {
	dirents := dirents(*inDir)
	mkOutDir(*outDir, PERMS)

	for _, entry := range dirents {
		if entry.IsDir() {
			// nested directories will have the same name within the newly named
			// output directory
			subInDir := filepath.Join(*inDir, entry.Name())
			subOutDir := filepath.Join(*outDir, entry.Name())

			// recursively walk subdirs
			createOutFiles(&subInDir, &subOutDir)
		} else {
			fName := entry.Name()
			if !strings.HasSuffix(entry.Name(), ".md") {
				log.Printf("Ignoring non-markdown file: %s\n", fName)
				continue
			}
			log.Printf("Found file: %s/%s\n", *inDir, fName)
			genOutFile(entry.Name(), *inDir, *outDir)
		}
	}
}

// Generate an outfile for `fName` in `outDir`
func genOutFile(fName string, inDir string, outDir string) {
	fPath := filepath.Join(inDir, fName)
	fData, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatalf("Error reading from file: %s Error: %s\n", fName, err.Error())
	}

	htmlOutname := strings.TrimSuffix(fName, ".md")
	htmlOutname += ".html"
	outFilePath := filepath.Join(outDir, htmlOutname)
	mdOutput(fData, outFilePath)
}

func main() {
	outDir := flag.String("out-dir", "dist", "The output directory for the generated files")
	inDir := flag.String("in-dir", "inputs", "The directory containing the markdown files")
	flag.Parse()

	createOutFiles(inDir, outDir)
}
