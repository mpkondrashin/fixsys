package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	ApexOnePassword string
)

//go:embed sakfile.sys
var sakfile []byte

const (
	step1Flag = "step1.txt"
	step2Flag = "step2.txt"
	logFile   = "fixsys.log"
	aOnePath  = `C:\Program Files (x86)\Trend Micro\Security Agent\PccNTMon.exe`
)

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
}

func SakExists() error {
	log.Println("Check that sakfile.sys exists")
	if !FileExists(filepath.Join(DriversFolder, "sakfile.sys")) {
		return fmt.Errorf("%s is missing - existing", filepath.Join(DriversFolder, "sakfile.sys"))
	}
	return nil
}

func Restart() error {
	log.Print("restart")
	return Run("cmd", "/C", "shutdown", "/r")
}

func Run(command string, args ...string) error {
	log.Printf("Exec %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateFile(filePath string) error {
	log.Printf("Create %s", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func Step01() error {
	if err := SakExists(); err != nil {
		return err
	}
	if err := Run("bcdedit", "/set", "testsigning", "on"); err != nil {
		return err
	}
	if err := CreateFile(Path(step1Flag)); err != nil {
		return nil
	}
	return Restart()
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false

}

func UnloadAOne() error {
	log.Println("UnloadAOne")
	return Run(aOnePath, "-n", ApexOnePassword)
}

func StopDriver() error {
	log.Println("Unload driver")
	return Run("sc", "stop", "sakfile")

}

var DriversFolder = `C:\Windows\system32\drivers`

func RenameDriver() error {
	log.Println("Rename driver")
	orig := filepath.Join(DriversFolder, "sakfile.sys")
	new := filepath.Join(DriversFolder, "sakfile.sys.bak")
	return os.Rename(orig, new)

}

func WriteSakFile() error {
	log.Println("Write new sakfile.sys")
	path := filepath.Join(DriversFolder, "sakfile.sys")
	return os.WriteFile(path, sakfile, 0666)
}

func LoadAOne() error {
	log.Println("Load AOne")
	return Run(aOnePath)
}

func Step02() error {
	if err := UnloadAOne(); err != nil {
		return err
	}
	s := 25 * time.Second
	log.Printf("Wait to upload. Sleep for %s", s)
	time.Sleep(s)
	if err := StopDriver(); err != nil {
		if !strings.Contains(err.Error(), "exit status 1062") {
			return err
		}
		log.Println(err)
	}
	if err := RenameDriver(); err != nil {
		return err
	}
	if err := WriteSakFile(); err != nil {
		return err
	}
	go func() {
		s := 10 * time.Second
		log.Printf("Wait to Load. Sleep for %s", s)
		time.Sleep(s)
		if err := Run("bcdedit", "/set", "testsigning", "off"); err != nil {
			log.Print(err)
		}
		CreateFile(Path(step2Flag))
		os.Exit(0)
	}()
	if err := LoadAOne(); err != nil {
		return err
	}
	s = 20 * time.Second
	log.Printf("Wait to run bcedit. Sleep for %s", s)
	time.Sleep(s)
	return nil
}

func main() {
	SetupFolder()
	logFile, err := os.OpenFile(filepath.Join(folder, logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stderr))
	log.Println("Started")
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
		}
	}()
	if len(os.Args) != 2 {
		log.Fatal("Missing ApexOne passsword parameter")
	}
	ApexOnePassword = os.Args[1]

	if !FileExists(Path(step1Flag)) {
		log.Printf("%s not found", Path(step1Flag))
		err := Step01() // it should not return control
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if !FileExists(Path(step2Flag)) {
		log.Printf("%s not found", Path(step2Flag))
		err := Step02()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Printf("%s and %s exist - exiting", Path(step1Flag), Path(step2Flag))
}
