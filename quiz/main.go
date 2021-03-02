package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// sumQuiz returns how many score you reach in the quiz
func sumQuiz(readFile *csv.Reader) {
	scanner := bufio.NewScanner(os.Stdin)
	rightAns := 0 //right answers
	for {
		record, err := readFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) != 2 {
			log.Fatal("Wrong format")
		}
		fmt.Printf("%s = ", record[0])

		var val int
		//fmt.Scanf("%d", &val)
		scanner.Scan()
		text := scanner.Text()
		result := strings.Trim(record[1], " ")
		csvResult, _ := strconv.Atoi(result)
		val, errVal := strconv.Atoi(text)
		if errVal != nil {
			log.Fatal(errVal)
		}

		if val == csvResult {
			rightAns++
		} else {
			readRemains, errRR := readFile.ReadAll()
			if errRR != nil {
				log.Fatal(errRR)
			}
			totalQuizs := rightAns + len(readRemains) + 1 //plus current record
			fmt.Printf("You scored %d out of %d", rightAns, totalQuizs)
			break
		}
	}
}

//Problem structure
type Problem struct {
	question, answer string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))
	for i, line := range lines {
		ret[i] = Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func main() {

	timeLimit := flag.Int("limit", 3, "the time limit for the quiz")
	flag.Parse()

	path := ""

	fmt.Printf("Give me path csv file: ")
	fmt.Scan(&path)
	fd, err := os.Open(path)
	if err != nil {
		exit(fmt.Sprintf("Fail to open the CSV file: %s\n", path))
		os.Exit(1)
	}
	file := csv.NewReader(fd)

	lines, errLines := file.ReadAll()
	if errLines != nil {
		exit("fail to parse the provided csv file")
	}

	problems := parseLines(lines)
	fmt.Println(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	score := 0
problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
		answerCh := make(chan string)
		go func() {
			var inputAns string
			fmt.Scan(&inputAns)
			answerCh <- inputAns
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nExceed time!\n")
			break problemloop
		case answer := <-answerCh:
			if answer == problem.answer {
				score++
				fmt.Println("Correct!")
			} else {
				break problemloop
			}

		}

	}
	fmt.Printf("You scored %d out of %d problems", score, len(lines))

}
