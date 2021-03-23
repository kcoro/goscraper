package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

type RequestData struct {
	Title    string `json:"title"`
	Location string `json:"location"`
}

type Job struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Url      string `json:"url"`
}

var reqData = RequestData{Title: "", Location: ""}
var jobs = make([]Job, 0, 300)

func scrapeIndeed(c *colly.Collector) {
	c.OnHTML("div.jobsearch-SerpJobCard", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.jobtitle"))
		company := strings.TrimSpace(e.ChildText("span.company"))
		location := strings.TrimSpace(e.ChildText("span.location"))
		url := strings.TrimSpace(e.ChildAttr("h2.title > a", "href"))

		if !strings.Contains(url, "https://indeed.com") {
			url = "https://indeed.com" + url
		}

		job := Job{
			Title:    title,
			Company:  company,
			Location: location,
			Url:      url,
		}

		if location == "" {
			location = reqData.Location
		}
		if title != "" && url != "" {
			jobs = append(jobs, job)
		}
	})
}

func scrapeStack(c *colly.Collector) {
	c.OnHTML("div.-job", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.s-link"))
		company := strings.TrimSpace(e.ChildText("h3 > span"))
		// Remove newline char and all chars after newline in company name
		companyFinal := company[0:strings.Index(company, "\n")]
		location := strings.TrimSpace(e.ChildText("span.fc-black-500"))
		url := strings.TrimSpace(e.ChildAttr("a.s-link[href]", "href"))
		url = "https://stackoverflow.com" + url

		job := Job{
			Title:    title,
			Company:  companyFinal,
			Location: location,
			Url:      url,
		}

		if location == "" {
			location = reqData.Location
		}
		if title != "" && url != "" {
			jobs = append(jobs, job)
		}
	})
}

// Monster recently changed website structure
// need to redesign scraper
func scrapeMonster(c *colly.Collector) {
	c.OnHTML("div.flex-row", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a"))
		company := strings.TrimSpace(e.ChildText("div.company > span.name"))
		location := strings.TrimSpace(e.ChildText("div.location > span.name"))
		url := strings.TrimSpace(e.ChildAttr("a[href]", "href"))

		job := Job{
			Title:    title,
			Company:  company,
			Location: location,
			Url:      url,
		}

		if location == "" {
			location = reqData.Location
		}
		if title != "" && url != "" {
			jobs = append(jobs, job)
		}
	})
}

// Example request query: ?title=software-developer&location=Miami-FL
func Handler(w http.ResponseWriter, r *http.Request) {

	// Get request query params from url
	query := r.URL.Query()
	reqData.Title = query.Get("title")
	reqData.Location = query.Get("location")

	// Mormalize search query terms for scraping
	strings.Replace(reqData.Title, " ", "+", -1)
	strings.Replace(reqData.Location, " ", "+", -1)
	strings.Replace(reqData.Title, "-", "+", -1)
	strings.Replace(reqData.Location, "-", "+", -1)

	// Instantiate default collectors
	cStack := colly.NewCollector()
	cIndeed := colly.NewCollector()
	// Monster recently changed website structure
	// cMonster := colly.NewCollector()

	scrapeStack(cStack)
	scrapeIndeed(cIndeed)
	// Monster recently changed website structure
	// scrapeMonster(cMonster)

	cStack.Visit("https://stackoverflow.com/jobs?q=" + reqData.Title + "&l=" + reqData.Location + "USA&d=20&u=Miles")
	cIndeed.Visit("https://www.indeed.com/jobs?q=" + reqData.Title + "&l=" + reqData.Location + "&explvl=entry_level")
	// Monster recently changed website structure
	// need to redesign scraper
	// cMonster.Visit("https://www.monster.com/jobs/search/?q=" + reqData.Title + "&where=" + reqData.Location + "&intcid=skr_navigation_nhpso_searchMain")

	// Encode map[string]Job to json
	jobsJson, _ := json.Marshal(jobs)
	jobsJson = bytes.Replace(jobsJson, []byte("\\u0026"), []byte("&"), -1) // replace explicit unicode code with &

	// Set responsewriter's header to let client expect json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*") // get around cors
	// Send json with response writer
	fmt.Fprintf(w, string(jobsJson))

	jobs = make([]Job, 0, 300) // Clear array of Jobs for next request

	// Close request
	r.Body.Close()
}

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
