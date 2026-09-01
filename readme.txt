========================================
Go Web Crawler
========================================

Prerequisites:
- Go 1.21+ installed
- Git installed

Setup:
1. Clone the repository:
   git clone <YOUR_REPO_URL>
   cd <REPO_FOLDER_NAME>

2. Install / tidy Go dependencies:
   go mod tidy

Usage:
   go run . <BASE_URL> <MAX_CONCURRENCY> <MAX_PAGES>

Arguments:
   <BASE_URL>         Starting URL to crawl
   <MAX_CONCURRENCY>  Maximum number of concurrent goroutines
   <MAX_PAGES>        Maximum number of pages to crawl

Example:
   go run . "https://learnwebscraping.dev/practice/ecommerce/" 2 20

Build (Optional):
   go build -o crawler .
   ./crawler "https://learnwebscraping.dev/practice/ecommerce/" 2 20