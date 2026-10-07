package main

import (
	"fmt"
	"log"
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	App            string `json:"app"`
	SSMPrefix      string `json:"ssmPrefix"`
	EnvFilesPrefix string `json:"envFilesPrefix"`
}

func main() {
	f, err := os.ReadFile("rikami.yml")
	if err != nil {
		log.Fatalf("Could not read rikami.yml. %q", err)
	}

	var c Config
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		log.Fatalf("Error parsing rikami.yml. %q", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("You need to provide a Command")
		os.Exit(2)
	}

	command := os.Args[1]

	switch command {
	case "params":

		if len(os.Args) < 3 {
			fmt.Println("You need to provide a SubCommand")
			os.Exit(2)
		}
		subCommand := os.Args[2]
		c.Params(subCommand)
	}
}
