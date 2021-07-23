// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"fmt"
	"strings"
)

func PrintLines(lines []string, sep string) {
	if len(lines) == 0 {
		return
	}
	columns := strings.Split(lines[0], sep)
	ncColumns := len(columns)
	if ncColumns == 0 {
		return
	}
	columnsSize := make([]int, ncColumns)
	for _, l := range lines {
		columns := strings.Split(l, sep)
		for i, col := range columns {
			if len(col) > columnsSize[i] {
				columnsSize[i] = len(col)
			}
		}
	}
	for _, l := range lines {
		columns := strings.Split(l, sep)
		var outLine string
		for i, col := range columns {
			formatString := "%-" + fmt.Sprintf("%d", columnsSize[i]) + "s "
			outLine = outLine + fmt.Sprintf(formatString, col)
		}
		fmt.Println(outLine)
	}
}
