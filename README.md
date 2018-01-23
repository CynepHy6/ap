# apla packager

Simvolio and Protypo files packager to import on ecosystem

### Examples

#### Unpack file from "basic.sim" to "./basic/"

>ap -u -i basic.sim

#### Pack files from "basic/" to ./basic.json

>ap -i basic/

### Usage of "ap"

  --conditions string

      -c, conditions (default "ContractConditions(\"MainCondition\")")

  --input string

      -i, path for input files (default ".")

  --menu string

      -m, menu (default "default_menu")

  --output string

      -o, output filename for JSON (default "output")


  --table-permission string

      -t, permission for tables (default "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}")

--unpack

      -u, unpacking mode

--verbose

      work log

for actuals flags use "ap -h"