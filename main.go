package main

import "fmt"
import "os"
import "strings"
import "path/filepath"

//import "gopkg.in/yaml.v2"

const outputPath string = "./content"

type Recording struct {
	Title    string   `yaml:"title"`
	Date     string   `yaml:date`
	Author   string   `yaml:author`
	Tags     []string `yaml:tags`
	FileName string
}

func newFromFilePath(filePath string) *Recording {
	data := make([]byte, 200)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	count, err := file.Read(data)
	if err != nil {
		panic(err)
	}

	r := Recording{FileName: filePath}

	err = yaml.Unmarshal(data[:count], &r)
	if err != nil {
		panic(fmt.Errorf("There was an error reading file: %s with error: %v", filePath, err))
	}

	//r.FileName = filePath
	return &r
}

func (r *Recording) datePath() string {
	dateParts := strings.Split(r.Date, "-")

	if len(dateParts) != 3 {
		panic(fmt.Errorf("File does not have a properly formatted date: %s", r.Title))
	}

	datePathParts := []string{dateParts[2], dateParts[0], dateParts[1]}
	return filepath.Join(datePathParts...)
}

func (r *Recording) createFilePath() {
	current, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path := filepath.Join(current, outputPath, "recordings", r.datePath())
	fmt.Println(path)
	err = os.MkdirAll(path, os.ModePerm)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func globFiles() []string {
	files, err := filepath.Glob("./recordings/*.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return files
}

func main() {

	var recordings []Recording

	files := globFiles()

	/* convert into []map[string][]Recording{
		{
			"2019/01/01": []Recording{r1, r2}
		}
	}
	*/
	for _, file := range files {
		//r := Recording{}
		fmt.Printf("Converting file: %s\n", file)
		r := newFromFilePath(file)
		fmt.Println("Recording: ", r.FileName, "Date: ", r.datePath())
		fmt.Println("Creating output directory: ", r.FileName)
		r.createFilePath()
		recordings = append(recordings, *r)
	}

	fmt.Println("Recording Count: ", len(recordings))

}
