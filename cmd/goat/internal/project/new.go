package project

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"goat/internal/base"
	"os"
	"path"
)

type Project struct {
	Name string
	Path string
}

func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("%s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "Do you want to override the folder?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}

	fmt.Sprintf("Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)

	repo := base.NewRepo(layout, branch)
	fmt.Println(layout, branch)
	if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}

	e := os.Rename(path.Join(to, "cmd", "server"), path.Join(to, "cmd", p.Name))
	if e != nil {
		return e
	}
	base.Tree(to, dir)

	fmt.Printf("Project creation succeeded: %s\n", color.GreenString(p.Name))
	fmt.Print("Use the following command to start the project: \n\n")
	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go generate ./..."))
	fmt.Println(color.WhiteString("$ go build -o ./bin ./..."))
	fmt.Printf(color.WhiteString("$ ./bin/%s -conf ./configs\n", p.Name))
	fmt.Println("Thank you for using goat-kit")
	return nil
}
