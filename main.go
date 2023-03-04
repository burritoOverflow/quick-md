package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
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

func mdOutput(input []byte, fName string) {
	bfOut := blackfriday.Run(input)
	outStr := fmt.Sprintf("%s%s%s", HTML_START, bfOut, HTML_END)

	err := os.WriteFile(fName, []byte(outStr), PERMS)
	if err != nil {
		log.Fatalf("Error writing to file %s Error: %s\n", fName, err.Error())
	}
}

func mkOutDir(outdir string) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		err := os.Mkdir(outdir, PERMS)
		if err != nil {
			log.Fatalf("Error creating directory: '%s'\n", outdir)
		}

		log.Printf("Created output directory: %s\n", outdir)
	} else {
		log.Printf("Output directory '%s' exists\n", outdir)
	}

}

func main() {
	outDir := flag.String("out-dir", "dist", "The output directory for the generated files")
	inDir := flag.String("in-dir", "inputs", "The directory containing the markdown files")
	flag.Parse()

	mkOutDir(*outDir)
	inFiles, err := os.ReadDir(*inDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range inFiles {
		if !strings.HasSuffix(f.Name(), ".md") {
			log.Printf("Ignoring file: %s\n", f.Name())
			continue
		}

		log.Printf("Found file: %s\n", f.Name())
		fPath := fmt.Sprintf("%s/%s", *inDir, f.Name())
		fData, err := os.ReadFile(fPath)
		if err != nil {
			log.Fatalf("Error reading from file %s Error: %s\n", f.Name(), err.Error())
		}

		htmlOutname := strings.TrimSuffix(f.Name(), ".md")
		htmlOutname += ".html"
		mdOutput(fData, fmt.Sprintf("%s/%s", *outDir, htmlOutname))
	}

}
