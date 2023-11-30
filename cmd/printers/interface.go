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

// GenerateTable creates a metav1.Table structure from a slice of reflect.Values and the reflection type of the item.
// It generates table columns and rows based on the provided data and returns the resulting metav1.Table.
func GenerateTable(responseValues []reflect.Value, responseItemType reflect.Type) (*metav1.Table, error) {
	columns, mergeTos := generateColumns(responseItemType)
	//[Spec.Name Spec.Git.Gitlab.Name Spec.Git.Gitlab.Path Spec.Git.Gitlab.Visibility Spec.Git.Gitlab.Description]

	//build columns to display
	columnsDefinitions := buildPrintColumnsName(columns, mergeTos)

	//build table rows
	rows := buildTable(responseValues, columns, mergeTos)

	table := &metav1.Table{
		ColumnDefinitions: columnsDefinitions,
		Rows:              rows,
	}
	return table, nil
}

// generateColumns extracts Columns and mergeTo specifications from the reflection type of the resource.
// It returns the list of Columns and mergeTo specifications for further processing.
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

// findDeepField recursively explores the fields of a reflection type and extracts Columns and mergeTo specifications.
// It returns the list of Columns and mergeTo specifications for further processing.
func findDeepField(sType reflect.Type) (columns []*Columns, mergeTos []*mergeTo) {
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		kind := field.Type.Kind()
		var fieldType = field.Type
		var fieldName = field.Name
		column := field.Tag.Get(types.Column)
		mergeToColumn := field.Tag.Get(types.MergeTo)
		if kind == reflect.Struct {
			if column != "" {
				if strings.Contains(column, ":") {
					columnArr := strings.Split(column, ":")
					column = columnArr[0]
					if len(columnArr) > 1 {
						fieldName = fmt.Sprintf("%s:%s", fieldName, columnArr[1])
					}
				}
				columns, mergeTos = dealTagForColumns(column, mergeToColumn, fieldName, columns, mergeTos)
				continue
			}
		}
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
				columns, mergeTos = dealTagForColumns(column, mergeToColumn, fieldName, columns, mergeTos)
			}
		case reflect.Map:
			continue
		default:
			columns, mergeTos = dealTagForColumns(column, mergeToColumn, fieldName, columns, mergeTos)
		}
	}
	return
}

// addPrefixFieldNameToDeepField adds a prefix to the field names extracted from a deep structure.
// It is used in the context of nested structures to maintain proper field hierarchy.
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

// dealTagForColumns processes the print and mergeTo tags for a field and updates the corresponding lists.
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

// buildPrintColumnsName generates table column definitions based on a list of Columns and mergeTo specifications.
// It ensures unique column names and handles cases where columns need to be merged or skipped.
// The resulting table column definitions are returned as a slice of metav1.TableColumnDefinition.
func buildPrintColumnsName(columns []*Columns, mergeTos []*mergeTo) []metav1.TableColumnDefinition {
	// Initialize an empty slice to store column definitions
	columnsDefinitions := make([]metav1.TableColumnDefinition, 0, len(columns))

	// Map to track unique column names and handle duplicates
	checkSameColumnsName := make(map[string]struct{})

	// Iterate through each column
	for _, column := range columns {
		// Get the print name of the column
		columnName := column.PrintName

		// Check for duplicate column names and append a suffix if necessary
		if _, ok := checkSameColumnsName[columnName]; ok {
			columnName = fmt.Sprintf("%s-2", columnName)
		}
		checkSameColumnsName[columnName] = struct{}{}

		// Variables to track merging and skipping of columns
		var notPrintColumn string
		var fromPrintNames []string

		// Iterate through mergeTo specifications to handle merging and skipping
		for _, val := range mergeTos {
			if val.toPrintName == columnName {
				fromPrintNames = append(fromPrintNames, val.fromPrintName)
			}
			if val.fromPrintName == columnName {
				notPrintColumn = val.from
			}
		}

		// Skip the column if it is marked as not to be printed
		if notPrintColumn != "" {
			continue
		}

		// If merging is needed, update the column name
		if len(fromPrintNames) > 0 {
			columnName = fmt.Sprintf("%s / %s", columnName, strings.Join(fromPrintNames, " / "))
		}

		// Append the column definition to the result
		columnsDefinitions = append(columnsDefinitions, metav1.TableColumnDefinition{
			Name: columnName, Type: "string",
		})
	}

	// Return the resulting table column definitions
	return columnsDefinitions
}

