package main

import (
	"MiniDocker/cgroups"
	"MiniDocker/cgroups/subsystem"
	"MiniDocker/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(cmdArray []string, tty bool, res *subsystem.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		logrus.Errorf("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("parent start failed,err : %v", err)
		return
	}

	//添加资源限制
	cgroupManager := cgroups.NewCGroupManager("go-docker")
	//删除资源限制
	defer cgroupManager.Destroy()
	//设置资源限制
	cgroupManager.Set(res)
	//将容器进程，加入到各个subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)

	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
