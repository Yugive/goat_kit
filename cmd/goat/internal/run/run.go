package run

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run project",
	Long:  "Run project. Example: goat run",
	Run:   Run,
}

//Run the application
func Run(cmd *cobra.Command, args []string) {
	var dir string

	if len(args) > 0 {
		dir = args[0]
	}

	base, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if dir == "" {
		cmdPath, err := findCMD(base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		if len(cmdPath) == 0 {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", "The cmd directory cannot be found in the current directory")
			return
		} else if len(cmdPath) == 1 {
			for k, v := range cmdPath {
				dir = path.Join(k, v)
			}
		} else {
			var cmdPaths []string
			for k := range cmdPath {
				cmdPaths = append(cmdPaths, k)
			}
			prompt := &survey.Select{
				Message: "Which directory do you want to run?",
				Options: cmdPaths,
			}
			e := survey.AskOne(prompt, &dir)
			if e != nil || dir == "" {
				return
			}
			dir = path.Join(cmdPath[dir], dir)
		}
	}

	fd := exec.Command("go", "run", ".")
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = dir
	if err := fd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
	}
}

func findCMD(base string) (map[string]string, error) {
	var root bool
	next := func(dir string) (map[string]string, error) {
		cmdPath := make(map[string]string)
		err := filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
			if strings.HasSuffix(walkPath, "cmd") {
				paths, err := os.ReadDir(walkPath)
				if err != nil {
					return err
				}
				for _, fileInfo := range paths {
					if fileInfo.IsDir() {
						cmdPath[path.Join("cmd", fileInfo.Name())] = filepath.Join(walkPath, "..")
					}
				}
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return cmdPath, err
	}
	for i := 0; i < 5; i++ {
		tmp := base
		cmd, err := next(tmp)
		if err != nil {
			return nil, err
		}
		if len(cmd) > 0 {
			return cmd, nil
		}
		if root {
			break
		}
		_ = filepath.Join(base, "..")
	}

	return map[string]string{"": base}, nil
}
