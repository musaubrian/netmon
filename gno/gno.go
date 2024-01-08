package gno

import (
	"bufio"
	"fmt"
	"io"
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
	conf     Config
	commands []command
}

type Config struct {
	buildDir string
	binName  string
	src      string
}

const (
	ERROR = iota
	WARN
	INFO
	CMD
)

const separator = string(os.PathSeparator)

func logMsg(level int, msg string) {
	switch {
	case level == ERROR:
		log.Fatalf("[ERROR] %s\n", msg)
	case level == INFO:
		log.Printf("[INFO]  %s\n", msg)
	case level == WARN:
		log.Printf("[WARN]  %s\n", msg)
	case level == CMD:
		log.Printf("[CMD]   %s\n", msg)
	default:
		log.Fatalf("[ERROR] %s\n", msg)
	}
}

func listFiles(pattern string) string {
	var matches string
	files, err := filepath.Glob(pattern)
	if err != nil {
		logMsg(ERROR, fmt.Sprintf("No matching  pattern files found for the pattern %s", pattern))
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
		cmdOpts += v + " "
	}
	return strings.TrimSpace(cmdOpts)
}

func backToPrevWorkDir(cwd string) {
	if err := os.Chdir(cwd); err != nil {
		logMsg(ERROR, err.Error())
	}
}

func New() *gnoDets {
	return &gnoDets{}
}

// Overrides the default config created by New()
func (g *gnoDets) BootstrapBuild(buildDirLocation string, bin string, source string) {
	g.conf.buildDir = buildDirLocation
	g.conf.binName = bin
	g.conf.src = source
	if len(g.conf.buildDir) == 0 {
		logMsg(ERROR, "Build directory not provided")
	} else {
		err := os.Mkdir(g.conf.buildDir, 0o770)
		if err != nil {
			logMsg(INFO, fmt.Sprintf("`%s` already exists, skipping", g.conf.buildDir))
		} else {
			logMsg(INFO, "Created build directory")
		}
	}
}

// Copy resources to the final build dir
func (g gnoDets) CopyResources(src string) {
	copyResources(src, g.conf.buildDir)
}

// Add commands to be executed
func (g *gnoDets) AddCommand(name string, opts ...string) {
	c := &command{
		name: name,
		opts: opts,
	}
	g.commands = append(g.commands, *c)

}

func printOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logMsg(ERROR, "Error reading output: "+err.Error())
		return
	}
}

// Run commands synchronously
func (g gnoDets) RunCommandsSync() {
	if len(g.commands) >= 1 {
		for _, v := range g.commands {
			ms := fmt.Sprintf("Running command: `%s %s`\n", v.name, formatOpts(v.opts))
			logMsg(CMD, ms)
			cmd := exec.Command(v.name, v.opts...)

			res, err := cmd.StdoutPipe()
			if err != nil {
				logMsg(ERROR, err.Error())
			}

			stdErr, err := cmd.StderrPipe()
			if err != nil {
				logMsg(ERROR, err.Error())
			}
			if err := cmd.Start(); err != nil {
				logMsg(ERROR, "Could not start cmd")
			}

			go printOutput(res)
			go printOutput(stdErr)

			if err := cmd.Wait(); err != nil {
				logMsg(ERROR, "IN WAIT: "+err.Error())
			}

		}
	} else {
		logMsg(INFO, "No commands, skipping")
	}
}

// Builds the binary and runs the commands Synchronously
// So they need to be ordered correctly
func (g gnoDets) Build() {
	buildBinary(g)
}

func buildBinary(g gnoDets) {
	if len(g.conf.buildDir) < 1 {
		logMsg(ERROR, "Build not bootstrapped")
	}
	src := listFiles(g.conf.src)
	if len(src) < 1 {
		src = g.conf.src
	}

	cwd, err := os.Getwd()
	if err != nil {
		logMsg(ERROR, err.Error())
	}

	defer backToPrevWorkDir(cwd)
	p, err := filepath.Abs(".")
	if err != nil {
		logMsg(INFO, p)
		logMsg(ERROR, err.Error())
	}

	binLoc := filepath.Join(p, g.conf.buildDir, g.conf.binName)
	if err := os.Chdir(src); err != nil {
		logMsg(ERROR, err.Error())
	}
	logMsg(INFO, "Switched to "+src)
	cmd := exec.Command("go", "build", "-o", binLoc, g.conf.src)
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd)
		logMsg(CMD, string(res))
		logMsg(ERROR, "Failed to build binary "+err.Error())
	}
	fullBin := strings.Split(binLoc, separator)
	binLoc = filepath.Join(fullBin[len(fullBin)-2], fullBin[len(fullBin)-1])

	logMsg(INFO, fmt.Sprintf("Built Binary -> %s", binLoc))
}

func copyResources(src string, dest string) {
	if dest == src {
		logMsg(ERROR, "Cannot copy a folder into itself!")
	}
	files, err := os.ReadDir(src)
	if err != nil {
		logMsg(WARN, err.Error())
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
		logMsg(ERROR, err.Error())
	}
	err = os.MkdirAll(filepath.Dir(dest), 0770)
	if err != nil {
		logMsg(ERROR, err.Error())
		logMsg(INFO, "Skipping  dir creation")
	}
	err = os.WriteFile(dest, content, 0644)
	if err != nil {
		logMsg(ERROR, err.Error())
	}
	logMsg(INFO, fmt.Sprintf("Copied %s -> %s", src, dest))
}
