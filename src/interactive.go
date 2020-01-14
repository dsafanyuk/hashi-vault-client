package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/mitchellh/cli"
)

type InteractiveCommand struct {
	Ui cli.Ui
}

func (c *InteractiveCommand) Run(args []string) int {

	prompt := promptui.Select{
		Label: "Select which service will be using this secret",
		Items: []string{"Builder API", "Ansible"},
	}

	_, service, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	prompt = promptui.Select{
		Label: "Which service is this secret for?",
		Items: []string{"AWS", "AWX", "ITRC"},
	}

	_, secretService, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	prompt = promptui.Select{
		Label: "What's the secret's name?",
		Items: []string{"api_key_id", "api_secret_key", "username"},
	}
	_, secretName, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	prompt = promptui.Select{
		Label: "Which environment is this secret going to be used?",
		Items: []string{"prod", "dev"},
	}
	_, env, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	fmt.Printf("Your created path %s/%s/%s/%s\n", service, env, secretService, secretName)

	c.Ui.Info("Success!")
	return 0
}

func (c *InteractiveCommand) Help() string {
	return `Usage: vc insert key1=value1 key2=value2...
  Inserts a new secret at the specified path with a set of key/value pairs.
`
}

func (c *InteractiveCommand) Synopsis() string {
	return "Insert an new secret"
}
