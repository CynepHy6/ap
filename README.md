# apla packager

Simvolio and Protypo files packager to import on ecosystem

### Examples

#### Unpack file from "basic.sim" to "output/"

>ap -u -i basic.sim

#### Pack files from "output/" to basic.sim.json

      cd output

      ap -o basic.sim

### Usage of "ap"

  --conditions string

      -c, conditions (default "ContractConditions(\"MainCondition\")")

  --input string

      -i, path for input files (default ".")

  --menu string

      -m, menu (default "default_menu")

  --output string

      -o, output filename for JSON (default "output")

  --prefix string

      -p, prefix for pages and contracts

  --table-permission string

      -t, permission for tables (default "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}")

--unpack
      -u, unpacking mode

--verbose
      -v, work log