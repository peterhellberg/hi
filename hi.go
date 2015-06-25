/*

Package hi allows you to find images for a given hashtag

*/
package hi

import (
	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"
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

// FindShuffledImages finds images and shuffles them
func (s *Scraper) FindShuffledImages() ([]Image, error) {
	rand.Seed(time.Now().UnixNano())

	images, err := s.FindImages()
	if err != nil {
		return images, err
	}

	for i := range images {
		j := rand.Intn(i + 1)
		images[i], images[j] = images[j], images[i]
	}

	return images, nil
}

// FindShuffledImages first creates a scraper, then finds images and shuffles them
func FindShuffledImages(hashtag string) ([]Image, error) {
	return NewScraper(hashtag).FindShuffledImages()
}

// FindImages first creates a scraper, then finds images
func FindImages(hashtag string) ([]Image, error) {
	return NewScraper(hashtag).FindImages()
}
