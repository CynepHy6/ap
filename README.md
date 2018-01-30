# apla packager

Simvolio and Protypo files packager to import on ecosystem

### Examples

#### Unpack file from "basic.sim" to "./basic/"

>ap -u -i basic.sim

#### Pack files from "basic/" to ./basic.json

>ap -i basic/

### Usage of "ap"

--input string

    -i, path for input files (default ".")

--output string

    -o, output filename for JSON (default "output" if "input" not found)

--unpack

    -u, unpacking mode


for view all flags please use "ap -h"