package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/hash"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/depsy"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// ShowSupportInfo prints verbose info about application, system, dependencies and
// important environment
func ShowSupportInfo(app, ver, gitRev string, gomod []byte) {
	fmtutil.SeparatorTitleColorTag = "{s-}"
	fmtutil.SeparatorFullscreen = false
	fmtutil.SeparatorColorTag = "{s-}"
	fmtutil.SeparatorSize = 80

	showApplicationInfo(app, ver, gitRev)
	showOSInfo()
	showEnvInfo()
	showDepsInfo(gomod)

	fmtutil.Separator(false)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showApplicationInfo shows verbose information about application
func showApplicationInfo(app, ver, gitRev string) {
	fmtutil.Separator(false, "APPLICATION INFO")

	printInfo(7, "Name", app)
	printInfo(7, "Version", ver)

	printInfo(7, "Go", fmtc.Sprintf(
		"%s {s}(%s/%s){!}",
		strings.TrimLeft(runtime.Version(), "go"),
		runtime.GOOS, runtime.GOARCH,
	))

	if gitRev != "" {
		if !fmtc.DisableColors && fmtc.IsTrueColorSupported() {
			printInfo(7, "Git SHA", gitRev+getHashColorBullet(gitRev))
		} else {
			printInfo(7, "Git SHA", gitRev)
		}
	}

	bin, _ := os.Executable()
	binSHA := hash.FileHash(bin)

	if binSHA != "" {
		binSHA = strutil.Head(binSHA, 7)
		if !fmtc.DisableColors && fmtc.IsTrueColorSupported() {
			printInfo(7, "Bin SHA", binSHA+getHashColorBullet(binSHA))
		} else {
			printInfo(7, "Bin SHA", binSHA)
		}
	}
}

// showEnvInfo shows info about environment
func showEnvInfo() {
	fmtutil.Separator(false, "ENVIRONMENT")

	printInfo(7, "Golang", getGOVersion())
	printInfo(7, "go-fuzz", getGOFuzzVersion())
}

// showDepsInfo shows information about all dependencies
func showDepsInfo(gomod []byte) {
	deps := depsy.Extract(gomod, false)

	if len(deps) == 0 {
		return
	}

	fmtutil.Separator(false, "DEPENDENCIES")

	for _, dep := range deps {
		if dep.Extra == "" {
			fmtc.Printf(" {s}%8s{!}  %s\n", dep.Version, dep.Path)
		} else {
			fmtc.Printf(" {s}%8s{!}  %s {s-}(%s){!}\n", dep.Version, dep.Path, dep.Extra)
		}
	}
}

// getHashColorBullet return bullet with color from hash
func getHashColorBullet(v string) string {
	if len(v) > 6 {
		v = strutil.Head(v, 6)
	}

	return fmtc.Sprintf(" {#" + strutil.Head(v, 6) + "}● {!}")
}

// printInfo formats and prints info record
func printInfo(size int, name, value string) {
	name = name + ":"
	size++

	if value == "" {
		fm := fmt.Sprintf("  {*}%%-%ds{!}  {s-}—{!}\n", size)
		fmtc.Printf(fm, name)
	} else {
		fm := fmt.Sprintf("  {*}%%-%ds{!}  %%s\n", size)
		fmtc.Printf(fm, name, value)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getGOVersion returns version of installed Go
func getGOVersion() string {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	goVer := strings.Trim(string(out), "\n\t\r")
	goVer = strutil.ReadField(goVer, 2, false, " ")
	goVer = strutil.Exclude(goVer, "go")

	return goVer
}

// getGOFuzzVersion returns go-fuzz module version
func getGOFuzzVersion() string {
	cmd := exec.Command("go", "list", "-m", "github.com/dvyukov/go-fuzz")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	fuzzVer := strings.Trim(string(out), "\n\t\r")
	fuzzVer = strutil.ReadField(fuzzVer, 1, false, " ")
	fuzzVer = strings.TrimLeft(fuzzVer, "v")

	return fuzzVer
}
