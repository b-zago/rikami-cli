package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (c *Config) Params(subCommand string) {
	switch subCommand {
	case "get":
		if len(os.Args) < 4 {
			fmt.Print("You need to provide a Path")
			os.Exit(2)
		}
		path := os.Args[3]

		param, err := GetParam(path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(param)
		// just get params as json wit optional flag to make output .env-like
	case "put":
		if len(os.Args) < 4 {
			fmt.Print("You need to provide an EnvPath")
			os.Exit(2)
		}
		if len(os.Args) < 5 {
			fmt.Print("You need to provide a SecretName")
			os.Exit(2)
		}

		envPath := os.Args[3]
		secretName := os.Args[4]

		envPathStripped, _ := strings.CutPrefix(filepath.Base(envPath), c.EnvFilesPrefix)
		envSSMPath := strings.ReplaceAll(envPathStripped, ".", "/")
		ssmParam := c.SSMPrefix + envSSMPath + "/" + c.App + "/" + secretName

		parsedEnv := ParseEnvAsJSON(envPath)
		err := PutParam(ssmParam, string(parsedEnv))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Uploaded to SSM path ", ssmParam)
		// upload to ssm by parsing file like .env.staging to json in ssm path /(config_prefix)/staging/(config_app_name)/(arg_secret_name)
	}
}
