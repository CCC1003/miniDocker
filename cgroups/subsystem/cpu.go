package subsystem

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
	apply bool
}

func (*CpuSubSystem) Name() string {
	return "cpu"
}

func (c *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path,err:%v", cgroupPath, err)
		return err
	}

	if res.CpuShare != "" {
		c.apply = true
		//0644 表示文件权限为 -rw-r--r--

		//-：表示这是一个普通文件。
		//rw-：文件所有者具有读写权限。
		//r--：与文件所有者同组的用户具有读权限。
		//r--：其他用户具有读权限。
		err := os.WriteFile(path.Join(subsystemCgroupPath, "cpu.shares"), []byte(res.CpuSet), 0644)
		if err != nil {
			logrus.Errorf("failed to write file cpu.shares,err:%+v", err)
			return err
		}
	}
	return nil
}
func (c *CpuSubSystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.Remove(subsystemCgroupPath)
}

func (c *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if c.apply {
		subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
		if err != nil {
			return err
		}
		tasksPath := path.Join(subsystemCgroupPath, "tasks")
		err = os.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("write pid to tasks,path:%s,pid:%d,err:%v", tasksPath, pid, err)
			return err
		}
	}
	return nil
}
