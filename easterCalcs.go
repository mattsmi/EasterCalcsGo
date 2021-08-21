/*
  Calculates the date of Easter according to the Julian, Revised, Jullian and Gregorian calendars.
  Usage given by executing the program with an argument of "-h".
*/
package main

import (
  "flag"
  "fmt"
  "os"
  "time"
  )

const (
  iEDM_JULIAN  = 1
  iEDM_ORTHODOX  = 2
  iEDM_WESTERN  = 3
  iFIRST_EASTER_YEAR  = 326
  iFIRST_VALID_GREGORIAN_YEAR  = 1583
  iLAST_VALID_GREGORIAN_YEAR  = 4099
  layoutISO = "2006-01-02"
)

func findInt(slice []int, val int) (int, bool) {
  /* Find takes a slice and looks for an element in it. If found it will
     return its key, otherwise it will return -1 and a bool of false.
  */
  for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

func pF15_CalcDateOfEaster(iYearToFind int, iDatingMethod int) (time.Time) {

  // default values for invalid arguments
  var imDay int = 0
  var imMonth int = 0
  // intermediate results (all integers)
  var iFirstDig int = 0
  var iRemain19 int = 0
  var iTempNum int = 0
  // tables A to E results (all integers)
  var iTableA int = 0
  var iTableB int = 0
  var iTableC int = 0
  var iTableD int = 0
  var iTableE int = 0

  //  Calculate Easter Sunday date
  // first 2 digits of year (integer division)
  iFirstDig = iYearToFind / 100
  // remainder of year / 19
  iRemain19 = iYearToFind % 19


  if (iDatingMethod == iEDM_JULIAN) || (iDatingMethod == iEDM_ORTHODOX) {
    //  calculate PFM date
    iTableA = ((225 - 11 * iRemain19) % 30) + 21

    //  find the next Sunday
    iTableB = (iTableA - 19) % 7
    iTableC = (40 - iFirstDig) % 7

    iTempNum = iYearToFind % 100
    iTableD = (iTempNum + (iTempNum / 4)) % 7

    iTableE = ((20 - iTableB - iTableC - iTableD) % 7) + 1
    imDay = iTableA + iTableE

    // convert Julian to Gregorian date
    if (iDatingMethod == iEDM_ORTHODOX) {
       // 10 days were # skipped#  in the Gregorian calendar from 5-14 Oct 1582
        iTempNum  = 10
        // Only 1 in every 4 century years are leap years in the Gregorian
        // calendar (every century is a leap year in the Julian calendar)
        if (iYearToFind > 1600) {
            iTempNum = iTempNum + iFirstDig - 16 - ((iFirstDig - 16) / 4)
        }
        imDay = imDay + iTempNum
    }
  } else {
    // That is iDatingMethod == iEDM_WESTERN
    //  calculate PFM date
    iTempNum = ((iFirstDig - 15) / 2) + 202 - 11 * iRemain19
    lFirstList := []int{21, 24, 25, 27, 28, 29, 30, 31, 32, 34, 35, 38}
    lSecondList := []int{33, 36, 37, 39, 40}
    _, foundInFirst := findInt(lFirstList, iFirstDig)
    _, foundInSecond := findInt(lSecondList, iFirstDig)
    if foundInFirst {
        iTempNum = iTempNum - 1
    } else if foundInSecond {
        iTempNum = iTempNum - 2
    }
    iTempNum = iTempNum % 30

    iTableA  = iTempNum + 21
    if iTempNum == 29 {
        iTableA = iTableA - 1
    }
    if ((iTempNum == 28) && (iRemain19 > 10)) {
        iTableA = iTableA - 1
    }

    //  find the next Sunday
    iTableB = (iTableA - 19) % 7

    iTableC = (40 - iFirstDig) % 4
    if iTableC == 3 {
        iTableC = iTableC + 1
    }
    if iTableC > 1 {
        iTableC = iTableC + 1
    }

    iTempNum = iYearToFind % 100
    iTableD = (iTempNum + iTempNum / 4) % 7

    iTableE = ((20 - iTableB - iTableC - iTableD) % 7) + 1
    imDay = iTableA + iTableE
  }


  //  return the date
  if imDay > 61 {
      imDay = imDay - 61
      imMonth = 5
      // for imMethod 2, Easter Sunday can occur in May
  } else if imDay > 31 {
      imDay = imDay - 31
      imMonth = 4
  } else {
      imMonth = 3
  }

  // Format the date and return it to the calling function.
  return time.Date(iYearToFind, time.Month(imMonth), imDay, 0, 0, 0, 0, time.UTC)

}

func main() {

  //declare variables
  var yearSought int
  var calendarToUse int
  thisYear, _, _ := time.Now().Date()

  //flags declaration

  flag.IntVar(&yearSought, "y", thisYear, "Specify the year, for which the date of Easter is sought.")
  flag.IntVar(&calendarToUse, "c", 3, "Which calendar to use (1: Julian; 2: Revised Julian; 3: Gregorian)?")
  flag.Parse()

  // check command-line arguments
  if (yearSought < iFIRST_EASTER_YEAR) || (yearSought > iLAST_VALID_GREGORIAN_YEAR){
    fmt.Printf("The year must be between AD %d and AD %d.\n", iFIRST_EASTER_YEAR, iLAST_VALID_GREGORIAN_YEAR)
    os.Exit(1)
  }
  if (calendarToUse != iEDM_JULIAN && calendarToUse != iEDM_ORTHODOX && calendarToUse != iEDM_WESTERN) {
    fmt.Printf("The calendar must of one of 1, 2, or 3: (1: Julian; 2: Revised Julian; 3: Gregorian).\n")
    os.Exit(1)
  }

  var easterDate time.Time
  easterDate = pF15_CalcDateOfEaster(yearSought, calendarToUse)

  fmt.Printf(easterDate.Format(layoutISO) + "\n")
}
