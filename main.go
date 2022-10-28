package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	action := chooseAction()

	for action == "d" || action == "r" {
		if action == "d" || action == "r" {
			start(action)
		}
		action = chooseAction()
	}

	os.Exit(0)
}

func start(action string) {
	isDownload := action == "d"

	if isDownload {
		var filename string

		fmt.Println("File name")
		fmt.Scan(&filename)

		file, err := os.Create(filename + ".mp4")

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		readFile(func(line string) {
			if strings.HasPrefix(line, "http") {
				part, err := downloadFilePart(line)

				if err != nil {
					log.Fatal(err)
				}

				if _, err := file.Write(part); err != nil {
					log.Fatal(err)
				}
			}
		})
	} else {
		readFile(func(line string) {
			fmt.Println(line)
		})
	}
}

func chooseAction() string {
	var action string

	fmt.Println("d - загрузить")
	fmt.Println("r - показать содержимое файла")
	fmt.Println("q - выйти")

	fmt.Scan(&action)

	return action
}

func readFile(callback func(line string)) {
	var playlistUrl string

	fmt.Println("Введите адрес файла *.m3u8")
	fmt.Scan(&playlistUrl)

	response, err := http.Get(playlistUrl)

	if err != nil {
		log.Fatal("Download error: ", err)
	}

	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "http") && !strings.HasPrefix(line, "#") {
			parsedUrl, err := url.Parse(playlistUrl)

			if err != nil {
				log.Fatal(err)
			}

			line = parsedUrl.Scheme + "://" + parsedUrl.Host + line
		}

		callback(line)
	}

	if scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func downloadFilePart(url string) ([]byte, error) {
	fmt.Printf("Downloading %s...\n", url)

	result := make([]byte, 0)
	response, err := http.Get(url)

	if err != nil {
		return result, err
	}

	return ioutil.ReadAll(response.Body)
}
