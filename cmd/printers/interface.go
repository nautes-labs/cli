// Copyright 2023 Nautes Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package printers

import (
	"fmt"
	"github.com/nautes-labs/cli/cmd/types"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

// terminalEscaper replaces ANSI escape sequences and other terminal special
// characters to avoid terminal escape character attacks (issue #101695).
var terminalEscaper = strings.NewReplacer("\x1b", "^[", "\r", "\\r")

// WriteEscaped replaces unsafe terminal characters with replacement strings
// and writes them to the given writer.
func WriteEscaped(writer io.Writer, output string) error {
	_, err := terminalEscaper.WriteString(writer, output)
	return err
}

type mergeTo struct {
	from          string
	fromPrintName string
	to            string
	toPrintName   string
}

type Columns struct {
	PrintName string
	FieldName string
}

// GenerateTable was used to generate a table, To display by print tag which the value is true.
func GenerateTable(responseValues []reflect.Value, responseItemType reflect.Type) (*metav1.Table, error) {
	columns, mergeTos := generateColumns(responseItemType)
	//[Spec.Name Spec.Git.Gitlab.Name Spec.Git.Gitlab.Path Spec.Git.Gitlab.Visibility Spec.Git.Gitlab.Description]

	//build columns to display
	columnsDefinitions := buildPrintColumnsName(columns, mergeTos)

	//build table rows
	rows, err := buildTable(responseValues, columns, mergeTos)
	if err != nil {
		return nil, err
	}

	table := &metav1.Table{
		ColumnDefinitions: columnsDefinitions,
		Rows:              rows,
	}
	return table, nil
}

// generateColumns generates prefix column name which was set true on print tag.
func generateColumns(resourceType reflect.Type) (columns []*Columns, mergeTos []*mergeTo) {
	value := reflect.New(resourceType)
	vi := value.Elem().Interface()
	valueType := reflect.TypeOf(vi)
	columns, mergeTos = findDeepField(valueType)
	for _, val := range mergeTos {
		for _, column := range columns {
			if column.PrintName == val.toPrintName {
				val.to = column.FieldName
			}
		}
	}
	return
}

// findDeepField finds deep field name which was set true on print tag.
func findDeepField(sType reflect.Type) (columns []*Columns, mergeTos []*mergeTo) {
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		kind := field.Type.Kind()
		var fieldType = field.Type
		var fieldName = field.Name
		if kind == reflect.Ptr {
			prtValue := reflect.New(fieldType)
			realValue := reflect.Indirect(prtValue)
			fieldType = realValue.Type().Elem()
		}
		switch fieldType.Kind() {
		case reflect.Struct:
			columns, mergeTos = addPrefixFieldNameToDeepField(fieldType, fieldName, columns, mergeTos)
		case reflect.Slice:
			sliceValue := reflect.New(fieldType)
			realValue := reflect.Indirect(sliceValue)
			fieldType = realValue.Type().Elem()
			switch fieldType.Kind() {
			case reflect.Struct:
				columns, mergeTos = addPrefixFieldNameToDeepField(fieldType, fieldName, columns, mergeTos)
			default:
				column := field.Tag.Get(types.Column)
				mergeToColumn := field.Tag.Get(types.MergeTo)
				columns, mergeTos = dealTagForColumns(column, mergeToColumn, fieldName, columns, mergeTos)
			}
		case reflect.Map:
			continue
		default:
			column := field.Tag.Get(types.Column)
			mergeToColumn := field.Tag.Get(types.MergeTo)
			columns, mergeTos = dealTagForColumns(column, mergeToColumn, fieldName, columns, mergeTos)
		}
	}
	return
}

// addPrefixFieldNameToDeepField add parent field name.
func addPrefixFieldNameToDeepField(fieldType reflect.Type, fieldName string, columns []*Columns, mergeTos []*mergeTo) ([]*Columns, []*mergeTo) {
	cols, merges := findDeepField(fieldType)
	for _, col := range cols {
		col.FieldName = fieldName + "." + col.FieldName
	}
	for _, merge := range merges {
		merge.from = fieldName + "." + merge.from
		merge.to = fieldName + "." + merge.to
	}
	columns = append(columns, cols...)
	mergeTos = append(mergeTos, merges...)
	return columns, mergeTos
}

// dealTagForColumns build column and mergeTo slice by print column and mergeTo column.
func dealTagForColumns(printColumn, mergeToColumn, fieldName string, columns []*Columns, mergeTos []*mergeTo) ([]*Columns, []*mergeTo) {
	if printColumn != "" {
		columns = append(columns, &Columns{
			PrintName: printColumn,
			FieldName: fieldName,
		})
	}
	if mergeToColumn != "" {
		mergeTos = append(mergeTos, &mergeTo{
			from:          fieldName,
			fromPrintName: printColumn,
			to:            "",
			toPrintName:   mergeToColumn,
		})
	}
	return columns, mergeTos
}

