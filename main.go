package main

import (
	"os"

	pixelart "github.com/hmarf/pixelArt-golang/pixelArt"
	"github.com/urfave/cli"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "pixel art"
	app.Usage = "Convert the image into a pixel art."
	app.Version = "0.0.1"
	app.Author = "hmarf"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Value: "None",
			Usage: "[required] string\n	input file path",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "./output.jpg",
			Usage: "string\n	output file path [default] ./output.jpg",
		},
		cli.IntFlag{
			Name:  "minP, m",
			Value: 50,
			Usage: "int\n	Minimum number of pixels",
		},
		cli.IntFlag{
			Name:  "color, c",
			Value: 16,
			Usage: "int\n	Number of colors used for pixel art",
		},
	}
	return app
}

func Action(c *cli.Context) {
	app := App()
	if c.String("input") == "None" {
		app.Run(os.Args)
		return
	}
	option := pixelart.Option{
		InputFile:  c.String("input"),
		OutputFile: c.String("output"),
		MinSize:    c.Int("minP"),
		Ncolor:     c.Int("color"),
	}
	pixelart.PixelArt(option)
}

func main() {
	app := App()
	app.Action = Action
	app.Run(os.Args)
}
