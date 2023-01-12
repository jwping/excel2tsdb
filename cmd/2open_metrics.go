package cmd

import (
	"excel2tsdb/internal/elog"
	"excel2tsdb/internal/util"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var (
	omCommand = &cobra.Command{
		Use:     "openmetrics",
		Aliases: []string{"om", "oms", "2om", "2oms"},
		Short:   "Convert excel to OpenMetrics format",
		Long: `Optional output to specified file (default print to screen).
If the specified path is a directory, then it loops through all supported format files(which is also optional).`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := toOpenMetricsRun(); err != nil {
				elog.Log.Error("openmetrics run faild", err)
			}
		},
	}

	outPath *string
)

func init() {
	outPath = omCommand.PersistentFlags().StringP("outpath", "o", "", "Specifies the file to output, or output to the screen if not specified.\nNote that this action will overwrite the specified file. Please exercise caution!")
	rootCmd.AddCommand(omCommand)
}

func toOpenMetricsRun() error {
	dataMetrics, err := util.GetMetrics(*excelPath, *deep, *worksheets, *excludedRows, *timestampColumn, *dataColumn, *formatTime, *millisecondTimestamp)
	if err != nil {
		elog.Log.Error("GetMetrics faild", err)
		return err
	}

	metricData := util.DataIntegration(*fqName, dataMetrics)

	if *outPath == "" {
		fmt.Println(string(metricData))
		return nil
	}

	if err = ioutil.WriteFile(*outPath, metricData, 0666); err != nil {
		elog.Log.Error("write file faild", err, "outpath", *outPath)
		return err
	}

	return nil
}
