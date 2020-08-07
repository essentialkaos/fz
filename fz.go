package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/signal"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/timeutil"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "fz"
	VER  = "0.0.3"
	DESC = "Tool for formatting go-fuzz output"
)

// Constants with options names
const (
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Info contains info about test
type Info struct {
	DateTime    string
	Workers     int
	Corpus      int
	CorpusDur   int64
	Crashers    int
	Restarts    string
	Execs       int
	ExecsPerSec int
	Cover       int
	Uptime      int64
}

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap is map with options
var optMap = options.Map{
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// ////////////////////////////////////////////////////////////////////////////////// //

var startInfo Info
var prevInfo Info

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main func
func main() {
	runtime.GOMAXPROCS(3)

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError("Options parsing errors:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		os.Exit(0)
	}

	if options.GetB(OPT_HELP) || !hasStdinData() {
		showUsage()
		os.Exit(0)
	}

	signal.Handlers{
		signal.INT:  signalHandler,
		signal.TERM: signalHandler,
		signal.QUIT: signalHandler,
	}.TrackAsync()

	processInput()
}

// processInput processes go-fuzz output passed to this tool
func processInput() {
	r := bufio.NewReader(os.Stdin)
	s := bufio.NewScanner(r)

	fmtc.TPrintf("{s-}Starting tests…{!}")

	for s.Scan() {
		data := s.Text()

		if isShutdownMessage(data) {
			time.Sleep(time.Minute)
		}

		info, ok := parseInfoLine(data)

		if !ok {
			fmtc.TPrintf("")
			printError(data)
			os.Exit(1)
		}

		renderInfo(info)

		prevInfo = info

		if startInfo.Uptime == 0 {
			startInfo = info
		}

		if startInfo.Cover == 0 && info.Cover > 0 {
			startInfo.Cover = info.Cover
		}
	}
}

// parseInfoLine parses line with go-fuzz output data
func parseInfoLine(data string) (Info, bool) {
	info := Info{}
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) < 7 {
		return info, false
	}

	corpusDate := strings.Trim(strutil.ReadField(dataSlice[1], 2, false, " "), "(")
	uptime := strutil.ReadField(dataSlice[6], 1, false, " ")

	info.DateTime = strutil.ReadField(dataSlice[0], 0, false, " ")
	info.DateTime += " " + strutil.ReadField(dataSlice[0], 1, false, " ")
	info.Workers, _ = strconv.Atoi(strutil.ReadField(dataSlice[0], 3, false, " "))
	info.Corpus, _ = strconv.Atoi(strutil.ReadField(dataSlice[1], 1, false, " "))
	info.CorpusDur = timeutil.ParseDuration(corpusDate)
	info.Crashers, _ = strconv.Atoi(strutil.ReadField(dataSlice[2], 1, false, " "))
	info.Restarts = strutil.ReadField(dataSlice[3], 1, false, " ")
	info.Execs, _ = strconv.Atoi(strutil.ReadField(dataSlice[4], 1, false, " "))
	execsPerSec := strings.Trim(strutil.ReadField(dataSlice[4], 2, false, " "), "(/sec)")
	info.ExecsPerSec, _ = strconv.Atoi(execsPerSec)
	info.Cover, _ = strconv.Atoi(strutil.ReadField(dataSlice[5], 1, false, " "))
	info.Uptime = timeutil.ParseDuration(uptime)

	return info, true
}

// renderInfo render line with info
func renderInfo(cur Info) {
	var crashersTag string

	workersTag := getIndicatorTag(cur.Workers, prevInfo.Workers)
	corpusTag := getIndicatorTag(cur.Corpus, prevInfo.Corpus)
	coverTag := getIndicatorTag(cur.Cover, prevInfo.Cover)

	if cur.Crashers != 0 {
		crashersTag = "{r}"
	}

	fmtc.TPrintf(
		"{s}%s{!} {s-}[%s]{!} {*}Workers:{!} "+workersTag+"%d{!} {s-}•{!} {*}Corpus:{!} "+corpusTag+"%s{!} {s-}(%s){!} {s-}•{!} {*}Crashers:{!} "+crashersTag+"%d {s-}•{!} {*}Restarts:{!} %s {s-}•{!} {*}Cover:{!} "+coverTag+"%s{!} {s-}•{!} {*}Execs:{!} %s{s}/s{!} {s-}(%s){!}",
		cur.DateTime, formatDuration(cur.Uptime), cur.Workers, fmtutil.PrettyNum(cur.Corpus),
		formatDuration(cur.CorpusDur), cur.Crashers, cur.Restarts, fmtutil.PrettyNum(cur.Cover),
		fmtutil.PrettyNum(cur.ExecsPerSec), fmtutil.PrettyNum(cur.Execs),
	)
}

// printResults prints tests results
func printResults() {
	if startInfo.Uptime == 0 {
		startInfo = prevInfo
	}

	corpus := formatResultNum(prevInfo.Corpus - startInfo.Corpus)
	crashers := formatResultNum(prevInfo.Crashers - startInfo.Crashers)
	cover := formatResultNum(prevInfo.Cover - startInfo.Cover)
	duration := timeutil.PrettyDuration(prevInfo.Uptime)
	execs := fmtutil.PrettyNum(prevInfo.Execs)

	fmtc.TPrintf(
		"{*}Duration:{!} %s {s-}•{!} {*}Execs:{!} %s {s-}•{!} {*}Corpus:{!} %s {s-}•{!} {*}Crashers:{!} %s {s-}•{!} {*}Cover:{!} %s\n",
		duration, execs, corpus, crashers, cover,
	)
}

// getIndicatorTag returns color tag based on difference between
// current and previous values
func getIndicatorTag(v1, v2 int) string {
	switch {
	case v1 > v2:
		return "{g}"
	case v1 < v2:
		return "{r}"
	default:
		return ""
	}
}

// formatDuration formats duration
func formatDuration(d int64) string {
	var min, sec int64

	min = d / 60
	sec = d % 60

	return fmtc.Sprintf("%01d:%02d", min, sec)
}

// formatResultNum formats number for results
func formatResultNum(v int) string {
	if v <= 0 {
		return "0"
	}

	return fmtc.Sprintf("+%d", v)
}

// signalHandler is signal handler
func signalHandler() {
	printResults()
	os.Exit(0)
}

// hasStdinData return true if there is some data in stdin
func hasStdinData() bool {
	stdin, err := os.Stdin.Stat()

	if err != nil {
		return false
	}

	if stdin.Mode()&os.ModeCharDevice != 0 {
		return false
	}

	return true
}

// isShutdownMessage returns true if data contains shutdown message
func isShutdownMessage(data string) bool {
	return strings.Contains(data, "shutting down...")
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage print usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("go-fuzz … |& fz")

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(genUsage(), "fz"))
	case "fish":
		fmt.Printf(fish.Generate(genUsage(), "fz"))
	case "zsh":
		fmt.Printf(zsh.Generate(genUsage(), optMap, "fz"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout print info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "Essential Kaos",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/fz", update.GitHubChecker},
	}

	about.Render()
}