// buildPrintColumnsName build columns to display.
func buildPrintColumnsName(columns []*Columns, mergeTos []*mergeTo) []metav1.TableColumnDefinition {
	columnsDefinitions := make([]metav1.TableColumnDefinition, 0, len(columns))
	checkSameColumnsName := make(map[string]struct{})
	for _, column := range columns {
		columnName := column.PrintName
		if _, ok := checkSameColumnsName[columnName]; ok {
			columnName = fmt.Sprintf("%s-2", columnName)
		}
		checkSameColumnsName[columnName] = struct{}{}

		var notPrintColumn string
		var fromPrintNames []string
		for _, val := range mergeTos {
			if val.toPrintName == columnName {
				fromPrintNames = append(fromPrintNames, val.fromPrintName)
			}
			if val.fromPrintName == columnName {
				notPrintColumn = val.from
			}
		}
		if notPrintColumn != "" {
			continue
		}
		if len(fromPrintNames) > 0 {
			columnName = fmt.Sprintf("%s / %s", columnName, strings.Join(fromPrintNames, " / "))
		}
		columnsDefinitions = append(columnsDefinitions, metav1.TableColumnDefinition{
			Name: columnName, Type: "string",
		})
	}
	return columnsDefinitions
}

// buildTable builds table which is includes table header and table row.
func buildTable(responseValues []reflect.Value, columns []*Columns, mergeTos []*mergeTo) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0)
	//rebuild column name
	var fieldColumns []string
	//[Name Git.Gitlab.Name Git.Gitlab.Path Git.Gitlab.Visibility Git.Gitlab.Description]
	for _, column := range columns {
		fieldColumns = append(fieldColumns, column.FieldName)
	}
	for _, value := range responseValues {
		itemRows, err := buildTableRows(value, fieldColumns, mergeTos)
		if err != nil {
			return nil, err
		}
		rows = append(rows, itemRows...)
	}
	return rows, nil
}

// buildTableRows builds table row from field columns name which is to display by stdout.
func buildTableRows(responseValue reflect.Value, columns []string, mergeTos []*mergeTo) ([]metav1.TableRow, error) {
	row := metav1.TableRow{}
	var rows []metav1.TableRow

	var mergeValueToColumnName = make(map[string][]string)
	var rebuildColumns []string
	for _, column := range columns {
		var skipMergeFromColumn bool
		var mergeFromValueArr []string
		var mergeToColumnName string
		for _, val := range mergeTos {
			if val.from == column { //skip current column value, it's from merge to column
				skipMergeFromColumn = true
				break
			} else if val.to == column {
				mergeFromValueArr = append(mergeFromValueArr, val.from)
				mergeToColumnName = val.to
			}
		}
		if skipMergeFromColumn {
			continue
		}
		rebuildColumns = append(rebuildColumns, column)
		//deal donot column
		columnNameArr := strings.Split(column, ".")
		columnValue := getValueByColumnName(responseValue, columnNameArr)

		//get from value to merge
		if len(mergeFromValueArr) > 0 {
			for _, fromValue := range mergeFromValueArr {
				columnNameArr = strings.Split(fromValue, ".")
				mergeValue := getValueByColumnName(responseValue, columnNameArr)
				mergeValueToColumnName[mergeToColumnName] = append(mergeValueToColumnName[mergeToColumnName], mergeValue)
			}
		}

		row.Cells = append(row.Cells, columnValue)
	}
	rows = append(rows, row)
	//deal empty columns
	if len(mergeValueToColumnName) > 0 {

		rowEmpty := metav1.TableRow{}
		cellsLen := len(rebuildColumns)
		for i := 0; i < cellsLen; i++ {
			if cellValues, ok := mergeValueToColumnName[rebuildColumns[i]]; ok {
				rowEmpty.Cells = append(rowEmpty.Cells, strings.Join(cellValues, ","))
			} else {
				rowEmpty.Cells = append(rowEmpty.Cells, "")
			}
		}
		rows = append(rows, rowEmpty)

		// add empty row
		rowEmptyOther := metav1.TableRow{}
		for i := 0; i < cellsLen; i++ {
			rowEmptyOther.Cells = append(rowEmptyOther.Cells, "")
		}
		rows = append(rows, rowEmptyOther)
	}
	return rows, nil
}