// buildTable builds table which is includes table header and table row.
func buildTable(responseValues []reflect.Value, columns []*Columns, mergeTos []*mergeTo) []metav1.TableRow {
	rows := make([]metav1.TableRow, 0)
	//rebuild column name
	var fieldColumns []string
	//[Name Git.Gitlab.Name Git.Gitlab.Path Git.Gitlab.Visibility Git.Gitlab.Description]
	for _, column := range columns {
		fieldColumns = append(fieldColumns, column.FieldName)
	}
	for _, value := range responseValues {
		itemRows := buildTableRows(value, fieldColumns, mergeTos)
		rows = append(rows, itemRows...)
	}
	return rows
}

// buildTableRows builds table row from field columns name which is to display by stdout.
func buildTableRows(responseValue reflect.Value, columns []string, mergeTos []*mergeTo) []metav1.TableRow {
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
	return rows
}

// getValueByColumnName retrieves the value from a nested structure using reflection
// based on the provided column names. It handles cases such as structs, slices, and pointers.
// It returns an empty string if the value is not found or if the input is invalid.
func getValueByColumnName(responseValue reflect.Value, column []string) (responseStr string) {
	// Check if the responseValue is invalid or a map, return an empty string in such cases
	if responseValue.Kind() == reflect.Invalid || responseValue.Kind() == reflect.Map {
		return ""
	}

	// If responseValue is a nil pointer, return an empty string
	if responseValue.Kind() == reflect.Ptr && responseValue.IsNil() {
		return ""
	}

	// If responseValue is a pointer, dereference it
	if responseValue.Kind() == reflect.Ptr {
		responseValue = reflect.Indirect(responseValue)
	}

	// If responseValue is a slice, delegate to getValueOfSliceByStruct
	if responseValue.Kind() == reflect.Slice {
		return getValueOfSliceByStruct(responseValue, column)
	}

	// Check if the first column name is valid and not a nested field
	if !responseValue.FieldByName(column[0]).IsValid() && !strings.Contains(column[0], ":") {
		return ""
	}

	// If there is only one column, extract the column and subfield names
	if len(column) == 1 {
		var columnName, subFieldName string
		if strings.Contains(column[0], ":") {
			columnArr := strings.Split(column[0], ":")
			columnName = columnArr[0]
			if len(columnArr) > 1 {
				subFieldName = columnArr[1]
			}
		} else {
			columnName = column[0]
		}
		return getValueByColumnKind(responseValue, columnName, subFieldName)
	}

	// Recursively call the function with the next level of nested structure
	responseValue = responseValue.FieldByName(column[0])
	column = column[1:]
	return getValueByColumnName(responseValue, column)
}

// getValueByColumnKind extracts and converts the value of a specific column from a structure
// based on its kind (type) using reflection. It also handles subfields in case of nested structures.
// It returns a string representation of the value.
func getValueByColumnKind(responseValue reflect.Value, column, subFieldName string) (responseStr string) {
	// Check if the responseValue is invalid, return an empty string if true
	if responseValue.Kind() == reflect.Invalid {
		return ""
	}

	// Obtain the kind (type) of the specified column in the structure
	kind := responseValue.FieldByName(column).Kind()

	// Get the interface{} representation of the column's value
	fieldValue := responseValue.FieldByName(column).Interface()

	// Switch based on the kind of the column and handle each type accordingly
	switch kind {
	case reflect.String:
		responseStr = fieldValue.(string)
	case reflect.Struct:
		// If the column is a struct, delegate to getValueOfStruct for further processing
		responseStr = getValueOfStruct(responseValue, column, subFieldName)
	case reflect.Slice:
		// If the column is a slice, delegate to getValueOfSlice for further processing
		responseStr = getValueOfSlice(responseValue, column)
	case reflect.Bool:
		// Convert bool to string representation
		bRes := fieldValue.(bool)
		responseStr = fmt.Sprintf("%t", bRes)
	case reflect.Int:
		// Convert int to string representation
		bRes := fieldValue.(int)
		responseStr = fmt.Sprintf("%d", bRes)
	case reflect.Int32:
		// Convert int32 to string representation
		bRes := fieldValue.(int32)
		responseStr = fmt.Sprintf("%d", bRes)
	case reflect.Int64:
		// Convert int64 to string representation
		bRes := fieldValue.(int64)
		responseStr = fmt.Sprintf("%d", bRes)
	case reflect.Float32:
		// Convert float32 to string representation
		bRes := fieldValue.(float32)
		responseStr = fmt.Sprintf("%f", bRes)
	case reflect.Float64:
		// Convert float64 to string representation
		bRes := fieldValue.(float64)
		responseStr = fmt.Sprintf("%f", bRes)
	default:
		// If the kind is not recognized, set the responseStr to an empty string
		responseStr = ""
	}

	// Return the string representation of the column's value
	return responseStr
}

