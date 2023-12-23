package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Json to Struct we want to make use of the following things
type Data struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	URL         string `json:"url"`
	MemoryLimit int    `json:"memoryLimit"`
	TimeLimit   int    `json:"timeLimit"`
	Tests       []struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"tests"`
}

// Environmental varibles which you can set
var PORT, FILE_DIRECTORY_PATH, LANGUAGE, STDIN, STDOUT, STDERR string

// Error Handling for the code
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// Make Files
func makeNewFile(Path string) {
	f, err := os.Create(Path)
	checkError(err)
	defer f.Close()
}

// Checking if a directory exists or Making it
func checkDirExistsOrMakeNewDir(Path string) {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(Path, 0755)
			checkError(err)
		} else {
			panic(err)
		}
	}
	os.Chdir(Path)
}

func makeFilesInSystem(data Data, templateFile []byte) {
	// Spliting data as no spaces reqd
	Site := strings.Split(data.Group, " ")

	// Making a directory for the site of contest
	Path := Site[0] + "/"
	checkDirExistsOrMakeNewDir(Path)

	// Directory for with the contest name
	Path = strings.Join(Site[2:], " ")
	checkDirExistsOrMakeNewDir(Path)

	// Directory for the Problem name
	checkDirExistsOrMakeNewDir(data.Name)

	// Now Splitting data for name of the problem
	Site = strings.Split(data.Name, " ")

	// Code File with the Problem number like A, B, C1, C2, D... or 1, 2, 3 as per contest
	codeFile := Site[0][:len(Site[0])-1] + LANGUAGE
	makeNewFile(codeFile)

	// Writing our template on to the code file
	err := os.WriteFile(codeFile, templateFile, 0644)
	checkError(err)

	for index, test := range data.Tests {
		// Index for the testcases as 1, 2, 3...
		indexToString := strconv.Itoa(index + 1)

		// Making input file
		inputFile := indexToString + STDIN
		makeNewFile(inputFile)
		err := os.WriteFile(inputFile, []byte(test.Input), 0644)
		checkError(err)

		// Making output file
		if len(STDOUT) > 0 {
			outputFile := indexToString + STDOUT
			makeNewFile(outputFile)
			err = os.WriteFile(outputFile, []byte(test.Output), 0644)
			checkError(err)
		}

		// Making Error file for the Programming file
		if len(STDERR) > 0 {
			errorFile := indexToString + STDERR
			makeNewFile(errorFile)
		}
	}
}

func main() {
	godotenv.Load(".env")

	PORT = os.Getenv("PORT")
	FILE_DIRECTORY_PATH = os.Getenv("FILE_DIRECTORY_PATH")
	LANGUAGE = os.Getenv("LANGUAGE")
	STDIN = os.Getenv("STDIN")
	STDOUT = os.Getenv("STDOUT")
	STDERR = os.Getenv("STDERR")

	// Taking Argument to make code file with template in whatever language we want
	var err error
	var templateFile []byte
	if len(os.Args) > 1 {
		argsOfProgram := os.Args[1]
		templateFile, err = os.ReadFile(argsOfProgram)
		checkError(err)
	}

	// Starting a gin http router
	router := gin.Default()

	// GET function for the endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"PORT": "STARTED",
		})
	})

	// POST function for the endpoint
	router.POST("/", func(c *gin.Context) {
		// Data Send to the Server
		// Decoding the Json File using Unmarshall and storing it in a struct
		data := Data{}
		err = c.BindJSON(&data)
		checkError(err)

		fmt.Println(data)

		// Getting the current directory to resign the current path after making codes file
		owd, err := os.Getwd()
		checkError(err)

		// File path to the directory we want to make codes
		os.Chdir(FILE_DIRECTORY_PATH)

		// Function to make the files on FILE_DIRECTORY_PATH
		makeFilesInSystem(data, templateFile)

		// Changing our directory back to where we started
		os.Chdir(owd)
	})

	// Starting a http server
	router.Run(PORT)
}
