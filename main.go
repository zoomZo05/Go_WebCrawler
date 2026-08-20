package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2  {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2  {
		fmt.Println( "too many arguments provided" )
		os.Exit(1)
	}

	firstArg := os.Args[1]

	// result,err := getHTML(firstArg)
	// if err != nil{
	// 	fmt.Printf("%s", err)
	// }



	// fmt.Printf("starting crawl of: %s", result)

	pages := map[string]int{}

	crawlPage(firstArg,firstArg, pages)

	for page, count := range pages {
		fmt.Printf("Found %d internal links to %s\n", count, page)
	}
}
