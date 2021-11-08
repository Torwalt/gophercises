package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

type Task struct {
	Question   string
	Answer     string
	UserAnswer string
}

func (t *Task) isCorrect() bool {
	return t.Answer == t.UserAnswer
}

func transformToTask(r []string) Task {
	rec := Task{
		Question: r[0],
		Answer:   r[1],
	}
	return rec
}

func askQuiz(tasks []Task, correct *int) bool {
	for _, task := range tasks {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("What is %v? ", task.Question)

		ua, _ := reader.ReadString('\n')
		ua = strings.TrimSpace(ua)
		ua = strings.ToLower(ua)

		task.UserAnswer = ua
		if task.isCorrect() {
			*correct++
		}
	}
	return true
}

func main() {
	problemPath := flag.String("problem", "problems.csv", "Path to the problem file")
	qtimer := flag.Int("timer", 30, "Quiz timer in seconds")
	shuffle := flag.Bool("shuffle", false, "Shuffle the questions")
	flag.Parse()

	records := readCsvFile(*problemPath)
	tasks := []Task{}
	for _, record := range records {
		task := transformToTask(record)
		tasks = append(tasks, task)
	}

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(tasks), func(i, j int) { tasks[i], tasks[j] = tasks[j], tasks[i] })
	}

	c1 := make(chan bool, 1)
	correct := 0

	go func() {
		done := askQuiz(tasks, &correct)
		c1 <- done
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case <-c1:
	case <-time.After(time.Duration(*qtimer) * time.Second):
		fmt.Println("\nRan out of time :(")
	}

	all := len(tasks)
	fmt.Printf("Total correct: %v, total: %v\n", correct, all)

}
