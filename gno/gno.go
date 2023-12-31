package gno

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type command struct {
	name string
	opts []string
}

type gnoDets struct {
	buildDir string
	binName  string
	src      string
	commands []command
}

const separator = string(os.PathSeparator)

func logMsg(msg string, level string) {
	switch {
	case level == "error":
		log.Fatalf("[ERROR] %s\n", msg)
	case level == "info":
		log.Printf("[INFO]  %s\n", msg)
	case level == "warn":
		log.Printf("[WARN]  %s\n", msg)
	case level == "cmd":
		log.Printf("[CMD]   %s\n", msg)
	default:
		log.Fatalf("[ERROR] %s\n", msg)
	}
}

func listFiles(pattern string) string {
	var matches string
	files, err := filepath.Glob(pattern)
	if err != nil {
		logMsg(fmt.Sprintf("No matching files found for the pattern %s", pattern), "error")
	}
	for _, v := range files {
		if !strings.Contains(v, "_test") {
			matches = matches + v + string(" ")
		}
	}
	return strings.TrimSpace(matches)
}

func formatOpts(opts []string) string {
	var cmdOpts string
	for _, v := range opts {
		cmdOpts += v
	}
	return cmdOpts
}

func backToPrevWorkDir(cwd string) {
	if err := os.Chdir(cwd); err != nil {
		logMsg(err.Error(), "error")
	}
}

func New() *gnoDets {
	return &gnoDets{}
}

// Sets up the build location
// Provide the location to put the build artefacts if any
func (g *gnoDets) BootstrapBuild(buildDirLocation string, bin string, source string) {
	g.buildDir = buildDirLocation
	g.binName = bin
	g.src = source
	if len(g.buildDir) == 0 {
		logMsg("Build directory not provided", "error")
	} else {
		err := os.Mkdir(g.buildDir, 0o770)
		if err != nil {
			logMsg(err.Error(), "warn")
			logMsg("Skipping build dir creation", "info")
		} else {
			logMsg("Created build directory", "info")
		}
	}
}

// Copy resources to the final build dir
func (g gnoDets) CopyResources(src string) {
	copyResources(src, g.buildDir)
}

// Add commands to be executed
// These are run synchronously, so they need to be ordered correctly
//
// TO run in a non blocking way, run them in the background with `&`
// example:
//
// g.AddCommand("templ", "--generate &")
//
// or
//
// g.AddCommand("templ", "--generate", "&")
func (g *gnoDets) AddCommand(name string, opts ...string) {
	c := &command{
		name: name,
		opts: opts,
	}
	g.commands = append(g.commands, *c)

}

func spaceCmdOpts(opts ...string) []string {
	var spacedOut []string
	var c string
	for _, v := range opts {
		c = " " + v
		spacedOut = append(spacedOut, c)
	}
	return spacedOut
}

func runCommands(g gnoDets) {
	if len(g.commands) >= 1 {
		for _, v := range g.commands {
			opts := spaceCmdOpts(v.opts...)
			ms := fmt.Sprintf("Running command: `%s %s`", v.name, formatOpts(v.opts))
			logMsg(ms, "cmd")
			res, err := exec.Command(v.name, opts...).CombinedOutput()
			if err != nil {
				logMsg(string(res), "warn")
				logMsg(err.Error(), "error")
			} else {
				fmt.Println(string(res))
			}
		}
	} else {
		logMsg("No commands, skipping", "info")
	}
}

// Builds the binary and runs the commands Synchronously
func (g gnoDets) Build() {
	buildBinary(g)
	runCommands(g)
}

func buildBinary(g gnoDets) {
	src := listFiles(g.src)
	if len(src) < 1 {
		src = g.src
	}

	cwd, err := os.Getwd()
	if err != nil {
		logMsg(err.Error(), "error")
	}

	defer backToPrevWorkDir(cwd)
	p, err := filepath.Abs(".")
	if err != nil {
		logMsg(p, "info")
		logMsg(err.Error(), "error")
	}

	binLoc := filepath.Join(p, g.buildDir, g.binName)
	if err := os.Chdir(src); err != nil {
		logMsg(err.Error(), "error")
	}
	logMsg("Switched to "+src, "info")
	cmd := exec.Command("go", "build", "-o", binLoc, ".")
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd)
		logMsg(string(res), "cmd")
		logMsg("Failed to build binary "+err.Error(), "error")
	}
	fullBin := strings.Split(binLoc, separator)
	binLoc = filepath.Join(fullBin[len(fullBin)-2], fullBin[len(fullBin)-1])

	logMsg(fmt.Sprintf("Built Binary -> %s", binLoc), "info")
}

func copyResources(src string, dest string) {
	if dest == src {
		logMsg("Cannot copy a folder into itself!", "error")
	}
	files, err := os.ReadDir(src)
	if err != nil {
		logMsg(err.Error(), "warn")
		copyFile(src, filepath.Join(dest, src))
	}
	for _, f := range files {
		var destPath string
		srcPath := filepath.Join(src, f.Name())
		if !strings.Contains(dest, strings.Split(srcPath, "/")[0]) {
			destPath = filepath.Join(dest, strings.Split(srcPath, "/")[0], f.Name())
		} else {
			destPath = filepath.Join(dest, f.Name())
		}
		if f.IsDir() {
			copyResources(srcPath, destPath)
		} else {
			copyFile(srcPath, destPath)
		}
	}
}

func copyFile(src string, dest string) {

	content, err := os.ReadFile(src)
	if err != nil {
		logMsg(err.Error(), "error")
	}
	err = os.MkdirAll(filepath.Dir(dest), 0770)
	if err != nil {
		logMsg(err.Error(), "error")
		logMsg("Skipping dir creation", "info")
	}
	err = os.WriteFile(dest, content, 0644)
	if err != nil {
		logMsg(err.Error(), "error")
	}
	logMsg(fmt.Sprintf("Copied %s -> %s", src, dest), "info")
}
