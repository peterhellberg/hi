/*

Package hi allows you to find images for a given hashtag

Usage

A small usage example

    package main

    import (
    	"fmt"

    	"github.com/peterhellberg/hi"
    )

    func main() {
    	images, err := hi.FindShuffledImages("pixel_dailies")

    	if err == nil && len(images) > 0 {
    		fmt.Println(images[0].URL)
    	}
    }

*/
package hi

import (
	"errors"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	// ErrNoImagesFound is the error returned when no images were found
	ErrNoImagesFound = errors.New("no images found")
)

// Image contains the fields for an image
type Image struct {
	URL     string
	ItemID  string
	TweetID string
	Height  string
	Width   string
	Ratio   string
}

// Scraper holds references to the URL and parsed goquery document
type Scraper struct {
	URL      string
	Document *goquery.Document
}

// NewScraper creates a new scraper
func NewScraper(hashtag string) *Scraper {
	return &Scraper{
		URL: "https://twitter.com/hashtag/" + hashtag + "?f=images",
	}
}

// FindImages finds images
func (s *Scraper) FindImages() ([]Image, error) {
	images := []Image{}

	if s.Document == nil {
		doc, err := goquery.NewDocument(s.URL)
		if err != nil {
			return nil, err
		}

		s.Document = doc
	}

	s.Document.Find("span.AdaptiveStreamGridImage").Each(func(i int, s *goquery.Selection) {
		if dataURL, ok := s.Attr("data-url"); ok {
			images = append(images, Image{
				URL:     dataURL,
				ItemID:  s.AttrOr("data-item-id", ""),
				TweetID: s.AttrOr("data-tweet-id", ""),
				Height:  s.AttrOr("data-height", ""),
				Width:   s.AttrOr("data-width", ""),
			})
		}
	})

	return images, nil
}

// FindImage finds image
func (s *Scraper) FindImage() (Image, error) {
	return singleImage(s.FindImages())
}

// FindShuffledImages finds images and shuffles them
func (s *Scraper) FindShuffledImages() ([]Image, error) {
	return shuffledImages(s.FindImages())
}

// FindShuffledImage finds shuffled image
func (s *Scraper) FindShuffledImage() (Image, error) {
	return singleImage(s.FindShuffledImages())
}

// FindShuffledImages first creates a scraper, then finds images and shuffles them
func FindShuffledImages(hashtag string) ([]Image, error) {
	return NewScraper(hashtag).FindShuffledImages()
}

// FindShuffledImage first creates a scraper, then finds shuffled image
func FindShuffledImage(hashtag string) (Image, error) {
	return NewScraper(hashtag).FindShuffledImage()
}

// FindImages first creates a scraper, then finds images
func FindImages(hashtag string) ([]Image, error) {
	return NewScraper(hashtag).FindImages()
}

// FindImage first creates a scraper, then finds image
func FindImage(hashtag string) (Image, error) {
	return NewScraper(hashtag).FindImage()
}

func singleImage(images []Image, err error) (Image, error) {
	if err != nil {
		return Image{}, err
	}

	if len(images) == 0 {
		return Image{}, ErrNoImagesFound
	}

	return images[0], nil
}

func shuffledImages(images []Image, err error) ([]Image, error) {
	if err != nil {
		return images, err
	}

	rand.Seed(time.Now().UnixNano())

	for i := range images {
		j := rand.Intn(i + 1)
		images[i], images[j] = images[j], images[i]
	}

	return images, nil
}
