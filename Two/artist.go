package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

type imgLinks struct {
	Raw, Full, Regular, Small, Thumb string
}

type imgUrls struct {
	Urls imgLinks
}

func main() {
	// getImage()
	img := readStolen("output/stolen.jpeg")

	h := img.Bounds().Max.Y
	w := img.Bounds().Max.X
	palette := getPalette(img, 20)
	// dc := gg.NewContextForImage(img)
	dc := gg.NewContext(w, h)

	for i := 0; i < w; i += 10 {
		for j := 0; j < h; j += 10 {
			r, g, b, _ := img.At(i, j).RGBA()
			closest := color.RGBA{R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 0xff}

			dc = paintDot(dc, float64(i), float64(j), 5, getClosestColor(palette, closest))
		}
	}
	dc.SaveJPG("output/out.jpeg", 80)
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
	}

	response, e := http.Get(list.Urls.Regular)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := os.Create("output/stolen.jpeg")
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

func readStolen(url string) image.Image {
	infile, err := os.Open(url)
	if err != nil {
		fmt.Println("Couldn't open stolen goods smh")
		panic(err)
	}
	defer infile.Close()

	thumbnail, _, err := image.Decode(infile)
	if err != nil {
		fmt.Println("Big problem with decoding image.")
		panic(err)
	}

	b := thumbnail.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), thumbnail, b.Min, draw.Src)

	return m
}

func getPalette(img image.Image, k int) []color.RGBA {
	var palette []color.RGBA

	width := img.Bounds().Max.Y
	out, _ := prominentcolor.KmeansWithAll(k, img, prominentcolor.ArgumentDefault, uint(width)/10, prominentcolor.GetDefaultMasks())

	for _, rgb := range out {

		palette = append(palette, color.RGBA{R: uint8(rgb.Color.R),
			G: uint8(rgb.Color.G),
			B: uint8(rgb.Color.B),
			A: 0xff})
	}

	paletteImg := image.NewRGBA(image.Rect(0, 0, 100*k, 100))
	for i := 0; i < k; i++ {
		for j := 0; j < 100; j++ {
			for l := 0; l < 100; l++ {
				paletteImg.Set(j+100*i, l, palette[i])
			}
		}
	}
	f, err := os.Create("output/palette.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, paletteImg)

	return palette
}

func paintDot(dc *gg.Context, x float64, y float64, r float64, shade color.RGBA) *gg.Context {
	dc.SetRGBA255(int(shade.R), int(shade.G), int(shade.B), int(shade.A))
	dc.DrawPoint(x, y, r)
	dc.Fill()

	return dc
}

func getClosestColor(palette []color.RGBA, shade color.RGBA) color.RGBA {
	var closest color.RGBA
	var minDst float64 = 100

	c1 := colorful.Color{float64(shade.R) / 255.0, float64(shade.G) / 255.0, float64(shade.B) / 255.0}

	for i, clr := range palette {
		dst := c1.DistanceLab(colorful.Color{float64(clr.R) / 255.0,
			float64(clr.G) / 255.0, float64(clr.B) / 255.0})
		// fmt.Println(minDst, dst)
		if dst < minDst {
			closest = palette[i]
			minDst = dst
		}
	}
	return closest
}
