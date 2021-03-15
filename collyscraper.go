package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Job struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Url      string `json:"url"`
}

var jobs = make([]Job, 0, 200)

// For testing only
func printResults(job Job) {
	fmt.Println(job.Title, " -> ", job.Company, " -> ", job.Location, " -> ", job.Url)
}

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

		jobs = append(jobs, job)
	})
}

func scrapeIndeed(c *colly.Collector) {
	c.OnHTML("div.jobsearch-SerpJobCard", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.jobtitle"))
		company := strings.TrimSpace(e.ChildText("span.company"))
		location := strings.TrimSpace(e.ChildText("span.location"))
		url := strings.TrimSpace(e.ChildAttr("a[href]", "href"))
		url = "https://www.indeed.com/viewjob" + url

		job := Job{
			Title:    title,
			Company:  company,
			Location: location,
			Url:      url,
		}

		jobs = append(jobs, job)
	})
}

func scrapeStack(c *colly.Collector) {
	c.OnHTML("div.-job", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.s-link"))
		company := strings.TrimSpace(e.ChildText("h3 > span"))
		location := strings.TrimSpace(e.ChildText("span.fc-black-500"))
		url := strings.TrimSpace(e.ChildAttr("a.s-link[href]", "href"))
		url = "https://stackoverflow.com" + url

		job := Job{
			Title:    title,
			Company:  company,
			Location: location,
			Url:      url,
		}

		jobs = append(jobs, job)
	})
}

func errorHandler(c *colly.Collector) {
	var collyError error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		collyError = err
	})

	if collyError != nil {
		// Return a response including error
	}
}

func main() {
	// Instantiate default collectors
	cMonster := colly.NewCollector()
	cIndeed := colly.NewCollector()
	cStack := colly.NewCollector()

	// var numVisited = 0 // Specify max num of requests
	// c.OnRequest(func(r *colly.Request) {
	// 	if numVisited > 100 {
	// 		r.Abort()
	// 	}
	// 	numVisited++
	// })

	errorHandler(cMonster)
	errorHandler(cIndeed)
	errorHandler(cStack)

	scrapeMonster(cMonster)
	scrapeIndeed(cIndeed)
	scrapeStack(cStack)

	cMonster.Visit("https://www.monster.com/jobs/search/?q=software-engineer&where=Raleigh__2C-NC&intcid=skr_navigation_nhpso_searchMain")
	cIndeed.Visit("https://www.indeed.com/jobs?q=software+engineer&l=Raleigh,+NC&explvl=entry_level")
	cStack.Visit("https://stackoverflow.com/jobs?q=Software+Engineer&l=North+Carolina%2C+USA&d=20&u=Miles")

	// Print results for testing only
	for i := 0; i < len(jobs); i++ {
		printResults(jobs[i])
	}
}
