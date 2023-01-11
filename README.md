# excel2tsdb

> Convert and export excel file data to OpenMetrics format and support tsdb data export through promtool

## usage:
```shell
$ go build

$ ./excel2tsdb
Supports conversion to OpenMetrics or direct conversion to tsdb data.
Supported format: XLAM/XLSM/XLSX/XLTM/XLTX.
Note that the specified excel file must have time and data columns exported in the metric format counter.

Usage:
  excel2tsdb [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  openmetrics Convert excel to OpenMetrics format
  openmetrics Convert excel to OpenMetrics format

Flags:
  -d, --dataColumn int           Specifies the data column (Start at 1) (default 2)
      --deep                     Enable directory recursive search, ignore this if the specified path is excel (enabled by default)      (default true)
      --excludedRows ints        Specify the rows to exclude (Start at 1)
      --formatTime string        If the timestamp is listed as a string time format, specify the layout format (for example, 2006-01-02 15:04:05, please follow golang's time layout, that is, the date year month day must be 2006 01, 02, and the time must be 15:04:05, it can also look like this: 2006 01-02 15.04.05)
  -h, --help                     help for excel2tsdb
      --metric string            Specify the name of the metric, default is (excel2tsdb) (default "excel2tsdb")
      --millisecondTimestamp     Specifies whether the timestamp column is millisecond level (default is false)
      --path string              Specifies the excel file path or directory
  -t, --timestampColumn int      Specifies the time column (Start at 1) (default 1)
      --worksheets stringArray   Specifies an array of worksheets (default [Sheet1])

Use "excel2tsdb [command] --help" for more information about a command.
```