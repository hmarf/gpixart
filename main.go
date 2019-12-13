package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

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
	width := rct.Dx()
	height := rct.Dy()
	fmt.Println(width, height)
	fmt.Println(img.At(0, 0).RGBA())
	imgResize := resize.Resize(200, 300, img, resize.Lanczos3)

	// 出力用ファイル作成(エラー処理は略)
	file, _ = os.Create("./output.jpg")
	defer file.Close()

	// JPEGで出力(100%品質)
	if err := jpeg.Encode(file, imgResize, &jpeg.Options{100}); err != nil {
		panic(err)
	}
}
