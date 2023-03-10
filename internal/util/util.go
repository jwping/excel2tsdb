package util

import (
	"bytes"
	"excel2tsdb/internal/elog"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Get all legitimate excel files from a file or directory
func getExcel(excelPath string, deep bool) ([]string, error) {
	fi, err := os.Stat(excelPath)
	if err != nil {
		return nil, err
	}

	var (
		excelfiles   = make([]string, 0)
		supportedMap = map[string]struct{}{
			"XLAM": {},
			"XLSM": {},
			"XLSX": {},
			"XLTM": {},
			"XLTX": {},
		}
	)

	checkExcel := func(excelPath string) bool {
		excelfilenames := strings.Split(excelPath, ".")
		if len(excelfilenames) != 2 {
			elog.Log.Warn("excel name irregularity!", "filepath", excelPath)
			return false
		}
		suffix := strings.ToUpper(strings.ToLower(excelfilenames[1]))
		if _, ok := supportedMap[suffix]; ok {
			return true
		}

		return false
	}

	if fi.IsDir() {
		files, err := os.ReadDir(excelPath)
		if err != nil {
			elog.Log.Error("read dir faild", err, "dirpath", excelPath)
			return nil, err
		}

		for _, file := range files {
			if file.IsDir() {
				if deep {
					efs, err := getExcel(path.Join(excelPath, file.Name()), deep)
					if err != nil {
						elog.Log.Warn("in GetExcel faild", "error", err.Error(), "path", path.Join(excelPath, file.Name()))
						continue
					}
					excelfiles = append(excelfiles, efs...)
				}
			} else if checkExcel(path.Join(excelPath, file.Name())) {
				excelfiles = append(excelfiles, path.Join(excelPath, file.Name()))
			}
		}
	} else {
		if checkExcel(excelPath) {
			excelfiles = append(excelfiles, excelPath)
		}
	}

	return excelfiles, nil
}

type dataItem [2]interface{}

func getExcelMetrics(excelPath string, worksheets []string, excludedRowMaps map[int]struct{}, timestampColumn int, dataColumn int, formatTime string, millisecondTimestamp bool) ([]dataItem, error) {
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		elog.Log.Warn("excel open faild", "error", err.Error(), "path", excelPath)
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			elog.Log.Warn("excelfile close faild", "error", err.Error(), "path", excelPath)
			return
		}
	}()

	metricDatas := make([]dataItem, 0)

	for _, worksheet := range worksheets {
		// Get all the rows in the Sheet1.
		rows, err := f.GetRows(worksheet)
		if err != nil {
			elog.Log.Warn("excel GetRows faild", "error", err.Error(), "path", excelPath, "worksheet", worksheet)
			continue
		}

		for i, row := range rows {
			if _, ok := excludedRowMaps[i]; !ok {
				if row[dataColumn-1] == "" || row[timestampColumn-1] == "" {
					elog.Log.Warn("There is a null value, ignore it", "row number", i, "dataColumn", row[dataColumn-1], "timestampColumn", row[timestampColumn-1])
					continue
				}

				var timestamp int64
				if formatTime != "" {
					timeP, err := time.Parse(formatTime, row[timestampColumn-1])
					if err != nil {
						elog.Log.Warn("time layout Parse faild", "error", err.Error(), "data", row[timestampColumn-1])
						continue
					}

					timestamp = timeP.UnixMilli()
					// timestamp = fmt.Sprintf("%d.%d\n", timestampI/1000, timestampI%1000)
				} else {
					timestamp, err = strconv.ParseInt(row[timestampColumn-1], 10, 64)
					if err != nil {
						elog.Log.Warn("millisecond timestamp to int64 faild", "error", err.Error(), "data", row[timestampColumn-1])
						continue
					}

					if !millisecondTimestamp {
						timestamp *= 1000
						// timestamp = fmt.Sprintf("%d.%d\n", timestampI/1000, timestampI%1000)
					}
				}

				metricDatas = append(metricDatas, dataItem{row[dataColumn-1], timestamp})
			}
		}
	}

	return metricDatas, nil
}

func GetMetrics(excelPath string, deep bool, worksheets []string, excludedRows []int, timestampColumn int, dataColumn int, formatTime string, millisecondTimestamp bool) ([]dataItem, error) {
	excelPaths, err := getExcel(excelPath, deep)
	if err != nil {
		return nil, err
	}

	excludedRowMaps := make(map[int]struct{})
	if excludedRows != nil && len(excludedRows) != 0 {
		for _, excludedRow := range excludedRows {
			excludedRowMaps[excludedRow-1] = struct{}{}
		}
	}

	if worksheets == nil || len(worksheets) == 0 {
		worksheets = []string{"Sheet1"}
	}

	metricDatas := make([]dataItem, 0)
	for _, excelPath := range excelPaths {
		metricData, err := getExcelMetrics(excelPath, worksheets, excludedRowMaps, timestampColumn, dataColumn, formatTime, millisecondTimestamp)
		if err != nil {
			elog.Log.Warn("getExcelMetrics faild", "error", err.Error(), "path", excelPath)
			continue
		}
		metricDatas = append(metricDatas, metricData...)
	}

	// sort
	sort.SliceStable(metricDatas, func(i, j int) bool {
		return metricDatas[i][1].(int64) < metricDatas[j][1].(int64)
	})

	return metricDatas, nil
}

func DataIntegration(fqName string, dataItems []dataItem) []byte {
	outBuf := bytes.Buffer{}

	outBuf.WriteString(fmt.Sprintf("# HELP %s Number Automatically generated by excel2tsdb.\n# TYPE %s counter\n", fqName, fqName))

	for _, dataItem := range dataItems {
		outBuf.WriteString(fmt.Sprintf("%s{generate=\"excel2tsdb\"} %s %d.%d\n", fqName, dataItem[0].(string), dataItem[1].(int64)/1000, dataItem[1].(int64)%1000))
	}
	outBuf.WriteString("# EOF")

	return outBuf.Bytes()
}
