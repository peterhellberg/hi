package hi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewScraper(t *testing.T) {
	s := NewScraper("foobar")

	if s.URL != "https://twitter.com/hashtag/foobar?f=images" {
		t.Errorf(`unexpected url %q`, s.URL)
	}
}

func TestFindImages(t *testing.T) {
	for _, tt := range []struct {
		body     []byte
		count    int
		imageURL string
	}{
		{twoImagesBody, 2, "bar"},
		{sixImagesBody, 6, "corge"},
		{noImagesBody, 0, ""},
	} {
		server, scraper := testServerAndScraper(tt.body)
		defer server.Close()

		images, err := scraper.FindImages()

		if err != nil {
			t.Fatalf(`unexpected error: %v`, err)
		}

		if got := len(images); got != tt.count {
			t.Errorf(`unexpected number of images: %d, want %d`, got, tt.count)
		}

		if tt.count > 0 {
			if img := images[len(images)-1]; img.URL != tt.imageURL {
				t.Errorf(`unexpected image URL: %q, want %q`, img.URL, tt.imageURL)
			}
		}
	}
}

func TestFindImage(t *testing.T) {
	for _, tt := range []struct {
		body     []byte
		err      error
		imageURL string
	}{
		{twoImagesBody, nil, "foo"},
		{sixImagesBody, nil, "foo"},
		{noImagesBody, ErrNoImagesFound, ""},
	} {
		server, scraper := testServerAndScraper(tt.body)
		defer server.Close()

		img, err := scraper.FindImage()

		if err != tt.err {
			t.Fatalf(`unexpected error: %v`, err)
		}

		if img.URL != tt.imageURL {
			t.Errorf(`unexpected image URL: %q, want %q`, img.URL, tt.imageURL)
		}
	}
}

func TestFindShuffledImage(t *testing.T) {
	for _, tt := range []struct {
		body []byte
		err  error
	}{
		{twoImagesBody, nil},
		{sixImagesBody, nil},
		{noImagesBody, ErrNoImagesFound},
	} {
		server, scraper := testServerAndScraper(tt.body)
		defer server.Close()

		img, err := scraper.FindShuffledImage()

		if err != tt.err {
			t.Fatalf(`unexpected error: %v`, err)
		}

		if err == nil && img.URL == "" {
			t.Errorf(`no image url`)
		}
	}
}

func testServerAndScraper(body []byte) (*httptest.Server, *Scraper) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write(body)
		},
	))

	return server, &Scraper{URL: server.URL}
}

var noImagesBody = []byte(`
<html>
	<body>
	</body>
</html>`)

var twoImagesBody = []byte(`
<html>
	<body>
		<span class="AdaptiveStreamGridImage" data-url="foo"></span>
		<span class="AdaptiveStreamGridImage" data-url="bar"></span>
	</body>
</html>`)

var sixImagesBody = []byte(`
<html>
	<body>
		<span class="AdaptiveStreamGridImage" data-url="foo"></span>
		<span class="AdaptiveStreamGridImage" data-url="bar"></span>
		<span class="AdaptiveStreamGridImage" data-url="baz"></span>
		<span class="AdaptiveStreamGridImage" data-url="qux"></span>
		<span class="AdaptiveStreamGridImage" data-url="quux"></span>
		<span class="AdaptiveStreamGridImage" data-url="corge"></span>
	</body>
</html>`)
