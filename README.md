# apla packager

Simvolio and Protypo files packager to import on ecosystem

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

        work log

  -c string

        --conditions (default "ContractConditions(\"MainCondition\")")

  -i string

        --input (default ".")

  -m string

        --menu (default "default_menu")

  -o string

        --output (default "output")

  -p string

        --prefix

  -t string

        --table-permission (default "{\"insert\":\"true\",\"update\":\"true\",
        \"new_column\":\"true\"}")
  -u    --unpack

  -v    --verbose
  

### Examples

#### Unpack file

>ap -u -i basic.sim

#### Pack files in dir

>cd output
>ap -o basic.sim