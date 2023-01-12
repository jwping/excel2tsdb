package cmd

import (
	"context"
	"excel2tsdb/internal/elog"
	"excel2tsdb/internal/util"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	tCommand = &cobra.Command{
		Use:     "tsdb",
		Aliases: []string{"tsdb", "prometheus", "2tsdb", "2ts"},
		Short:   "excel converted to prometheus tsdb data",
		Long: `The converted prometheus tsdb data files can be placed directly in the prometheus data directory.
Note that prometheus needs to be restarted after importing the tsdb, and historical data may be cleaned up due to the data retention time set by prometheus (default: 15d) if it is too late. Specific see Prometheus "--storage.tsdb.retention". The time parameter.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := toTsdbRun(); err != nil {
				elog.Log.Error("tsdb run faild", err)
			}
		},
	}

	outDir       *string
	timeout      *int
	promtoolPath *string
)

func init() {
	promtoolPath = tCommand.PersistentFlags().StringP("promtoolpath", "p", "/usr/bin/promtool", "Specify the promtool path")
	outDir = tCommand.PersistentFlags().StringP("outdir", "o", "./", "Specify the directory where Prometheus tsdb data will be exported (default output to current directory)")
	timeout = tCommand.PersistentFlags().Int("timeout", 30, "Timeout period converted to tsdb format (unit: minute)")
	rootCmd.AddCommand(tCommand)
}

func toTsdbRun() error {
	dataMetrics, err := util.GetMetrics(*excelPath, *deep, *worksheets, *excludedRows, *timestampColumn, *dataColumn, *formatTime, *millisecondTimestamp)
	if err != nil {
		return err
	}

	metricData := util.DataIntegration(*fqName, dataMetrics)

	f, err := ioutil.TempFile(os.TempDir(), ".openmetrics.tmp")

	if err != nil {
		elog.Log.Error("Failed to create a temporary file. Procedure", err)
		return err
	}

	_, err = f.Write(metricData)
	if err != nil {
		elog.Log.Error("Failed to write data to the temporary file", err)
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*time.Duration(*timeout))
	cmd := exec.CommandContext(ctx, *promtoolPath, "tsdb", "create-blocks-from", "openmetrics", f.Name(), *outDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		elog.Log.Error("Description Failed to convert to tsdb", err, "tempfile", f.Name())
		return err
	}

	elog.Log.Info("Export successfully! Manually copy all tsdb data in the export directory to the data directory of Prometheus and restart it.", "metric_name", *fqName)

	if err = f.Close(); err != nil {
		elog.Log.Error("tempfile clost faild", err, "tempfile", f.Name())
		return err
	}

	if err = os.Remove(f.Name()); err != nil {
		elog.Log.Warn("remove tempfile faild", "error", err.Error(), "path", f.Name())
	}

	return nil
}
