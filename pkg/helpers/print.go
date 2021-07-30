// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

const (
	YamlFormat            string = "yaml"
	JsonFormat            string = "json"
	CustomColumnsFormat   string = "columns="
	ColumnsSeparator      string = ","
	SupportedOutputFormat string = YamlFormat + "|" + JsonFormat + "|" + CustomColumnsFormat + "..."
)

func IsOutputFormatSupported(format string) bool {
	return len(format) == 0 ||
		strings.HasPrefix(format, CustomColumnsFormat) ||
		strings.ToLower(format) == YamlFormat ||
		strings.ToLower(format) == JsonFormat
}

func Print(o interface{}, format string, noHeaders bool, f func(interface{}) ([]map[string]string, error)) error {
	switch strings.ToLower(format) {
	case YamlFormat:
		return printYaml(o)
	case JsonFormat:
		return printJson(o)
	default:
		m, err := f(o)
		if err != nil {
			return err
		}
		return printText(m, format, noHeaders)
	}
}

func printYaml(o interface{}) error {
	b, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}

func printJson(o interface{}) error {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func printText(o []map[string]string, format string, noHeaders bool) error {
	format = strings.TrimPrefix(format, CustomColumnsFormat)
	lines, err := printArrayColumns(o, format, noHeaders)
	if err != nil {
		return err
	}
	printLines(lines)
	return nil
}

func printArrayColumns(o []map[string]string, format string, noHeaders bool) ([]string, error) {
	lines := make([]string, 0)
	if !noHeaders {
		line, err := generateLine(o, format, true)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	for _, elem := range o {
		line, err := generateLine(elem, format, false)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func generateLine(o interface{}, format string, header bool) (string, error) {
	formatColumns := strings.Split(format, ColumnsSeparator)
	columns := make([]string, 0)
	for _, c := range formatColumns {
		if header {
			columns = append(columns, c)
		} else {
			columns = append(columns, o.(map[string]string)[c])
		}
	}
	return strings.Join(columns, ColumnsSeparator), nil
}

func printLines(lines []string) {
	if len(lines) == 0 {
		return
	}
	columns := strings.Split(lines[0], ColumnsSeparator)
	ncColumns := len(columns)
	if ncColumns == 0 {
		return
	}
	columnsSize := make([]int, ncColumns)
	for _, l := range lines {
		columns := strings.Split(l, ColumnsSeparator)
		for i, col := range columns {
			if len(col) > columnsSize[i] {
				columnsSize[i] = len(col)
			}
		}
	}
	for _, l := range lines {
		columns := strings.Split(l, ColumnsSeparator)
		var outLine string
		for i, col := range columns {
			formatString := "%-" + fmt.Sprintf("%d", columnsSize[i]) + "s "
			outLine = outLine + fmt.Sprintf(formatString, col)
		}
		fmt.Println(outLine)
	}
}