// getValueOfSlice extracts and processes the values of a slice column from a structure using reflection.
// It retrieves the string representation of each string element in the slice, up to a maximum of 5 elements,
// concatenates them with commas, and returns the resulting string. If there are more than 5 elements,
// the string is truncated and appended with "...".
func getValueOfSlice(responseValue reflect.Value, column string) (responseStr string) {
	// Get the length of the slice column
	instanceLen := responseValue.FieldByName(column).Len()

	// Initialize a slice to store string representations of each element
	var result []string

	// Iterate through each element in the slice
	for i := 0; i < instanceLen; i++ {
		// Get the i-th element from the slice
		item := responseValue.FieldByName(column).Index(i)

		// Check if the element is of kind string, skip if not
		if item.Kind() != reflect.String {
			continue
		}

		// Append the string representation of the element to the result slice
		result = append(result, item.String())
	}

	// If there are more than 5 elements, truncate the result and append "..."
	if len(result) > 5 {
		responseStr = fmt.Sprintf("%s...", strings.Join(result[0:5], ","))
	} else {
		// If there are 5 or fewer elements, concatenate them with commas
		responseStr = strings.Join(result, ",")
	}

	// Return the resulting string representation of the slice column
	return responseStr
}

// getValueOfSliceByStruct extracts and processes unique values of a slice of structures based on the specified column names.
// It uses reflection to navigate the nested structures and retrieves the values of the specified columns for each structure.
// The unique values are then sorted and concatenated with commas, and if there are more than 5 unique values, the result
// is truncated and appended with "...".
func getValueOfSliceByStruct(responseValue reflect.Value, column []string) (responseStr string) {
	// Get the length of the slice
	instanceLen := responseValue.Len()

	// Initialize a slice to store unique string representations of specified columns
	var result []string

	// Use a map to track unique values to avoid duplicates
	var duplicateValue = make(map[string]struct{})

	// Iterate through each element in the slice
	for i := 0; i < instanceLen; i++ {
		// Get the i-th element from the slice
		item := responseValue.Index(i)

		// Check if the element is of kind struct
		if item.Kind() == reflect.Struct {
			// Get the value of the specified column for the current structure
			columnValue := getValueByColumnName(item, column)

			// Store the unique column values in the map
			duplicateValue[columnValue] = struct{}{}
		}
	}

	// Iterate through unique values and append them to the result slice
	for idx := range duplicateValue {
		result = append(result, idx)
	}

	// Sort the result to maintain order
	sort.Strings(result)

	// If there are more than 5 unique values, truncate the result and append "..."
	if len(result) > 5 {
		responseStr = fmt.Sprintf("%s...", strings.Join(result[0:5], ","))
	} else {
		// If there are 5 or fewer unique values, concatenate them with commas
		responseStr = strings.Join(result, ",")
	}

	// Return the resulting string representation of the unique values in the slice of structures
	return responseStr
}

// getValueOfStruct extracts and processes the values of a specified column within a struct using reflection.
// It iterates through the fields of the specified column, optionally considering a subfield, and retrieves
// the string representation of each field. The unique values are then sorted and concatenated with commas,
// and if there are more than 5 unique values, the result is truncated and appended with "...".
// The final result includes the count of unique values within square brackets.
func getValueOfStruct(responseValue reflect.Value, column, subFieldName string) (responseStr string) {
	// Get the value of the specified column within the struct
	value := responseValue.FieldByName(column)

	// Initialize a slice to store unique string representations of field values
	var result []string

	// Iterate through each field in the specified column
	for i := 0; i < value.NumField(); i++ {
		// Get the i-th field
		field := value.Field(i)

		// Get the name of the field (considering subfields)
		name := getValueByColumnKind(field.Elem(), subFieldName, subFieldName)

		// If the name is not empty, append it to the result slice
		if name != "" {
			result = append(result, name)
		}
	}

	// Sort the result to maintain order
	sort.Strings(result)

	// If there are more than 5 unique values, truncate the result and append "..."
	if len(result) > 5 {
		responseStr = fmt.Sprintf("%s...", strings.Join(result[0:5], ","))
	} else {
		// If there are 5 or fewer unique values, concatenate them with commas
		responseStr = strings.Join(result, ",")
	}

	// Format the final response string to include the count of unique values within square brackets
	responseStr = fmt.Sprintf("[%d] %s", len(result), responseStr)

	// Return the resulting string representation of the values within the struct
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
					printable := val
					truncated := false
					breakchar := strings.IndexAny(printable, "\f\n\r")
					if breakchar >= 0 {
						truncated = true
						printable = printable[:breakchar]
					}
					_ = WriteEscaped(output, printable)
					if truncated {
						fmt.Fprint(output, "...")
					}
				default:
					_ = WriteEscaped(output, fmt.Sprint(val))
				}
			}
		}
		fmt.Fprintln(output)
	}
	return nil
}
