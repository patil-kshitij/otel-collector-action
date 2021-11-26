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
	go runCollector()
	if err != nil {
		fmt.Println("error in running collector :", err)
	}
	time.Sleep(30 * time.Second)
	err = executeCommand()
	time.Sleep(30 * time.Second)
	fmt.Println("error while executing command :", err)
	//collectorCMD.Process.
	// fmt.Println("Collector Output :", collectorOutput.String())
	// fmt.Println("collector err :", collectorErr.String())
	err = collectorCMD.Process.Kill()
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
	cmd := exec.Command("cp", "config.yaml", "/go/src/otel/bin/config.yaml")
	//cmd.Dir = "./otel/bin/"
	return cmd.Run()
}

func runCollector() error {
	collectorCMD = exec.Command("./cmd-otelcol", "--config=/go/src/config.yaml")
	collectorCMD.Dir = "./otel/bin"
	// collectorCMD.Stdout = &collectorOutput
	// collectorCMD.Stdin = &collectorErr
	output, err := collectorCMD.CombinedOutput()
	fmt.Println("Collector Output :", string(output))
	return err
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
