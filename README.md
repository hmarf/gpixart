# gpixart( pixel art )
## Overview
Convert the image into a pixel art.
Size and number of colors can be specified.

<img src="https://github.com/hmarf/gpixart/blob/master/img/summaryImage.png?raw=true" width="700px">

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
The resized-image is also output at the same time. The image is a full-color pixel art.

<img src="https://github.com/hmarf/gpixart/blob/master/img/resize.png?raw=true" width="150px">

- Convert './image/pokemon.jpg' to './output.jpg' 40 pixel, 16 color
```
gpixart -i ./image/pokemon.jpg -o ./outout.jpg -m 40 -c 16
```
<img src="https://github.com/hmarf/gpixart/blob/master/img/16color.png?raw=true" width="150px">

- Convert './image/pokemon.jpg' to './output.jpg' 40 pixel, 4 color
```
gpixart -i ./image/pokemon.jpg -o ./outout.jpg -m 40 -c 4
```
<img src="https://github.com/hmarf/gpixart/blob/master/img/4color.png?raw=true" width="150px">

## Reference
- http://dot-e-nanika.com/
