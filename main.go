package main

import "fmt"
import "os"
import "strings"
import "path/filepath"
import "gopkg.in/yaml.v2"

const outputPath string = "./content"

type Recording struct {
	Title    string   `yaml:"title"`
	Date     string   `yaml:date`
	Author   string   `yaml:author`
	Tags     []string `yaml:tags`
	FileName string
}

func (r Recording) fromFilePath(filePath string) Recording {
	data := make([]byte, 200)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	count, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data[:count], &r)
	if err != nil {
		fmt.Println("There was an error reading file: ", filePath, " with error: ", err)
		os.Exit(1)
	}

	r.FileName = filePath
	return r
}

func (r Recording) datePath() string {
	dateParts := strings.Split(r.Date, "-")

	if len(dateParts) != 3 {
		fmt.Println("File does not have a properly formatted date: ", r.Title)
		os.Exit(1)
	}

	datePathParts := []string{dateParts[2], dateParts[0], dateParts[1]}
	return filepath.Join(datePathParts...)
}

func createFilePath(r Recording) {
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
		r := Recording{}
		fmt.Printf("Converting file: %s\n", file)
		r = r.fromFilePath(file)
		fmt.Println("Recording: ", r.FileName, "Date: ", r.datePath())
		fmt.Println("Creating output directory: ", r.FileName)
		createFilePath(r)
		recordings = append(recordings, r)
	}

	fmt.Println("Recording Count: ", len(recordings))

}
