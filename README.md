This is a Go project which scrapes data from job boards Monster, Indeed, and Stack Overflow.
The program runs a web server listening for requests. When a request is received it takes query params and scrapes job sites
with the corresponding job title and location. The program returns a json response listing results.

## Getting Started

Make sure you have Go installed on your machine.
```bash
go version
```

Clone, Build and Run Instructions
- Note on Windows, use backslash '\\' instead of forwardslash '/' for all commands.

```bash
# keep your go projects in your go path's go/src
# cd into your paths go/src
cd ~/go/src/
# clone this repo with https
git clone https://github.com/kcoro/goscraper.git
# cd into new project directory
cd ~/go/src/goscraper
# build the .exe
go build
# run the .exe
./jobscraper
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the default result.

## How to make Queries
This project accepts query parameters in the URL.
Query labels include: title, location.

 - Structure of a query: http://localhost:3000/?title=foo&location=bar
 - Example: [http://localhost:3000/?title=software-engineer&location=Miami-FL](http://localhost:3000/?title=software-engineer&location=Miami-FL)
