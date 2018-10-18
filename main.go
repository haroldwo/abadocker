package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemHieMnt = "/sys/fs/cgroup/memory" //default Hierarchy of memory subsystem

func main() {
	if os.Args[0] == "/proc/self/exe" { //The static resource file of current process
		fmt.Printf("current pid %d", syscall.Getpid())
		os.Mkdir(path.Join(cgroupMemHieMnt, "testmemlmt"), 0755)
		ioutil.WriteFile(path.Join(cgroupMemHieMnt, "testmemlmt", "tasks"), []byte(strconv.Itoa(syscall.Getpid())), 0644)
		ioutil.WriteFile(path.Join(cgroupMemHieMnt, "testmemlmt", "memory.limit_in_bytes"), []byte("100m"), 0644) //Set limitation to 100 MB of memory
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, //Use SysProcAttr flags which has been set in kernel, these are the core features of container.
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.Process.Wait() //Wait for checking limitation. You can use "top" in another shell.
}
