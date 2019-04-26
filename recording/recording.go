package recording

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const outputPath string = "./content"

type Recording struct {
	Title    string   `yaml:"title" json:"title"`
	Date     string   `yaml:"date" json:"date"`
	Author   string   `yaml:"author" json:"author"`
	Tags     []string `yaml:"tags" json:"tags,omitempty"`
	HtmlPath string   `json:"htmlPath"`
	FileName string   `json:"-"`
	FilePath string   `json:"-"`
	BaseName string   `json:"-"`
}

//Create a new Recording based on given filePath
func NewRecording(filePath string) *Recording {
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

	r.HtmlPath = r.createHtmlPath()

	return &r
}

func (r *Recording) fileWithoutExt() string {
	name := strings.TrimSuffix(r.FileName, filepath.Ext(r.FileName))
	name = strings.Replace(name, "recordings/", "", 1)
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

func (r *Recording) createHtmlPath() string {
	path := filepath.Join(r.DatePath(), "index.html")
	return path
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

func RenderJsonFile(recordings []Recording) {
	file, err := os.Create(filepath.Join(outputPath, "./recordings.json"))
	if err != nil {
		panic(fmt.Errorf("There was an error opening the JSON output file: %v", err))
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	err = enc.Encode(recordings)
	if err != nil {
		panic(fmt.Errorf("There was an error encoding the JSON file: %v", err))
	}
}
