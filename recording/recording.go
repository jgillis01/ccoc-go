package recording

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const outputPath string = "./content"

type Recording struct {
	Title    string   `yaml:"title"`
	Date     string   `yaml:date`
	Author   string   `yaml:author`
	Tags     []string `yaml:tags`
	FileName string
	FilePath string
	BaseName string
}

//Create a new Recording based on given filePath
func NewFromFilePath(filePath string) *Recording {
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
	r.BaseName = r.fileWithoutExt()

	err = yaml.Unmarshal(data[:count], &r)
	if err != nil {
		panic(fmt.Errorf("There was an error reading file: %s with error: %v", filePath, err))
	}

	return &r
}

func (r *Recording) fileWithoutExt() string {
	name := strings.TrimSuffix(r.FileName, filepath.Ext(r.FileName))
	return name
}

func (r *Recording) DatePath() string {
	dateParts := strings.Split(r.Date, "-")

	if len(dateParts) != 3 {
		panic(fmt.Errorf("File does not have a properly formatted date: %s", r.Title))
	}

	datePathParts := []string{dateParts[2], dateParts[0], dateParts[1]}
	return filepath.Join(datePathParts...)
}

func (r *Recording) CreateFilePath() {
	current, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		panic(fmt.Errorf("Could not get working directory"))
		os.Exit(1)
	}

	path := filepath.Join(current, outputPath, "recordings", r.DatePath())
	r.FilePath = path
	err = os.MkdirAll(path, os.ModePerm)

	if err != nil {
		panic(fmt.Errorf("Could not create directory at path: %s, Error: %v", r.FilePath, err))
		os.Exit(1)
	}
}

func GlobFiles() []string {
	files, err := filepath.Glob("./recordings/*.yaml")
	if err != nil {
		panic(fmt.Errorf("There was an error trying to retrieve yaml files: %v", err))
		os.Exit(1)
	}
	return files
}

func RenderTemplate(recordings []Recording) {
	filePath := filepath.Join(recordings[0].FilePath, "index.html")
	fmt.Println("Creating: ", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Errorf("There was an error creating file: %s, Error: %v", filePath, err))
	}

	defer file.Close()

	tmpl, err := template.ParseFiles("./templates/recordings.html")
	if err != nil {
		panic(fmt.Errorf("There was an error loading the template: %v", err))
	}
	err = tmpl.Execute(file, recordings)
	if err != nil {
		panic(fmt.Errorf("There was an error rendering the template: %v", err))
	}
}