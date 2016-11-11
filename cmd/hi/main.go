package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/peterhellberg/hi"
)

func main() {
	limit := flag.Int("l", 3, "limit number of images")
	shuffle := flag.Bool("s", false, "shuffle the images")
	jsonOutput := flag.Bool("json", false, "set output format to JSON")

	flag.Parse()

	hashtag := flag.Arg(0)

	if hashtag == "" {
		fmt.Println("Missing hashtag")
		os.Exit(1)
	}

	var (
		limitedImages = []hi.Image{}
		findImages    = hi.FindImages
	)

	if *shuffle {
		findImages = hi.FindShuffledImages
	}

	if images, err := findImages(hashtag); err == nil {
		l := *limit

		for i, img := range images {
			if i == l {
				break
			}

			limitedImages = append(limitedImages, img)
		}
	}

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(limitedImages)
	} else {
		for _, img := range limitedImages {
			fmt.Println(img.URL)
		}
	}
}