// getValueByColumnName gets column value by column name.
func getValueByColumnName(responseValue reflect.Value, column []string) (responseStr string) {
	if responseValue.Kind() == reflect.Invalid || responseValue.Kind() == reflect.Map {
		return ""
	}
	if responseValue.Kind() == reflect.Ptr && responseValue.IsNil() {
		return ""
	}
	if responseValue.Kind() == reflect.Ptr {
		responseValue = reflect.Indirect(responseValue)
	}
	if responseValue.Kind() == reflect.Slice {
		return getValueOfSliceByStruct(responseValue, column)
	}
	if !responseValue.FieldByName(column[0]).IsValid() {
		return ""
	}
	if len(column) == 1 {
		return getValueByColumnKind(responseValue, column[0])
	}
	responseValue = responseValue.FieldByName(column[0])
	column = column[1:]
	return getValueByColumnName(responseValue, column)
}

// getValueByColumnKind gets value which the kind is different.
func getValueByColumnKind(responseValue reflect.Value, column string) (responseStr string) {
	kind := responseValue.FieldByName(column).Kind()
	fieldValue := responseValue.FieldByName(column).Interface()
	switch kind {
	case reflect.String:
		responseStr = fieldValue.(string)
	case reflect.Slice:
		responseStr = getValueOfSlice(responseValue, column)
	case reflect.Bool:
		bRes := fieldValue.(bool)
		responseStr = fmt.Sprintf("%t", bRes)
	case reflect.Int:
		bRes := fieldValue.(int)
		responseStr = fmt.Sprintf("%d", bRes)
	case reflect.Int32:
		bRes := fieldValue.(int32)
		responseStr = fmt.Sprintf("%d", bRes)
	case reflect.Int64:
		bRes := fieldValue.(int64)
		responseStr = fmt.Sprintf("%d", bRes)
	case reflect.Float32:
		bRes := fieldValue.(float32)
		responseStr = fmt.Sprintf("%f", bRes)
	case reflect.Float64:
		bRes := fieldValue.(float64)
		responseStr = fmt.Sprintf("%f", bRes)
	default:
		responseStr = ""
	}
	return responseStr
}

// getValueOfSlice gets value when the column type is slice.
func getValueOfSlice(responseValue reflect.Value, column string) (responseStr string) {
	instanceLen := responseValue.FieldByName(column).Len()
	var result []string
	for i := 0; i < instanceLen; i++ {
		item := responseValue.FieldByName(column).Index(i)
		if item.Kind() != reflect.String {
			continue
		}
		result = append(result, item.String())
	}
	if len(result) > 5 {
		responseStr = fmt.Sprintf("%s...", strings.Join(result[0:5], ","))
	} else {
		responseStr = strings.Join(result, ",")
	}
	return responseStr
}

// getValueOfSliceByStruct gets value when the column type is slice.
func getValueOfSliceByStruct(responseValue reflect.Value, column []string) (responseStr string) {
	instanceLen := responseValue.Len()
	var result []string
	var duplicateValue = make(map[string]struct{})
	for i := 0; i < instanceLen; i++ {
		item := responseValue.Index(i)
		if item.Kind() == reflect.Struct {
			columnValue := getValueByColumnName(item, column)
			duplicateValue[columnValue] = struct{}{}
		}
	}

	for idx, _ := range duplicateValue {
		result = append(result, idx)
	}
	sort.Strings(result)
	if len(result) > 6 {
		responseStr = fmt.Sprintf("%s...", strings.Join(result[0:6], ","))
	} else {
		responseStr = strings.Join(result, ",")
	}
	return responseStr
}

// PrintTable prints a table to the provided output respecting the filtering rules for options
// for wide columns and filtered rows. It filters out rows that are Completed. You should call
// decorateTable if you receive a table from a remote server before calling printTable.
func PrintTable(table *metav1.Table, output io.Writer) error {
	if _, found := output.(*tabwriter.Writer); !found {
		w := GetNewTabWriter(output)
		output = w
		defer w.Flush()
	}

	first := true
	for _, column := range table.ColumnDefinitions {
		if first {
			first = false
		} else {
			fmt.Fprint(output, "\t")
		}
		fmt.Fprint(output, strings.ToUpper(column.Name))
	}
	fmt.Fprintln(output)
	for _, row := range table.Rows {
		first := true
		for i, cell := range row.Cells {
			if i >= len(table.ColumnDefinitions) {
				break
			}
			if first {
				first = false
			} else {
				fmt.Fprint(output, "\t")
			}
			if cell != nil {
				switch val := cell.(type) {
				case string:
					print := val
					truncated := false
					breakchar := strings.IndexAny(print, "\f\n\r")
					if breakchar >= 0 {
						truncated = true
						print = print[:breakchar]
					}
					WriteEscaped(output, print)
					if truncated {
						fmt.Fprint(output, "...")
					}
				default:
					WriteEscaped(output, fmt.Sprint(val))
				}
			}
		}
		fmt.Fprintln(output)
	}
	return nil
}
