package main

import (
	"fmt"
	rec "github.com/jgillis01/ccoc-go/recording"
)

func main() {

	var recordings = map[string][]rec.Recording{}
	var allRecords = []rec.Recording{}

	files := rec.GlobFiles()

	for _, file := range files {
		fmt.Printf("Processing: %s\n", file)
		r := rec.NewRecording(file)
		r.CreateFilePath()
		recordings[r.DatePath()] = append(recordings[r.DatePath()], *r)
		allRecords = append(allRecords, *r)
	}

	for _, records := range recordings {
		fmt.Println("Creating: ", records[0].HtmlPath)
		rec.RenderTemplate(records)
	}

	// Reverse  records slice
	for i := len(allRecords)/2 - 1; i >= 0; i-- {
		opp := len(allRecords) - 1 - i
		allRecords[i], allRecords[opp] = allRecords[opp], allRecords[i]
	}

	fmt.Println("Creating:  JSON search file")
	rec.RenderJsonFile(allRecords)

	fmt.Println("Recording Count: ", len(allRecords))

}
