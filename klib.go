package main

var monsterURL = "https://www.monster.com/jobs/search/?q="
var indeedURL = "https://www.indeed.com/jobs?q="
var stackURL = "https://stackoverflow.com/jobs?q="
var allUrls = []string{monsterURL, indeedURL, stackURL}

var foundTitle string
var foundCompany string
var foundLocation string
var foundURL string

func buildMonsterURL(title string, location string) string {
	return "https://www.monster.com/jobs/search/?q=" + title + "&where=" + location
}

func buildIndeedURL(title string, location string) string {
	return "https://www.indeed.com/jobs?q=" + title + "&l=" + location + "&explvl=entry_level"
}

func buildStackURL(title string, location string) string {
	return "https://stackoverflow.com/jobs?q=" + title + "&l=" + location + "%2C+USA&d=20&u=Miles"
}
