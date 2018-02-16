# apla packager

Simvolio and Protypo files packager for import to ecosystem.
Utilite that can convert import json bundle from/to files of simvolio, protypo, csv, json.

## Examples

### Unpack file from "basic.sim" to "./basic/"

>ap basic.sim

### Pack files from "basic/" to ./basic.json

>ap basic/

## build

### linux

    go build

### on windows

    go build -ldflags -H=windowsgui

### on linux for windows

    env GOARCH=amd64 GOOS=windows CGO_ENABLED=1 CC=/usr/bin/x86_64-w64-mingw32-gcc CXX=/usr/bin/x86_64-w64-mingw32-g++  go build -ldflags -H=windowsgui

## struct.dot

Is created when you pack/unpack. Shows the structure of an application. Can be opened using [graphviz](http://graphviz.org/download/) or [webgraphviz](http://webgraphviz.com/)
