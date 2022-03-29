package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/signal"
	"github.com/essentialkaos/ek/v12/timeutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/fz/gofuzz"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "fz"
	VER  = "1.0.1"
	DESC = "Tool for formatting go-fuzz output"
)

// Constants with options names
const (
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap is map with options
var optMap = options.Map{
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

var prev gofuzz.Line
var startTime time.Time
var corpusTime time.Time

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main func
func main() {
	runtime.GOMAXPROCS(3)

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		os.Exit(genCompletion())
	}

	if options.Has(OPT_GENERATE_MAN) {
		os.Exit(genMan())
	}

	configureUI()

	if options.GetB(OPT_VER) {
		os.Exit(showAbout())
	}

	if options.GetB(OPT_HELP) || !hasStdinData() {
		os.Exit(showUsage())
	}

	configureSignalHandlers()
	processInput()
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// configureSignalHandlers configures signal handlers
func configureSignalHandlers() {
	signal.Handlers{
		signal.INT:  signalHandler,
		signal.TERM: signalHandler,
		signal.QUIT: signalHandler,
	}.TrackAsync()
}

// processInput processes go-fuzz output passed to this tool
func processInput() {
	r := bufio.NewReader(os.Stdin)
	s := bufio.NewScanner(r)

	startTime = time.Now()
	corpusTime = time.Now()

	fmtc.TPrintf("{s-}Starting tests…{!}")

	for s.Scan() {
		data := s.Text()

		if isShutdownMessage(data) {
			time.Sleep(time.Minute)
		}

		line, err := gofuzz.Parse(data)

		if err != nil {
			fmtc.TPrintf("")
			printError(err.Error())
			os.Exit(1)
		}

		renderInfo(line)

		prev = line
	}
}

// renderInfo render line with info
func renderInfo(cur gofuzz.Line) {
	var crashersTag string

	workersTag := getIndicatorTag(cur.Workers, prev.Workers)
	corpusTag := getIndicatorTag(cur.Corpus, prev.Corpus)
	coverTag := getIndicatorTag(cur.Cover, prev.Cover)

	if cur.Crashers != 0 {
		crashersTag = "{r}"
	}

	if cur.Corpus != prev.Corpus {
		corpusTime = time.Now()
	}

	fmtc.TPrintf(
		"{s}%s{!} {s-}[%s]{!} {*}Workers:{!} "+workersTag+"%d{!} {s-}•{!} {*}Corpus:{!} "+corpusTag+"%s{!} {s-}(%s){!} {s-}•{!} {*}Crashers:{!} "+crashersTag+"%s {s-}•{!} {*}Restarts:{!} {s}1/{!}%s {s-}•{!} {*}Cover:{!} "+coverTag+"%s{!} {s-}•{!} {*}Execs:{!} %s{s}/s{!} {s-}(%s){!}",
		timeutil.Format(cur.DateTime, "%Y/%m/%d %H:%M:%S"),
		timeutil.ShortDuration(time.Since(startTime), false),
		cur.Workers, fmtutil.PrettyNum(cur.Corpus),
		timeutil.ShortDuration(time.Since(corpusTime), false),
		fmtutil.PrettyNum(cur.Crashers),
		fmtutil.PrettyNum(cur.Restarts),
		fmtutil.PrettyNum(cur.Cover),
		fmtutil.PrettyNum(cur.ExecsPerSec),
		fmtutil.PrettyNum(cur.Execs),
	)
}

// printResults prints tests results
func printResults() {
	corpus := fmtutil.PrettyNum(prev.Corpus)
	crashers := fmtutil.PrettyNum(prev.Crashers)
	cover := fmtutil.PrettyNum(prev.Cover)
	duration := timeutil.ShortDuration(time.Since(startTime))
	execs := fmtutil.PrettyNum(prev.Execs)

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

// genCompletion generates completion for different shells
func genCompletion() int {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(genUsage(), APP))
	case "fish":
		fmt.Printf(fish.Generate(genUsage(), APP))
	case "zsh":
		fmt.Printf(zsh.Generate(genUsage(), optMap, APP))
	default:
		return 1
	}

	return 0
}

// genMan generates man page
func genMan() int {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(),
		),
	)

	return 0
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("go-fuzz … |& fz")

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	return info
}

// genAbout generates info about version
func genAbout() *usage.About {
	return &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/fz", update.GitHubChecker},
	}
}

// showUsage prints usage info
func showUsage() int {
	genUsage().Render()
	return 0
}

// showAbout prints info about version
func showAbout() int {
	genAbout().Render()
	return 0
}
