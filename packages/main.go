package main

import (
	"fmt"

	"./carbon"
)

func main() {
	fmt.Println(carbon.Now().MonthName())
	fmt.Println(carbon.Now().PreviousMonthStartDay().MonthName())

	fmt.Println(carbon.Now().DayNumber())
	fmt.Println(carbon.Now().MonthNumber())
	fmt.Println(carbon.Now().YearNumber())
}
