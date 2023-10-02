package main

// #cgo pkg-config: libseccomp
/*
#include <stdlib.h>
#include <seccomp.h>

const uint32_t C_SCMP_ARCH_X86_64 = SCMP_ARCH_X86_64;
*/
// import "C"
import (
	"log"
	"net/http"

	libseccomp "github.com/seccomp/libseccomp-golang"
	"golang.org/x/sys/unix"
)

// type nativeArch uint32

func main() {
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
	// filter.AddArch(libseccomp.ArchAMD64)
	// filter.AddArch(libseccomp.ArchNative)
	// fmt.Println("Sys call Id", syscallID)
	syscallID, err := libseccomp.GetSyscallFromName("socket")
	if err != nil {
		log.Fatal(err)
	}

	filter.AddRule(syscallID, libseccomp.ActKillProcess)
	if err != nil {
		log.Fatal(err)
	}

	if err := filter.Load(); err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get("https://google.com")
	if err != nil {
		log.Fatal("Request err", resp)
	}
	log.Println(resp.Body)
}
