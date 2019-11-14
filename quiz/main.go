package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const WELCOME_MESSAGE = "Welcome to Math quiz !"
const ENTER_MESSAGE = "Press [Enter] to start the quiz"
const USER_QUESTION_COUNT = 10

type Problem struct {
	question string
	answer   string
}

func (p *Problem) validateAnswer(userAnswer string) bool {
	return p.answer == userAnswer
}

func (p *Problem) ask(quiz *Quiz) {
	fmt.Println("What is ", p.question, "?")

	var ans string
	_, err := fmt.Scan(&ans)

	if err != nil {
		log.Fatalln(err)
	}

	if p.validateAnswer(ans) {
		quiz.score = quiz.score + 1
	}
}

type Quiz struct {
	problems []Problem
	score    int
}

func main() {
	filename := flag.String("file", "problems.csv", "Filename of the CSV which contains questions")
	timeout := flag.Int("timeout", 30, "Timer for the whole quiz")
	shuffle := flag.Bool("shuffle", true, "Shuffle the questions")

	flag.Parse()

	problems := mapCSVToStruct(*filename)

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

	quiz := Quiz{
		problems: problems,
		score:    0,
	}

	fmt.Println(WELCOME_MESSAGE)
	fmt.Println(ENTER_MESSAGE)

	waitForEnter()

	timer := time.NewTimer(time.Duration(*timeout) * time.Second)
	done := make(chan bool)

	go func() {
		for i := 0; i < USER_QUESTION_COUNT; i++ {
			problem := quiz.problems[i]
			problem.ask(&quiz)
		}

		done <- true
	}()

	select {
	case <-done:
		fmt.Println("You have got ", quiz.score, "/", USER_QUESTION_COUNT)
	case <-timer.C:
		fmt.Println("Timed out .....")
	}

}

func waitForEnter() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func mapCSVToStruct(filename string) []Problem {
	var qa []Problem

	csvFile, err := os.Open(filename)

	if err != nil {
		log.Fatalln("Could not open the file ", err)
	}

	reader := csv.NewReader(csvFile)

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
		}

		qa = append(qa, Problem{question: row[0], answer: row[1]})

	}

	if len(qa) == 0 {
		log.Fatalln("There were no questions in the file")
	}

	return qa
}
