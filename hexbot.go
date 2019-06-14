package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	h   int
	w   int
	cnt int
)

type coordinate struct {
	X int
	Y int
}

type pixelData struct {
	Value       string
	Coordinates coordinate
}

type pixelArray struct {
	Colors []pixelData
}

type imgLinks struct {
	Raw, Full, Regular, Small, Thumb string
}

type imgUrls struct {
	Urls imgLinks
}

func main() {
	var images []*image.RGBA
	var colours pixelArray

	h = 100
	w = 100
	cnt = 1000

	fmt.Println("Go Gopher boi! x20")

	getImage()

	for i := 0; i < 20; i++ {
		img := image.NewRGBA(image.Rect(0, 0, h, w))

		url := fmt.Sprintf("https://api.noopschallenge.com/hexbot?count=%d&width=%d&height=%d%s",
			cnt, h, w, "")

		images = append(images, img)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("HTTP request has gophailed smh")
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal([]byte(data), &colours)

			for _, pixel := range colours.Colors {
				pixVal, _ := parseHexColor(pixel.Value)
				images[i].Set(pixel.Coordinates.X, pixel.Coordinates.Y, pixVal)
			}

		}
	}

	saveGIF("giffun.gif", images)

	fmt.Println("Crocodile Done-deey")
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

func saveImg(fileName string, img image.Image) {
	f, err := os.Create("output/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func saveGIF(fileName string, images []*image.RGBA) {

	outGif := &gif.GIF{}
	f, _ := os.OpenFile("output/"+fileName, os.O_WRONLY|os.O_CREATE, 0600)
	for _, simage := range images {
		palettedImage := image.NewPaletted(image.Rect(0, 0, h, w), palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, simage, image.Rect(0, 0, h, w).Min, draw.Over)
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 0)
	}
	defer f.Close()
	gif.EncodeAll(f, outGif)
}

func getImage() {
	var list imgUrls
	endpoint := "https://api.unsplash.com/photos/random?client_id=4b30f506ef4e2e506abe9edd3156eb33dc99194ddeb1de27bbd73aac14c7da84"
	response, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("HTTP request has gophailed smh")
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal([]byte(data), &list)
		fmt.Printf("Look smh --> %+v", list.Urls.Thumb)
	}

	response, e := http.Get(list.Urls.Thumb)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := os.Create("output/stolen.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Stole the image.")
}
