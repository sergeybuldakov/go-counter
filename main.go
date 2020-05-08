package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func getGoCount(url string, lockChan chan bool, countChan chan int) {
	defer func() {
		<-lockChan
	}()

	resp, err := http.Get(url)
	if err != nil {
		countChan <- 0
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		countChan <- 0
		return
	}

	defer resp.Body.Close()

	countChan <- strings.Count(string(body), "go")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	entrance := os.Stdin
	scanner := bufio.NewScanner(entrance)
	scanner.Scan()

	count := 0

	urls := strings.Split(scanner.Text(), " ")
	log.Println(urls)
	lockChan := make(chan bool, 5)
	countChan := make(chan int, len(urls))

	for _, url := range urls {
		log.Println(url)
		lockChan <- true
		go getGoCount(url, lockChan, countChan)
	}

	for i := 0; i < len(urls); i++ {
		count += <-countChan
	}

	log.Println(count)
}
