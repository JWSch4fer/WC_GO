# WC_GO

Web Crawler written in Go

## Overview

WC_GO is a lightweight, concurrent web crawler built in Go. It recursively crawls a website starting from a specified base URL, collects internal links, and generates a formatted report of the pages discovered in decending order of visits.

## Features

- **Concurrent Crawling:** Leverages Go's goroutines and channels to crawl multiple pages simultaneously.
- **Configurable Concurrency:** Control the number of concurrent HTTP requests to balance speed and server load.
- **Timeout Handling:** Uses Go's context package to enforce a crawl timeout.
- **Maximum Page Limit:** Specify the maximum number of pages to crawl to avoid runaway processes.
- **Thread-Safe Operations:** Utilizes mutexes and wait groups to safely manage shared data.
- **Sorted Reporting:** Generates a report that sorts pages by the number of internal links (highest first) and alphabetically when counts are equal.

## Prerequisites

- [Go](https://golang.org) (version 1.24 or later recommended)

## Installation

Clone the repository and build the project:

```bash
git clone https://github.com/yourusername/WC_GO.git
cd WC_GO
go build -o wc_go main.go
```

##Usage

Run the crawler from the command line using the following syntax:

```
./wc_go <website> <maxConcurrency> <maxPages>
```

- <website> is the starting URL to crawl.
- <maxConcurrency> is the maximum number of concurrent HTTP requests.
- <maxPages> is the maximum number of pages to crawl.


## Testing
To run the tests, execute:
```
go test ./...
```

## License
Distributed under the MIT License. See LICENSE for more information.
