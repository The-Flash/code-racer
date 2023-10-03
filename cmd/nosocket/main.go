package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	libseccomp "github.com/seccomp/libseccomp-golang"
	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(errors.New("no arguments provided"))
		os.Exit(1)
	}
	unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0)
	filter, err := libseccomp.NewFilter(libseccomp.ActAllow)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	filter.AddArch(libseccomp.ArchX32)
	filter.AddArch(libseccomp.ArchX86)

	syscallblacklist := []string{"socket"}

	for _, syscallStr := range syscallblacklist {
		socketSyscallID, err := libseccomp.GetSyscallFromName(syscallStr)
		if err != nil {
			log.Fatal(err)
		}

		filter.AddRule(socketSyscallID, libseccomp.ActKillProcess)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := filter.Load(); err != nil {
		log.Fatal(err)
	}
	// start child process
	output, err := exec.Command(os.Args[1], os.Args[2:]...).Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Stderr.WriteString(string(e.Stderr))
			os.Exit(e.ExitCode())
		default:
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	}
	fmt.Print(string(output))
}
