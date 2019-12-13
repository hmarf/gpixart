package main

import (
	"fmt"
	image "image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func rgbaToUnit8s(img *image.RGBA) []color.RGBA {
	n_pixels := (img.Rect.Max.X - img.Rect.Min.X) * (img.Rect.Max.Y - img.Rect.Min.Y)
	vcolor := make([]color.RGBA, n_pixels)
	index := 0
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			vcolor[index], _ = img.At(x, y).(color.RGBA)
			index += 1
		}
	}
	return vcolor
}

func distance(color1 color.RGBA, color2 color.RGBA) float64 {
	r1, g1, b1, a1 := color1.R, color1.G, color1.B, color1.A
	r2, g2, b2, a2 := color2.R, color2.G, color2.B, color2.A
	r1f := float64(r1)
	g1f := float64(g1)
	b1f := float64(b1)
	a1f := float64(a1)

	r2f := float64(r2)
	g2f := float64(g2)
	b2f := float64(b2)
	a2f := float64(a2)

	distance := math.Sqrt(math.Pow(r2f-r1f, 2) + math.Pow(g2f-g1f, 2) + math.Pow(b2f-b1f, 2) + math.Pow(a2f-a1f, 2))
	return distance
}

func kmeans(img *image.RGBA, oImage *image.RGBA, n_cluster int) {

	vcolor := rgbaToUnit8s(img)
	n_pixels := len(vcolor)
	vcluster := make([]color.RGBA, n_cluster)
	residual := float32(n_pixels)

	rand.Seed(time.Now().UnixNano())
	vtype := make([]int, n_pixels)
	for i := 0; i < len(vtype); i++ {
		vtype[i] = rand.Intn(n_cluster)
	}

	n_iter := 0
	for residual > 0 && n_iter < 30 {
		residual = 0
		for i := 0; i < n_cluster; i++ {
			vtype_cluster_i := make([]int, 0)
			for index, type_cluster := range vtype {
				if type_cluster == i {
					vtype_cluster_i = append(vtype_cluster_i, index)
				}
			}
			if len(vtype_cluster_i) == 0 {
				continue
			}
			n_vtype_cluster_i := float64(len(vtype_cluster_i))
			r_sum, g_sum, b_sum, a_sum := 0.0, 0.0, 0.0, 0.0
			for _, type_cluster_i := range vtype_cluster_i {
				color_ := vcolor[type_cluster_i]
				r_sum += float64(color_.R) / n_vtype_cluster_i
				g_sum += float64(color_.G) / n_vtype_cluster_i
				b_sum += float64(color_.B) / n_vtype_cluster_i
				a_sum += float64(color_.A) / n_vtype_cluster_i
			}
			vcluster[i] = color.RGBA{uint8(r_sum), uint8(g_sum), uint8(b_sum), uint8(a_sum)}
		}

		for vtype_index, color_ := range vcolor {
			cluster_index_min := vtype[vtype_index]
			distance_min := 1000.0
			for cluster_index, cluster := range vcluster {
				distance := distance(color_, cluster)
				if distance < distance_min {
					distance_min = distance
					cluster_index_min = cluster_index
				}
			}
			if cluster_index_min != vtype[vtype_index] {
				residual += 1
			}
			vtype[vtype_index] = cluster_index_min
		}

		n_iter += 1
	}

	for index := 0; index < n_pixels; index++ {
		vcolor[index] = vcluster[vtype[index]]
	}
	updataImageByUint8s(img, vcolor)
}

func updataImageByUint8s(img *image.RGBA, vcolor []color.RGBA) {
	index := 0
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			img.Set(x, y, vcolor[index])
			index += 1
		}
	}
}

func makeOutputImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func resizeImage(img image.Image, width uint, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func resizeAndMakeImage(img image.Image, width uint, height uint, n_cluster int) *image.RGBA {
	oImage := makeOutputImage(int(width), int(height))
	if _, ok := img.(*image.NRGBA); ok {
		if aa, ok := resizeImage(img, width, height).(*image.RGBA); ok {
			kmeans(aa, oImage, n_cluster)
			return aa
		}
	} else if _, ok := img.(*image.RGBA); ok {
		if aa, ok := resizeImage(img, width, height).(*image.RGBA); ok {
			kmeans(aa, oImage, n_cluster)
			return aa
		}
	} else if _, ok := img.(*image.YCbCr); ok {
		b := img.Bounds()
		m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
		if aa, ok := resizeImage(m, width, height).(*image.RGBA); ok {
			kmeans(aa, oImage, n_cluster)
			return aa
		}
	}
	return nil
}

func calcurateImageSize(h, w int) (newH, newW int) {
	minSize := 50
	if w > h {
		newW = minSize
		newH = h / w * minSize
	} else if w < h {
		newH = minSize
		newW = w / h * minSize
	} else {
		newW = minSize
		newH = minSize
	}
	return
}

func main() {

	file, err := os.Open("./image/pokemon.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Print(err)
	}

	rct := img.Bounds()
	height, width := calcurateImageSize(rct.Dy(), rct.Dx())
	newImage := resizeAndMakeImage(img, uint(height), uint(width), 2)
	file, err = os.Create("./output.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, newImage, &jpeg.Options{100}); err != nil {
		panic(err)
	}
}
