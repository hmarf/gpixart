package main

import (
	"os"

	pixelart "github.com/hmarf/pixelArt-golang/pixelArt"
	"github.com/urfave/cli"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "pixel art"
	// app.Usage = "Trunks is a simple command line tool for HTTP load testing."
	app.Version = "0.1.2"
	app.Author = "hmarf"
	return app
}

func Action(c *cli.Context) {
	// app := App()
	option := pixelart.Option{
		InputFile:  "image/pokemon.png",
		OutputFile: "aaa.jpg",
		MinSize:    50,
		Ncolor:     4,
	}
	pixelart.PixelArt(option)
}

func main() {
	app := App()
	app.Action = Action
	app.Run(os.Args)
}
