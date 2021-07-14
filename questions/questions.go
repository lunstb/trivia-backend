package questions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Categories struct {
	Categories []*Category `json:"categories"`
}

type Question struct {
	Question string `json:"question"`
	Answer   int    `json:"answer"`
	Unit     string `json:"unit"`
	Source   string `json:"source"`
}

type Category struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Questions   []*Question `json:"questions"`
}

func GetCategories() Categories {
	questionsFile, err := os.Open("./questions/questions.json")

	if err != nil {
		fmt.Println(err)
	}

	defer questionsFile.Close()

	byteValue, _ := ioutil.ReadAll(questionsFile)

	var categories Categories

	json.Unmarshal(byteValue, &categories)

	return categories
}
