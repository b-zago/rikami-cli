package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	SSMConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
)

func ParseEnvAsJSON(envPath string) []byte {
	f, err := os.ReadFile(envPath)
	if err != nil {
		log.Fatalf("Could not read %s. %q", envPath, err)
	}
	f = bytes.TrimSpace(f)

	lines := strings.Split(string(f), "\n")

	envs := make(map[string]string, len(lines))

	for _, line := range lines {
		key, value, _ := strings.Cut(line, "=")
		envs[key] = value
	}

	envJSON, err := json.Marshal(envs)
	if err != nil {
		log.Fatalf("Could not parse env file %s to JSON. %q", envPath, err)
	}

	return envJSON
}

func GetSSMClient(ctx context.Context) *ssm.Client {
	cfg, err := SSMConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("could not read AWS creds")
	}
	return ssm.NewFromConfig(cfg)
}

func PutParam(path, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := GetSSMClient(ctx)

	_, err := c.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(path),
		Overwrite: aws.Bool(true),
		Type:      types.ParameterTypeSecureString,
		Value:     aws.String(value),
	})
	if err != nil {
		if e, ok := errors.AsType[smithy.APIError](err); ok {
			// so it wont log something it shouldnt by accident
			return fmt.Errorf("could not put secret parameter %q\ncode=%s fault=%s", path, e.ErrorCode(), e.ErrorFault())
		} else {
			return fmt.Errorf("unknown error for parameter %q", path)
		}
	}

	return nil
}

func GetParam(arn string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := GetSSMClient(ctx)

	param, err := c.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(arn),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("get ssm parameter error %q: %w", arn, err)
	}

	return *param.Parameter.Value, nil
}
