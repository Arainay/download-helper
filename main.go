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

// https://strm.yandex.ru/vh-kp-converted/vod-content/498426e87560e93698b9778f131ef154/9446757x1654455914x7b8f6708-3108-43ad-9223-4d8556c997d1/hls-v3/ysign1=9d5c05143f29d8f0533a56601c108fb32084c93aee93be4415ac6aa302a13ccf,abcID=1358,from=ott-kp,pfx,sfx,ts=63613d49/master_sdr_hd_avc_aac.m3u8?partner-id=NaN&video-category-id=0&imp-id=undefined&gzip=1&from=discovery&vsid=99dfc4f2bc8f958cc93205d6ee9bca03c7232c03603bxWEBx9727x1666453566&slots=436979%2C0%2C77%3B659596%2C0%2C47&testIds=436979%2C659596&session_data=1&preview=1&t=1666453566863
// https://strm.yandex.ru/vh-kp-converted/vod-content/498426e87560e93698b9778f131ef154/9446757x1654455914x7b8f6708-3108-43ad-9223-4d8556c997d1/hls-v3/ysign1=c5b95151cc60382256ce97f91e4af5c24a20f75c6a47cc197340c8dd10d06be6,abcID=1358,pfx,sfx,ts=63619ef1/video_sdr_avc_1080p_6000_audio_eng_aac_2_192.m3u8?from=discovery&chunks=1&vsid=99dfc4f2bc8f958cc93205d6ee9bca03c7232c03603bxWEBx9727x1666453566

// https://strm.yandex.ru/vh-kp-converted/vod-content/498426e87560e93698b9778f131ef154/9446757x1654455914x7b8f6708-3108-43ad-9223-4d8556c997d1/hls-v3/ysign1=655ad8342a79a5646b8efdef797709f9b1d2ce82d309018247f98b415ff5f476,abcID=1358,pfx,sfx,ts=63619ef1/video_sdr_avc_360p_365_audio_eng_aac_2_192.m3u8?from=discovery&chunks=1&vsid=99dfc4f2bc8f958cc93205d6ee9bca03c7232c03603bxWEBx9727x1666453566&redundant=48

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
