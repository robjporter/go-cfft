package main

import (
	"fmt"
	"time"

	"./xTools/as"
	"./xTools/carbon"
	"./xTools/permissions"
	"./xTools/xstrings"
)

func main() {
	xstringsFunc()
	permissionsFunc()
	carbonFunc()
	asFunc()
}

func asFunc() {
	fmt.Println("AS FUNC======================================================")

}

func xstringsFunc() {
	fmt.Println(xstrings.Center("Center Me", "=", 50))
	fmt.Println(xstrings.Center("Center Me", "<>", 50))
}

func permissionsFunc() {
	anAcl := permissions.Acl{}
	anAcl.Grant("Admin", "Read", "Write", "Execute", "Water the flowers")
	fmt.Println(anAcl.Can("Admin", "Execute"))
	fmt.Println(anAcl.Can("Admin", "Water the flowers"))
	fmt.Println(anAcl.Can("Admin", "Test"))
	anAcl.Revoke("Admin", "Water the flowers")
	fmt.Println(anAcl.Can("Admin", "Water the flowers"))
}

func carbonFunc() {
	fmt.Println("ORDINAL: ")
	fmt.Println(carbon.Now().Ordinal())
	fmt.Println(carbon.Now().OrdinalOnly())

	fmt.Println("\nMONTH NAME: ")
	fmt.Println(carbon.Now().MonthName())
	fmt.Println("\nPREVIOUS MONTH NAME: ")
	fmt.Println(carbon.Now().PreviousMonthStartDay().MonthName())

	fmt.Println("\nDAY NUMBER: ")
	fmt.Println(carbon.Now().DayNumber())
	fmt.Println("\nMONTH NUMBER: ")
	fmt.Println(carbon.Now().MonthNumber())
	fmt.Println("\nYEAR NUMBER: ")
	fmt.Println(carbon.Now().YearNumber())

	fmt.Println("\nTOMORROW: ")
	fmt.Println(carbon.Now().Tomorrow().DayNumber())
	fmt.Println(carbon.Now().Tomorrow())
	fmt.Println(carbon.Now().Tomorrow().StartOfDay())
	fmt.Println(carbon.Now().Tomorrow().EndOfDay())

	fmt.Println("\nYESTERDAY: ")
	fmt.Println(carbon.Now().Yesterday().DayNumber())
	fmt.Println(carbon.Now().Yesterday())
	fmt.Println(carbon.Now().Yesterday().StartOfDay())
	fmt.Println(carbon.Now().Yesterday().EndOfDay())

	fmt.Println("\nCALENDAR YEAR QUARTER: ")
	fmt.Println(carbon.Now().Quarter())
	fmt.Println("\nFINANCIAL YEAR QUARTER STARTING IN AUGUST: ")
	fmt.Println(carbon.Now().Quarter(time.August))

	fmt.Println("\nIS LEAP YEAR: ")
	fmt.Println(carbon.Now().IsLeapYear())
	fmt.Println()

	fmt.Printf("CURRENT MONTH: %s DAYS: %d\n", carbon.Now().MonthName(), carbon.Now().DaysInMonth())
	fmt.Printf("PREVIOUS MONTH: %s DAYS: %d\n", carbon.Now().PreviousMonth().MonthName(), carbon.Now().PreviousMonth().DaysInMonth())

	fmt.Println("\nDAYS LEFT IN WEEK: ")
	fmt.Println(carbon.Now().DaysLeftInWeek())

	fmt.Println("\nDAYS LEFT IN MONTH: ")
	fmt.Println(carbon.Now().DaysLeftInMonth())

	fmt.Println("\nDAYS LEFT IN YEAR: ")
	fmt.Println(carbon.Now().DaysLeftInYear())

	fmt.Println("\n30 DAYS TO HOURS: ")
	fmt.Println(carbon.Now().DaysToHours(30))

	fmt.Println("\nIS WEEKEND: ")
	fmt.Println(carbon.Now().IsWeekend())

	fmt.Println("\nIS WEEKDAY: ")
	fmt.Println(carbon.Now().IsWeekday())
}
