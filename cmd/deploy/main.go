package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var address string

var folder string

func Path(fileName string) string {
	return filepath.Join(folder, fileName)
}

func SetupFolder() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	folder = filepath.Dir(path)
	log.Printf("Folder: %s", folder)
}

func Run(command string, args ...string) error {
	log.Printf("Exec %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func RunFix(address, adminUsername, adminPassword, apexOnePassword string) error {
	log.Printf("Address: %s", address)
	psExecName := "PSExec64.exe"
	psExecPath := Path(psExecName)
	cmd := []string{
		psExecPath,
		"-c",
		"-f",
		"-u", adminUsername,
		"-p", adminPassword,
		"-w", `C:\`,
		`\\` + address,
		"fixsys", apexOnePassword,
	}
	return Run(psExecPath, cmd...)
}

func main() {
	SetupFolder()
	if len(os.Args) != 5 {
		log.Fatalf("usage: %s address username password apex_unload_password", os.Args[0])
	}
	address = os.Args[1]
	adminUsername := os.Args[2]
	adminPassword := os.Args[3]
	apexOnePassword := os.Args[4]
	err := RunFix(address, adminUsername, adminPassword, apexOnePassword)
	if err != nil {
		log.Println(err)
	}
}
