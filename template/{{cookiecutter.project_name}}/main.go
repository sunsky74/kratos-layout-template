package main

import (
	cmd "{{cookiecutter.project_name}}/cmd/{{cookiecutter.project_name}}"
)

func main() {
	app, f := cmd.NewApp()
	defer f()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
