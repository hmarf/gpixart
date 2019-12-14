# gpixart( pixel art )
## Overview
Convert the image into a pixel art.

<img src="https://github.com/hmarf/pixelArt-golang/blob/master/img/upload.png?raw=true" width="700px">

## Usage

```
NAME:
   gpixart - Convert the image into a pixel art.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   hmarf

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value, -i value   [required] string
                             input file path (default: "None")
   --output value, -o value  string
                             output file path [default] ./output.jpg (default: "./output.jpg")
   --minP value, -m value    int
                             Minimum number of pixels (default: 50)
   --color value, -c value   int
                             Number of colors used for pixel art (default: 16)
   --help, -h                show help
   --version, -v             print the version
```

## Example

Convert './image/pokemon.jpg' to './output.jpg' 32 pixel

```
gpixart -i ./image/pokemon.jpg -o ./outout.jpg -m 30 -c 32
```

