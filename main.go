package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Println("box starting...")
	fmt.Printf("arg1: %s\n", os.Args[1])
	if "run" == os.Args[1] {
		run()
	} else if "init" == os.Args[1] {
		initContainer()
	} else {
		fmt.Printf("unknow arg: %s\n", os.Args[1])
	}
}

func run() {
	fmt.Println("box.run")
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      syscall.Getuid(),
				Size:        1100,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      syscall.Getgid(),
				Size:        1100,
			},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
	os.Exit(-1)
}

func initContainer() {
	fmt.Println("box.init")
	mountFlag := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(mountFlag), "")
	if err := syscall.Exec("/bin/sh", []string{"-it"}, os.Environ()); err != nil {
		fmt.Printf("exec error: %s\n", err.Error())
	}
}
