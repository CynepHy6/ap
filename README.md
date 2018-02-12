# apla packager

Simvolio and Protypo files packager to import on ecosystem

## Examples

### Unpack file from "basic.sim" to "./basic/"

>ap -u -i basic.sim

### Pack files from "basic/" to ./basic.json

>ap -i basic/

## Usage of "ap"

Without flags will start GUI. For view all flags please use "ap -h"

--input string

    -i, path for input files (default ".")

--output string

    -o, output filename for JSON (default "output" if "input" not found)

--unpack

    -u, unpacking mode


## build

### linux

    go build

### windows

    go build -ldflags -H=windowsgui

## struct.dot

Is created when you unpack. Shows the structure of an application. Can be opened using [GraphViz](http://graphviz.org/download/)