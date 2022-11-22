package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"
	"strings"

	_ "image/png"
	_ "image/jpeg"

	"github.com/PuerkitoBio/goquery"
	"github.com/ritchie46/GOPHY/img2gif"
)

func main() {
	usrPtr := flag.String("u", "carlos-mendible", "Your credly username")
	sizePtr := flag.Bool("s", false, "Generate a small gif")

	flag.Parse()

	user := *usrPtr
	size := "220x220"
	if *sizePtr {
		size = "110x110"
	}

	badges := []string{}

	fmt.Println("Reading badges for credly user " + user)
	fmt.Println("---")

	url := fmt.Sprintf("https://www.credly.com/users/%s/badges?sort=-state_updated_at", user)

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Get badges src
	doc.Find(".cr-standard-grid-item-content__image").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if size != "110x110" {
			src = strings.Replace(src, "110x110", size, 1)
		}
		title, _ := s.Parent().Parent().Attr("title")
		fmt.Println(title)
		badges = append(badges, src)
	})

	fmt.Println("---")
	fmt.Println("Creating animated gif")

	img := ReadImages(&badges)

	img_p := img2gif.EncodeImgPaletted(&img)

	fps := 1

	img2gif.WriteGif(&img_p, 100/fps, "credly.gif")

	fmt.Println("credly.gif successfully created")
}

// Read images from a slice with file locations.
func ReadImages(files *[]string) []image.Image {
	im := []image.Image{}

	for _, s := range *files {
		response, e := http.Get(s)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()

		img, _, err := image.Decode(response.Body)

		if err != nil {
			fmt.Println(err)
		}

		im = append(im, img)
	}

	return im
}

// Reference:
// https://github.com/ritchie46/GOPHY/blob/master/img2gif/img2gif.go
