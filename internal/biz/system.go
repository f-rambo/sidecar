package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "system.c"
*/
import "C"

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math"
	n "net"
	"os/exec"
	"strings"

	"github.com/f-rambo/ship/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// syscall

type System struct {
	ID         int64   `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	OS         string  `json:"os" gorm:"column:os; default:''; NOT NULL"`
	OSInfo     string  `json:"os_info" gorm:"column:os_info; default:''; NOT NULL"`
	ARCH       string  `json:"arch" gorm:"column:arch; default:''; NOT NULL"`
	CPU        int32   `json:"cpu" gorm:"column:cpu; default:0; NOT NULL"`
	Memory     float64 `json:"memory" gorm:"column:memory; default:0; NOT NULL"`
	GPU        int32   `json:"gpu" gorm:"column:gpu; default:0; NOT NULL"`
	GpuSpec    string  `json:"gpu_spec" gorm:"column:gpu_spec; default:''; NOT NULL"`
	DataDisk   int32   `json:"data_disk" gorm:"column:data_disk; default:0; NOT NULL"`
	Kernel     string  `json:"kernel" gorm:"column:kernel; default:''; NOT NULL"`
	Container  string  `json:"container" gorm:"column:container; default:''; NOT NULL"`
	Kubelet    string  `json:"kubelet" gorm:"column:kubelet; default:''; NOT NULL"`
	KubeProxy  string  `json:"kube_proxy" gorm:"column:kube_proxy; default:''; NOT NULL"`
	InternalIP string  `json:"internal_ip" gorm:"column:internal_ip; default:''; NOT NULL"`
	MachineID  string  `json:"machine_id" gorm:"column:machine_id; default:''; NOT NULL"`
}

type SystemRepo interface {
	GetSystem(ctx context.Context) (*System, error)
	SaveSystem(ctx context.Context, system *System) error
}

type SystemUsecase struct {
	SystemRepo SystemRepo
	log        *log.Helper
}

var ARCH_MAP = map[string]string{
	"x86_64":  "amd64",
	"aarch64": "arm64",
}

func NewSystemUseCase(systemRepo SystemRepo, logger log.Logger) *SystemUsecase {
	s := &SystemUsecase{
		SystemRepo: systemRepo,
		log:        log.NewHelper(logger),
	}
	return s
}

func (s *SystemUsecase) GetSystem(ctx context.Context) (*System, error) {
	system, err := s.SystemRepo.GetSystem(ctx)
	if err != nil {
		return nil, err
	}
	err = s.installSoftware("dmidecode")
	if err != nil {
		return nil, err
	}
	// check current user is root
	output, err := exec.Command("whoami").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	if strings.TrimSpace(string(output)) != "root" {
		// swtich to root user
		output, err = exec.Command("sudo", "-i").CombinedOutput()
		if err != nil {
			return nil, errors.Wrap(err, string(output))
		}
	}
	// get mac address
	output, err = exec.Command("dmidecode", "-s", "system-uuid").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	system.MachineID = strings.TrimSpace(string(output))
	// get system info
	output, err = exec.Command("uname", "-a").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	system.Kernel = strings.TrimSpace(string(output))
	// get system os
	output, err = exec.Command("cat", "/etc/os-release").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	system.OSInfo = strings.TrimSpace(string(output))
	output, err = exec.Command("uname", "-s").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	system.OS = strings.TrimSpace(string(output))
	output, err = exec.Command("uname", "-m").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	arch, ok := ARCH_MAP[strings.TrimSpace(string(output))]
	if !ok {
		return nil, errors.New("not support arch")
	}
	system.ARCH = arch
	cpu, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	system.CPU = int32(len(cpu))
	// memory info
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}
	system.Memory = math.Ceil(float64(v.Total) / 1024 / 1024 / 1024)
	// gpu info
	err = s.installSoftware("pciutils")
	if err != nil {
		return nil, err
	}
	output, err = exec.Command("lspci").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(output))
	}
	nvidiaLines := strings.Count(strings.ToLower(string(output)), "nvidia")
	system.GPU = int32(nvidiaLines)
	if system.GPU > 0 {
		// use output to get gpu spec
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(strings.ToLower(string(output)), "nvidia") {
				system.GpuSpec = strings.TrimSpace(line)
				break
			}
		}
	}
	// get data disk
	disk, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return nil, err
	}
	system.DataDisk = int32(math.Ceil(float64(disk.Total) / (1024 * 1024 * 1024)))
	// get internal ip
	netCards, err := net.InterfacesWithContext(ctx)
	if err != nil {
		return nil, err
	}
	for _, netCard := range netCards {
		if utils.ArrayContains(netCard.Flags, "up") && !utils.ArrayContains(netCard.Flags, "loopback") && len(netCard.Addrs) > 0 {
			ipCidr := netCard.Addrs[0].Addr
			ip, _, err := n.ParseCIDR(ipCidr)
			if err != nil {
				return nil, err
			}
			if ip != nil {
				system.InternalIP = ip.String()
			}
			break
		}
	}
	err = s.SystemRepo.SaveSystem(ctx, system)
	if err != nil {
		return nil, err
	}
	return system, nil
}

func (s *SystemUsecase) installSoftware(softwares ...string) error {
	// check linux version /etc/debian_version || /etc/redhat-release
	ok := utils.IsFileExist("/etc/debian_version")
	if ok {
		s.log.Info("debian")
		// check if software is already downloaded
		for _, software := range softwares {
			output, err := exec.Command("apt-cache", "policy", software).CombinedOutput()
			if err != nil {
				return errors.Wrap(err, string(output))
			}
			if !strings.Contains(string(output), "Installed: (none)") && string(output) != "" {
				s.log.Info("software already downloaded", software)
				continue
			}
			// update apt-get
			if err := s.runCommandWithLogging("apt", "update"); err != nil {
				return err
			}
			// install software
			if err := s.runCommandWithLogging("apt", "install", "-y", software); err != nil {
				return err
			}
		}
		return nil
	}
	ok = utils.IsFileExist("/etc/redhat-release")
	if ok {
		s.log.Info("redhat")
		// check if software is already downloaded
		for _, software := range softwares {
			output, err := exec.Command("yum", "list", "installed", software).CombinedOutput()
			if err != nil {
				return errors.Wrap(err, string(output))
			}
			if !strings.Contains(string(output), "Installed Packages") && string(output) != "" {
				s.log.Info("software already downloaded", software)
				continue
			}
			// update yum
			if err := s.runCommandWithLogging("sudo", "yum", "update"); err != nil {
				return err
			}
			// install software
			if err := s.runCommandWithLogging("sudo", "yum", "install", "-y", software); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("not support system")
}

func (s *SystemUsecase) runCommandWithLogging(command string, args ...string) error {
	s.log.Info("exec command: ", fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stderr pipe")
	}
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start command")
	}
	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			s.log.Info(scanner.Text())
		}
	}()
	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "command failed")
	}
	return nil
}

// exec command
func (s *SystemUsecase) execCommand(command string, args ...string) (output string, err error) {
	s.log.Info("exec command: ", fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	outputBytes, err := exec.Command(command, args...).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(outputBytes))
	}
	return string(outputBytes), err
}
