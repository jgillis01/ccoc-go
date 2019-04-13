package main

import (
	"fmt"
	rec "github.com/jgillis01/ccoc_go/recording"
)

func main() {

	var recordings = map[string][]rec.Recording{}

	files := rec.GlobFiles()

	for _, file := range files {
		fmt.Printf("Processing: %s\n", file)
		r := rec.NewRecording(file)
		r.CreateFilePath()
		recordings[r.DatePath()] = append(recordings[r.DatePath()], *r)
	}

	for _, records := range recordings {
		rec.RenderTemplate(records)
	}

}
