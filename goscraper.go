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
var jobs = make([]Job, 0, 200)

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

func scrapeIndeed(c *colly.Collector) {
	c.OnHTML("div.jobsearch-SerpJobCard", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.jobtitle"))
		company := strings.TrimSpace(e.ChildText("span.company"))
		location := strings.TrimSpace(e.ChildText("span.location"))
		url := strings.TrimSpace(e.ChildAttr("a[href]", "href"))
		url = "https://indeed.com/pagead/clk?mo=r&" + url[17:]

		job := Job{
			Title:    title,
			Company:  company,
			Location: location,
			Url:      url,
		}

		if location == "" {
			location = reqData.Location
		}
		if title != "" || url != "" {
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

// Example request query: ?title=software-developer&location=Miami-FL
func Handler(w http.ResponseWriter, r *http.Request) {
	// Get request query params from url
	query := r.URL.Query()
	reqData.Title = query.Get("title")
	reqData.Location = query.Get("location")
	// Indeed uses diff formatting, + between words not - like Monster and SO
	indeedTitle := reqData.Title
	indeedLocation := reqData.Location

	// Remove space from request param and replace with %20 for querying job sites
	strings.Replace(reqData.Title, " ", "-", -1)
	strings.Replace(reqData.Location, " ", "-", -1)
	// Remove - and replace with + for indeed search query
	strings.Replace(indeedTitle, "-", "+", -1)
	strings.Replace(indeedLocation, "-", "+", -1)

	// Instantiate default collectors
	cMonster := colly.NewCollector()
	cIndeed := colly.NewCollector()
	cStack := colly.NewCollector()

	scrapeMonster(cMonster)
	scrapeIndeed(cIndeed)
	scrapeStack(cStack)

	cMonster.Visit("https://www.monster.com/jobs/search/?q=" + reqData.Title + "&where=" + reqData.Location + "&intcid=skr_navigation_nhpso_searchMain")
	cIndeed.Visit("https://www.indeed.com/jobs?q=" + indeedTitle + "&l=" + indeedLocation + "&explvl=entry_level")
	cStack.Visit("https://stackoverflow.com/jobs?q=" + reqData.Title + "&l=" + reqData.Location + "USA&d=20&u=Miles")

	// Encode map[string]Job to json
	jobsJson, _ := json.Marshal(jobs)
	jobsJson = bytes.Replace(jobsJson, []byte("\\u0026"), []byte("&"), -1) // replace explicit unicode code with &
	fmt.Fprint(w, string(jobsJson))                                        // must explicityly convert json to string before sending
	jobs = make([]Job, 0, 200)                                             // Clear array of Jobs for next request

	// Close request
	r.Body.Close()
}

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}