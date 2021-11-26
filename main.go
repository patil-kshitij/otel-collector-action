package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	git "github.com/go-git/go-git/v5"
)

var (
	collectorCMD    *exec.Cmd
	collectorOutput bytes.Buffer
	collectorErr    bytes.Buffer
)

func main() {
	var err error
	//pull otel-code
	pullOtelCollector()

	//build
	buildCollector()

	//create config in bin dir
	createOtelConfig()

	//run bin
	err = runCollector()
	if err != nil {
		fmt.Println("error in running collector :", err)
	}
	fmt.Println("About to wait for collector")
	for collectorCMD.Process == nil {

	}
	time.Sleep(30 * time.Second)
	err = executeCommand()
	time.Sleep(30 * time.Second)
	fmt.Println("error while executing command :", err)
	//collectorCMD.Process.

	err = collectorCMD.Process.Kill()
	fmt.Println("Collector Output :", collectorOutput.String())
	fmt.Println("collector err :", collectorErr.String())
	if err != nil {
		fmt.Println("error in killing process")
	}

}

func pullOtelCollector() error {
	directory := "otel"
	url := "https://github.com/open-telemetry/opentelemetry-collector.git"
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	return err

}

func buildCollector() error {
	cmd := exec.Command("make", "build-binary-cmd-otelcol")
	cmd.Dir = "./otel"
	return cmd.Run()

}

func createOtelConfig() error {
	cmd := exec.Command("cp", "config.yaml", "./otel/bin/config.yaml")
	//cmd.Dir = "./otel/bin/"
	return cmd.Run()
}

func runCollector() error {
	collectorCMD = exec.Command("./cmd-otelcol", "--config=config.yaml")
	collectorCMD.Dir = "./otel/bin"
	collectorCMD.Stdout = &collectorOutput
	collectorCMD.Stdin = &collectorErr
	return collectorCMD.Start()
}

func executeCommand() error {
	cmd := exec.Command("go", "test", "./...")
	repoDir := os.Getenv("GITHUB_WORKSPACE")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("error in getting std output for command :", err)
		return err
	}
	fmt.Println("Command Output :", string(output))
	return nil
}
