package main

import (
	"flag"
	"fmt"

	"github.com/peterhellberg/hi"
)

var (
	limit   = flag.Int("l", 3, "limit number of images")
	shuffle = flag.Bool("s", true, "shuffle the images")
)

func main() {
	flag.Parse()

	hashtag := flag.Arg(0)

	if hashtag == "" {
		fmt.Println("Missing hashtag")
		return
	}

	findImages := hi.FindImages

	if *shuffle {
		findImages = hi.FindShuffledImages
	}

	if images, err := findImages(hashtag); err == nil {
		for i, img := range images {
			fmt.Println(img.URL)

			if i == *limit-1 {
				break
			}
		}
	}
}
