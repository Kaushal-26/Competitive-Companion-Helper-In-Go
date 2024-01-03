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

// Error Handling for the code
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// Make Files
func makeNewFile(Path string) {
	fmt.Println("Making a New File: ", Path)
	f, err := os.Create(Path)
	checkError(err)
	defer f.Close()
}

// Checking if a directory exists or Making it
func checkDirExistsOrMakeNewDir(Path string) {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Making a new Directory: ", Path)
			err := os.Mkdir(Path, 0755)
			checkError(err)
		} else {
			panic(err)
		}
	}
	os.Chdir(Path)
}

func makeFilesInSystem(data Data, templateFile []byte, LANGUAGE, STDIN, STDOUT string) {
	// Spliting data as no spaces reqd
	Site := strings.Split(data.Group, " ")

	// Making a directory for the site of contest
	Path := Site[0] + "/"
	checkDirExistsOrMakeNewDir(Path)

	// Directory for with the contest URL
	Site = strings.Split(data.URL, "/")
	checkDirExistsOrMakeNewDir(Site[len(Site)-3])

	// Code File with the Problem number like A, B, C1, C2, D... or 1, 2, 3 as per contest
	codeFile := Site[len(Site)-1] + LANGUAGE
	makeNewFile(codeFile)

	// Writing our template on to the code file
	if len(templateFile) > 0 {
		err := os.WriteFile(codeFile, templateFile, 0644)
		checkError(err)
	}

	for index, test := range data.Tests {
		// Index for the testcases as 1, 2, 3...
		indexToString := strconv.Itoa(index + 1)

		// Making input file if STDIN exists
		if len(STDIN) > 0 {
			inputFile := indexToString + STDIN
			makeNewFile(inputFile)
			err := os.WriteFile(inputFile, []byte(test.Input), 0644)
			checkError(err)
		}

		// Making output file if STDOUT exists
		if len(STDOUT) > 0 {
			outputFile := indexToString + STDOUT
			makeNewFile(outputFile)
			err := os.WriteFile(outputFile, []byte(test.Output), 0644)
			checkError(err)
		}
	}
}

func main() {
	// Your .env file where all the varibles are defined
	current_working_directory, err := os.Getwd()
	checkError(err)
	godotenv.Load(current_working_directory + "\\.env")

	// Environmental variables which you will take
	var PORT, FILE_DIRECTORY_PATH, LANGUAGE, STDIN, STDOUT string

	PORT = os.Getenv("PORT")
	FILE_DIRECTORY_PATH = os.Getenv("FILE_DIRECTORY_PATH")
	LANGUAGE = os.Getenv("LANGUAGE")
	STDIN = os.Getenv("STDIN")   // Not Compulsory
	STDOUT = os.Getenv("STDOUT") // Not Compulsory

	if PORT == "" || FILE_DIRECTORY_PATH == "" || LANGUAGE == "" {
		fmt.Println("Environmental Variables Not Defined!")
		return
	}

	var templateFile []byte
	// Taking Argument to make code file with template in whatever language we want
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

		// File path to the directory we want to make codes
		os.Chdir(FILE_DIRECTORY_PATH)

		// Function to make the files on FILE_DIRECTORY_PATH
		makeFilesInSystem(data, templateFile, LANGUAGE, STDIN, STDOUT)
	})

	// Starting a http server
	router.Run(PORT)
}
