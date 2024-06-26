package container

import (
	"MiniDocker/common"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// 创建一个会隔离namespace进程的command
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, _ := os.Pipe()
	//调用自身，传入 init 参数，也就是执行initCommand
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := NewWorkSpace(common.RootPath, common.MntPath, volume)
	if err != nil {
		logrus.Errorf("new work space, err: %v", err)
	}
	cmd.ExtraFiles = []*os.File{
		readPipe,
	}

	return cmd, writePipe
}
