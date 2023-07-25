package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var fileName string
	var timeLimit int
	var totalCorrect int

	flag.StringVar(&fileName, "f", "problems.csv", "Provide the csv file name")
	flag.IntVar(&timeLimit, "t", 30, "Provide the time limit in seconds")
	flag.Parse()

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	fileReader := csv.NewReader(file)
	records, err := fileReader.ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	problems := parseLines(records)

	fmt.Println("Press enter to start the quizz")
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		answerCh := make(chan string)
		go askAnswer(answerCh)

		select {
		case <-timer.C:
			fmt.Printf("\nTime is over")
			break problemloop
		case answer := <-answerCh:
			if answer == problem.answer {
				totalCorrect++
			}
		}
	}

	fmt.Printf("\nCorrect %d of %d", totalCorrect, len(records))
}

func askAnswer(ch chan<- string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	ch <- answer
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	question string
	answer   string
}
