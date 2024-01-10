package gofuzz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Line contains data from go-fuzz output line
type Line struct {
	DateTime    time.Time
	Workers     int
	Corpus      int
	Crashers    int
	Restarts    int
	Execs       int
	ExecsPerSec int
	Cover       int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Parse parses go-fuzz output line
func Parse(line string) (Line, error) {
	result := Line{}
	data := strings.Split(line, ",")

	if len(data) < 7 {
		return result, fmt.Errorf("Output line must contain at least 7 fields")
	}

	var err error

	data[3] = strutil.Exclude(data[3], "1/")
	data[4] = strutil.Exclude(data[4], "/sec")
	data[4] = strutil.Exclude(data[4], "(")
	data[4] = strutil.Exclude(data[4], ")")

	dateTime := strutil.ReadField(data[0], 0, false, " ") + " "
	dateTime += strutil.ReadField(data[0], 1, false, " ")

	result.DateTime, err = time.Parse("2006/01/02 15:04:05", dateTime)

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse date and time field: %w", err)
	}

	result.Workers, err = strconv.Atoi(strutil.ReadField(data[0], 3, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse workers field: %w", err)
	}

	result.Corpus, err = strconv.Atoi(strutil.ReadField(data[1], 1, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse corpus field: %w", err)
	}

	result.Crashers, err = strconv.Atoi(strutil.ReadField(data[2], 1, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse crashers field: %w", err)
	}

	result.Restarts, err = strconv.Atoi(strutil.ReadField(data[3], 1, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse restarts field: %w", err)
	}

	result.Execs, err = strconv.Atoi(strutil.ReadField(data[4], 1, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse execs field: %w", err)
	}

	result.ExecsPerSec, err = strconv.Atoi(strutil.ReadField(data[4], 2, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse execs per sec field: %w", err)
	}

	result.Cover, err = strconv.Atoi(strutil.ReadField(data[5], 1, false, " "))

	if err != nil {
		return Line{}, fmt.Errorf("Can't parse cover field: %w", err)
	}

	return result, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //
