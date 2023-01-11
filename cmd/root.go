package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "excel2tsdb",
		Short: "Convert excel files to Prometheus storage or export format",
		Long: `Supports conversion to OpenMetrics or direct conversion to tsdb data.
Supported format: XLAM/XLSM/XLSX/XLTM/XLTX.
Note that the specified excel file must have time and data columns exported in the metric format counter.
`,
	}

	excelPath            *string
	deep                 *bool
	worksheets           *[]string
	excludedRows         *[]int
	timestampColumn      *int
	dataColumn           *int
	fqName               *string
	formatTime           *string
	millisecondTimestamp *bool
)

func Execute() {

	excelPath = rootCmd.PersistentFlags().String("path", "", "Specifies the excel file path or directory")
	fqName = rootCmd.PersistentFlags().String("metric", "excel2tsdb", "Specify the name of the metric, default is (excel2tsdb)")
	deep = rootCmd.PersistentFlags().Bool("deep", true, "Enable directory recursive search, ignore this if the specified path is excel (enabled by default)	")
	worksheets = rootCmd.PersistentFlags().StringArray("worksheets", []string{"Sheet1"}, "Specifies an array of worksheets")
	excludedRows = rootCmd.PersistentFlags().IntSlice("excludedRows", nil, "Specify the rows to exclude (Start at 1)")
	timestampColumn = rootCmd.PersistentFlags().IntP("timestampColumn", "t", 1, "Specifies the time column (Start at 1)")
	dataColumn = rootCmd.PersistentFlags().IntP("dataColumn", "d", 2, "Specifies the data column (Start at 1)")
	formatTime = rootCmd.PersistentFlags().String("formatTime", "", "If the timestamp is listed as a string time format, specify the layout format (for example, 2006-01-02 15:04:05, please follow golang's time layout, that is, the date year month day must be 2006 01, 02, and the time must be 15:04:05, it can also look like this: 2006 01-02 15.04.05)")
	millisecondTimestamp = rootCmd.PersistentFlags().Bool("millisecondTimestamp", false, "Specifies whether the timestamp column is millisecond level (default is false)")

	rootCmd.MarkPersistentFlagRequired("path")
	rootCmd.Execute()
}
