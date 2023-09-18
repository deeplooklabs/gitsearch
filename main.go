package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	version = "1.0"
	site    = "deeplooklabs.com"
)

var (
	banner = fmt.Sprintf(`
       _ _                       _    
  __ _(_) |_ ___ ___ __ _ _ _ __| |_  
 / _`+"`"+` | |  _(_-</ -_) _`+"`"+` | '_/ _| ' \ 
 \__, |_|\__/__/\___\__,_|_| \__|_||_|    %s
 |___/                                

        %s
`, version, site)

	example        = "gitsearch -q \"tesla.com boto language:python\""
	query          string
	token          string
	files          bool
	repositorys    bool
	dontShowBanner bool
)

func main() {

	flag.StringVar(&query, "q", "", "Query to search on Github")
	flag.StringVar(&token, "t", "", "Token of Github API or use on env ex: GITHUB_TOKEN")
	flag.BoolVar(&files, "f", true, "Show files urls")
	flag.BoolVar(&repositorys, "r", false, "Show repository urls")
	flag.BoolVar(&dontShowBanner, "silent", false, "Dont show banner")

	flag.Parse()

	if query == "" {
		fmt.Println(banner)
		fmt.Println(example)
		return
	}

	if dontShowBanner == true {

	} else {
		fmt.Println(banner)
	}

	searchTerm := url.QueryEscape(query)

	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	sortBy := "updated"
	headers := map[string]string{"Authorization": "Token " + token}
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s&sort=%s", searchTerm, sortBy)
	response, err := sendRequest(url, headers)
	if err != nil {
		fmt.Printf("[WRN] An error occurred while making the request! to %s", url)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("[WRN] An error occurred while making the request! {status_code: %d}", response.StatusCode)
		return
	}

	re := regexp.MustCompile(`page=(\d+)[^>]*>; rel="last"`)
	matches := re.FindStringSubmatch(response.Header.Get("Link"))

	if len(matches) == 2 {
		lastPageStr := matches[1]
		lastPage, err := strconv.Atoi(lastPageStr)
		fmt.Printf("[INF] Total pages found: %d\n", lastPage)

		if err != nil {
			fmt.Println("Error during conversion")
			return
		}

		for page := 1; page <= lastPage; page++ {
			url := fmt.Sprintf("https://api.github.com/search/code?q=%s&sort=%s&page=%d", searchTerm, sortBy, page)
			response, err := sendRequest(url, headers)

			if err != nil {
				fmt.Printf("Erro na solicitação HTTP para a página %d: %v\n", page, err)
				continue
			}

			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("[WRN] An error occurred while reading the response!")
				return
			}

			json := string(data[:])

			if files {
				items := gjson.Get(json, "items.#.html_url").Array()
				for _, value := range items {
					url := value.String()
					lastSlashIndex := strings.LastIndex(url, "/")
					nameOfFile := url[lastSlashIndex+1:]
					fmt.Printf("[%s] %s\n", nameOfFile, url)
				}
			}

			if repositorys {
				items := gjson.Get(json, "items.#.repository.html_url").Array()
				for _, value := range items {
					url := value.String()
					fmt.Printf("%s\n", url)
				}
			}

			time.Sleep(700 * time.Millisecond)

		}

	}
}

func sendRequest(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
