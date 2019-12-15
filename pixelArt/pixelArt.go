package pixelart

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/nfnt/resize"
)

type Option struct {
	InputFile  string
	OutputFile string
	MinSize    int
	Ncolor     int
}

func rgbaToArray(img *image.RGBA) []color.RGBA {
	pixels := (img.Rect.Max.X - img.Rect.Min.X) * (img.Rect.Max.Y - img.Rect.Min.Y)
	vcolor := make([]color.RGBA, pixels)
	index := 0
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			vcolor[index], _ = img.At(x, y).(color.RGBA)
			index++
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

func kmeans(img *image.RGBA, oImage *image.RGBA, cluster int, size int) *image.RGBA {

	vcolor := rgbaToArray(img)
	npixels := len(vcolor)
	vcluster := make([]color.RGBA, cluster)
	residual := float32(npixels)

	rand.Seed(time.Now().UnixNano())
	vtype := make([]int, npixels)
	for i := 0; i < len(vtype); i++ {
		vtype[i] = rand.Intn(cluster)
	}

	niter := 0
	for residual > 0 && niter < 30 {
		residual = 0
		for i := 0; i < cluster; i++ {
			clusterInt := make([]int, 0)
			for index, typeCluster := range vtype {
				if typeCluster == i {
					clusterInt = append(clusterInt, index)
				}
			}
			if len(clusterInt) == 0 {
				continue
			}
			nclusterInt := float64(len(clusterInt))
			rS, gS, bS, aS := 0.0, 0.0, 0.0, 0.0
			for _, typeCluster := range clusterInt {
				color_ := vcolor[typeCluster]
				rS += float64(color_.R) / nclusterInt
				gS += float64(color_.G) / nclusterInt
				bS += float64(color_.B) / nclusterInt
				aS += float64(color_.A) / nclusterInt
			}
			vcluster[i] = color.RGBA{uint8(rS), uint8(gS), uint8(bS), uint8(aS)}
		}

		for vTypeIndex, color_ := range vcolor {
			clusterIndexMin := vtype[vTypeIndex]
			distanceMin := 1000.0
			for clusterIndex, cluster := range vcluster {
				distance := distance(color_, cluster)
				if distance < distanceMin {
					distanceMin = distance
					clusterIndexMin = clusterIndex
				}
			}
			if clusterIndexMin != vtype[vTypeIndex] {
				residual++
			}
			vtype[vTypeIndex] = clusterIndexMin
		}

		niter++
	}

	for index := 0; index < npixels; index++ {
		vcolor[index] = vcluster[vtype[index]]
	}
	return updataImage(img, vcolor, size)
}

func updataImage(img *image.RGBA, vcolor []color.RGBA, size int) *image.RGBA {
	imgSrc := img.Bounds()
	newImage := image.NewRGBA(image.Rect(0, 0, imgSrc.Dx()*size, imgSrc.Dy()*size))
	index := 0
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					newImage.Set(x*size+i, size*y+j, vcolor[index])
				}
			}
			index++
		}
	}
	return newImage
}

func makeOutputImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func resizeImage(img image.Image, nwidth uint, nheight uint) image.Image {
	resizeImage := resize.Resize(nwidth, nheight, img, resize.Lanczos3)
	file, err := os.Create("./resize.jpg")
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "creation of the save destination file failed.")
	}
	defer file.Close()
	if err := jpeg.Encode(file, resizeImage, &jpeg.Options{100}); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "Failed to save image.")
	}
	return resizeImage
}

func resizeAndMakeImage(img image.Image, width uint, height uint, cluster int, minSize int) *image.RGBA {
	oImage := makeOutputImage(int(width), int(height))
	if _, ok := img.(*image.NRGBA); ok {
		if aa, ok := resizeImage(img, width, height).(*image.RGBA); ok {
			newImage := kmeans(aa, oImage, cluster, minSize)
			return newImage
		}
	} else if _, ok := img.(*image.RGBA); ok {
		if aa, ok := resizeImage(img, width, height).(*image.RGBA); ok {
			newImage := kmeans(aa, oImage, cluster, minSize)
			return newImage
		}
	} else if _, ok := img.(*image.YCbCr); ok {
		b := img.Bounds()
		m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
		if aa, ok := resizeImage(m, width, height).(*image.RGBA); ok {
			newImage := kmeans(aa, oImage, cluster, minSize)
			return newImage
		}
	}
	return nil
}

func calcurateImageSize(h, w, minSize int) (newH, newW int) {
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

func PixelArt(o Option) {
	file, err := os.Open(o.InputFile)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "no such file or directory")
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "The file cannot be read. Please select an image file.")
		return
	}

	rct := img.Bounds()
	height := rct.Dy()
	width := rct.Dx()
	nheight, nwidth := calcurateImageSize(height, width, o.MinSize)
	newImage := resizeAndMakeImage(img, uint(nheight), uint(nwidth), o.Ncolor, o.MinSize)
	file, err = os.Create(o.OutputFile)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "creation of the save destination file failed.")
		return
	}
	defer file.Close()

	if err := jpeg.Encode(file, newImage, &jpeg.Options{100}); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "Failed to save image.")
		return
	}
}
