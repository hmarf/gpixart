package main

import (
	image "image"
	"image/color"
	"image/draw"
	jpeg "image/jpeg"
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
	// ベクトル
	n_pixels := len(vcolor)
	// クラスタ中心数 = 8 * ビット数
	// クラスタ中心のベクトル
	vcluster := make([]color.RGBA, n_cluster)
	// 残差
	residual := float32(n_pixels)

	// 初期化
	rand.Seed(time.Now().UnixNano())
	vtype := make([]int, n_pixels)
	for i := 0; i < len(vtype); i++ {
		vtype[i] = rand.Intn(n_cluster)
	}

	// k-means
	n_iter := 0
	for residual > 0 && n_iter < 30 {
		residual = 0
		// vclusterの更新
		// vtype から filterして特定のcluter中心に対応するcolorを取り出して平均を計算する
		for i := 0; i < n_cluster; i++ {
			// vtypeのうち，cluster i に属するindexのみを取り出す
			vtype_cluster_i := make([]int, 0)
			for index, type_cluster := range vtype {
				if type_cluster == i {
					vtype_cluster_i = append(vtype_cluster_i, index)
				}
			}
			// vtype_cluster_iが0個ならスルー
			if len(vtype_cluster_i) == 0 {
				continue
			}
			n_vtype_cluster_i := float64(len(vtype_cluster_i))
			// type_cluster_i
			r_sum, g_sum, b_sum, a_sum := 0.0, 0.0, 0.0, 0.0
			for _, type_cluster_i := range vtype_cluster_i {
				color_ := vcolor[type_cluster_i]
				r_sum += float64(color_.R) / n_vtype_cluster_i
				g_sum += float64(color_.G) / n_vtype_cluster_i
				b_sum += float64(color_.B) / n_vtype_cluster_i
				a_sum += float64(color_.A) / n_vtype_cluster_i
			}
			// クラスタ中心の色更新
			vcluster[i] = color.RGBA{uint8(r_sum), uint8(g_sum), uint8(b_sum), uint8(a_sum)}
		}

		// vtypeの更新
		for vtype_index, color_ := range vcolor {
			// どのclusterに距離が近いか
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

	// 色をcluster中心の色に書き換える
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

func resizeAndMakeImage(img image.Image, width uint, height uint, n_cluster int) *image.RGBA {
	img_resized := resize.Resize(width, height, img, resize.Lanczos3)
	oImage := makeOutputImage(int(width), int(height))
	if x, ok := img_resized.(*image.RGBA); ok {
		kmeans(x, oImage, n_cluster)
		return x
	} else if _, ok := img_resized.(*image.YCbCr); ok {
		b := img.Bounds()
		m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
		kmeans(m, oImage, n_cluster)
		return m
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
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	rct := img.Bounds()
	height, width := calcurateImageSize(rct.Dy(), rct.Dx())
	newImage := resizeAndMakeImage(img, uint(height), uint(width), 2)
	file, _ = os.Create("./image/output.jpg")
	defer file.Close()

	if err := jpeg.Encode(file, newImage, &jpeg.Options{100}); err != nil {
		panic(err)
	}
}
