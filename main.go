package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type QuestionAnswer struct {
	Question string
	Answer   string
}

func main() {
	filename := flag.String("filename", "problems.csv", "Specify filename of csv file containing quiz questions")
	gameDuration := flag.Int("gameDuration", 30, "The duration the player has to answer all questions")
	flag.Parse()

	questions, err := readQuestions(*filename)
	if err != nil {
		return
	}
	run(questions, *gameDuration)
}

func readQuestions(filename string) ([]QuestionAnswer, error) {
	fmt.Println("Loading Questions from file: ", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("File cannot be opened")
	}
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, errors.New("File cannot be read")
	}
	var questions []QuestionAnswer
	for _, record := range records {
		questions = append(questions, QuestionAnswer{record[0], record[1]})
	}
	defer file.Close()
	return questions, nil
}

func run(questions []QuestionAnswer, duration int) {
	correct, incorrect := 0, 0
	reader := bufio.NewReader(os.Stdin)

	timer := time.NewTimer(time.Duration(duration) * time.Second)
	defer timer.Stop()
	fmt.Printf("You have %d seconds to complete all questions\n", duration)

QuestionLoop:
	for i, q := range questions {
		fmt.Printf("Question %d: %s\n", i+1, q.Question)
		answerChannel := make(chan string)

		go func() {
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			answerChannel <- text
		}()

		select {
		case <-timer.C:
			fmt.Println("Times up! Game Over!")
			break QuestionLoop
		case text := <-answerChannel:
			if text == q.Answer {
				correct++
			} else {
				incorrect++
			}
		}
	}
	fmt.Printf("Correct: %d  Incorrect: %d\n", correct, incorrect)
}
