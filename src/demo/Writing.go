package main

import (
	"fmt"
	"excel"
)

func main() {

	categories := map[string]string{"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"}
	values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
	xlsx := excel.NewFile()

	for k, v := range categories {
		fmt.Println(k)
		fmt.Println(v)
		xlsx.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		xlsx.SetCellValue("Sheet1", k, v)
	}
	// Save xlsx file by the given path.
	err := xlsx.SaveAs("d:/Workbook.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}
