package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
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

	"github.com/essentialkaos/ek/v13/env"
	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/selfupdate"
	"github.com/essentialkaos/ek/v13/selfupdate/interactive"
	storage "github.com/essentialkaos/ek/v13/selfupdate/storage/basic"
	"github.com/essentialkaos/ek/v13/signal"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/apps"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/timeutil"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"

	"github.com/essentialkaos/fz/gofuzz"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "fz"
	VER  = "1.2.0"
	DESC = "Tool for formatting go-fuzz output"
)

// Constants with options names
const (
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_UPDATE       = "U:update"
	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap is map with options
var optMap = options.Map{
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_UPDATE:       {Type: options.MIXED},
	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

var prev gofuzz.Line
var startTime time.Time
var corpusTime time.Time

var colorTagApp, colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(3)

	preConfigureUI()

	_, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error("- "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithApps(apps.Golang()).
			WithChecks(checkForGoFuzz()).
			Print()
		os.Exit(0)
	case options.GetB(OPT_UPDATE):
		os.Exit(updateBinary())
	case options.GetB(OPT_HELP) || !hasStdinData():
		genUsage().Print()
		os.Exit(0)
	}

	configureSignalHandlers()
	processInput()
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{&}{#00ADD8}", "{#5DC9E2}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{&}{#38}", "{#74}"
	default:
		colorTagApp, colorTagVer = "{*}{&}{c}", "{c}"
	}
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
			terminal.Error(err)
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
		timeutil.Pretty(time.Since(startTime)).Short(true),
		cur.Workers, fmtutil.PrettyNum(cur.Corpus),
		timeutil.Pretty(time.Since(corpusTime)).Short(true),
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
	duration := timeutil.Pretty(time.Since(startTime)).Short()
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

// ////////////////////////////////////////////////////////////////////////////////// //

// checkForGoFuzz checks if go-fuzz binary present on the system
func checkForGoFuzz() support.Check {
	goFuzzBin := env.Which("go-fuzz")

	if goFuzzBin == "" {
		return support.Check{support.CHECK_ERROR, "go-fuzz", "Binary not found in PATH"}
	}

	return support.Check{support.CHECK_OK, "go-fuzz", fmt.Sprintf("Binary found (%s)", goFuzzBin)}
}

// updateBinary updates current binary to the latest version
func updateBinary() int {
	quiet := strings.ToLower(options.GetS(OPT_UPDATE)) == "quiet"
	updInfo, hasUpdate, err := storage.NewStorage("https://apps.kaos.ws").Check(APP, VER)

	if err != nil {
		if !quiet {
			terminal.Error("Can't update binary: %v", err)
		}

		return 1
	}

	if !hasUpdate {
		fmtc.If(!quiet).Println("{g}You are using the latest version of the app{!}")
		return 0
	}

	pubKey := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEnYHsOTvrKqeE97dsEt7Ge97+yUcvQJn1++s++FqShDyqwV8CcoKp0E6nDTc8SxInZ5wxwcScxSicfvC9S73OSg=="

	if quiet {
		err = selfupdate.Run(updInfo, pubKey, nil)
	} else {
		err = selfupdate.Run(updInfo, pubKey, interactive.Dispatcher())
	}

	if err != nil {
		return 1
	}

	return 0
}

// printCompletion prints completion for given shell
func printCompletion() int {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(genUsage(), APP))
	case "fish":
		fmt.Print(fish.Generate(genUsage(), APP))
	case "zsh":
		fmt.Print(zsh.Generate(genUsage(), optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("go-fuzz … |& fz")

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddRawExample(
		"go-fuzz -bin app-fuzz.zip |& fz",
		"Run fuzz test for app-fuzz.zip",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{"essentialkaos/fz", update.GitHubChecker}
	}

	return about
}
