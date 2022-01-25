package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type QuestionRecord struct {
	Question string
	Answer   string
}

func createQuestionsList(data [][]string) []QuestionRecord {
	var questionlist []QuestionRecord
	for i, line := range data {
		if i > 0 { // omit header line
			var rec QuestionRecord
			for j, field := range line {
				if j == 0 {
					rec.Question = field
				} else if j == 1 {
					rec.Answer = field
				}
			}
			questionlist = append(questionlist, rec)
		}
	}
	return questionlist
}

func main() {
	var n, score int
	var fname, user_ans string
	fmt.Print("Input number of questions: ")
	fmt.Scan(&n)
	fmt.Print("Input database file name: ")
	fmt.Scan(&fname)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Failed to open csv file: %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse csv file: %s", err)
	}

	questionlist := createQuestionsList(lines)

	fmt.Printf("%+v\n", questionlist)

	if n > len(questionlist) {
		log.Fatal("Not enough questions in the database")
	} else {
		score = 0
		for i := 0; i < n; i++ {
			fmt.Println(questionlist[i].Question)
			fmt.Scan(&user_ans)
			if user_ans == questionlist[i].Answer {
				fmt.Println("Correct!")
				score++
			} else {
				fmt.Println("Incorrect!")
			}
		}
		fmt.Printf("End of quiz. You answered %v out of %v questions correctly.", score, n)
	}

}
