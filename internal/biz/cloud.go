package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "cloud.c"
*/
import "C"

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/go-kratos/kratos/v2/log"
)

type Cloud struct {
	ID                       int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	KubeadmInitConfig        string `json:"kubeadm_init_config" gorm:"column:kubeadm_init_config"`
	KubeadmConfig            string `json:"kubeadm_config" gorm:"column:kubeadm_config"`
	KubeletService           string `json:"kubelet_service" gorm:"column:kubelet_service"`
	CRIOVersion              string `json:"crio_version" gorm:"column:crio_version"`
	ARCH                     string `json:"arch" gorm:"column:arch"`
	Token                    string `json:"token" gorm:"column:token"`
	DiscoveryTokenCACertHash string `json:"discovery_token_ca_cert_hash" gorm:"column:discovery_token_ca_cert_hash"`
	ControlPlaneEndpoint     string `json:"control_plane_endpoint" gorm:"column:control_plane_endpoint"`
	JoinConfig               string `json:"join_config" gorm:"column:join_config"`
}

type CloudRepo interface{}

type CloudUsecase struct {
	cloudRepo CloudRepo
	log       *log.Helper
}

func NewCloudUseCase(cloudRepo CloudRepo, logger log.Logger) *CloudUsecase {
	c := &CloudUsecase{
		cloudRepo: cloudRepo,
		log:       log.NewHelper(logger),
	}
	return c
}

// install kubeadm kubelet cri-o
func (uc *CloudUsecase) InstallKubeadmKubeletCriO(ctx context.Context, cloud *Cloud) error {
	crioFilePath := fmt.Sprintf("/tmp/crio.%s.v%s.tar.gz", cloud.ARCH, cloud.CRIOVersion)
	uc.log.WithContext(ctx).Info("install crio......")
	output, err := exec.Command("sudo", "tar", "-zxvf", crioFilePath, "-C", "$HOME/crio").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to install crio: %v, output: %s", err, string(output))
		return err
	}
	output, err = exec.Command("sudo", "chmod", "+x", "$HOME/crio/install").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to chmod crio: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("chmod crio success")
	output, err = exec.Command("sudo", "bash", "$HOME/crio/install").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to install crio: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("install crio success......")
	output, err = exec.Command("sudo", "mv", "/tmp/kubeadm", "/usr/local/bin/kubeadm").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to move kubeadm: %v, output: %s", err, string(output))
		return err
	}
	output, err = exec.Command("sudo", "chmod", "+x", "/usr/local/bin/kubeadm").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to chmod kubeadm: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("chmod kubeadm success")
	output, err = exec.Command("sudo", "mv", "/tmp/kubelet", "/usr/local/bin/kubelet").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to move kubelet: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("move kubelet success")
	output, err = exec.Command("sudo", "chmod", "+x", "/usr/local/bin/kubelet").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to chmod kubelet: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("chmod kubelet success")
	return nil
}

// add kubelet service and seting kubeadm config
func (uc *CloudUsecase) AddKubeletServiceAndSettingKubeadmConfig(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("add kubelet service and setting kubeadm config")
	output, err := exec.Command("echo", cloud.KubeadmConfig, "|", "sed", "s:/usr/bin:/usr/local/bin:g", "|", "sudo", "tee", "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to add kubelet service and setting kubeadm config: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("add kubelet service and setting kubeadm config success")
	output, err = exec.Command("echo", cloud.KubeletService, "|", "sed", "s:/usr/bin:/usr/local/bin:g", "|", "sudo", "tee", "/usr/lib/systemd/system/kubelet.service").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to add kubelet service: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("add kubelet service success")
	output, err = exec.Command("sudo", "systemctl", "daemon-reload").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to reload daemon: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("reload daemon success")
	output, err = exec.Command("sudo", "systemctl", "enable", "--now", "kubelet").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to enable kubelet: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("enable kubelet success")
	output, err = exec.Command("sudo", "systemctl", "status", "kubelet").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to get kubelet status: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("get kubelet status success")
	return nil
}

func (uc *CloudUsecase) JoinKubeadmWithJoinCommand(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("joining node to Kubernetes cluster...")

	joinCommand := fmt.Sprintf("sudo kubeadm join %s --token %s --discovery-token-ca-cert-hash %s", cloud.ControlPlaneEndpoint, cloud.Token, cloud.DiscoveryTokenCACertHash)
	cmd := exec.Command("bash", "-c", joinCommand)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to create stdout pipe: %v", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to create stderr pipe: %v", err)
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to start command: %v", err)
		return err
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
	for scanner.Scan() {
		line := scanner.Text()
		uc.log.WithContext(ctx).Info(line)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to join node to Kubernetes cluster: %v", err)
		return err
	}

	uc.log.WithContext(ctx).Info("node joined to Kubernetes cluster successfully")
	return nil
}

func (uc *CloudUsecase) JoinKubeadmWithJoinConfigFile(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("joining node to Kubernetes cluster using join config file...")

	joinCommand := fmt.Sprintf("sudo kubeadm join --config %s", cloud.JoinConfig)
	cmd := exec.Command("bash", "-c", joinCommand)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to create stdout pipe: %v", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to create stderr pipe: %v", err)
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to start command: %v", err)
		return err
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
	for scanner.Scan() {
		line := scanner.Text()
		uc.log.WithContext(ctx).Info(line)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to join node to Kubernetes cluster: %v", err)
		return err
	}

	uc.log.WithContext(ctx).Info("node joined to Kubernetes cluster successfully")
	return nil
}

func (uc *CloudUsecase) InitKubeadm(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("initializing Kubernetes cluster...")

	initCommand := fmt.Sprintf("sudo kubeadm init --config %s", cloud.KubeadmInitConfig)
	cmd := exec.Command("bash", "-c", initCommand)

	// Create a pipe for the command's stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to create stdout pipe: %v", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to create stderr pipe: %v", err)
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to start command: %v", err)
		return err
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
	for scanner.Scan() {
		uc.log.WithContext(ctx).Info(scanner.Text())
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to initialize Kubernetes cluster: %v", err)
		return err
	}

	uc.log.WithContext(ctx).Info("Kubernetes cluster initialized successfully")
	return nil
}

// seting ipv4 forward
func (uc *CloudUsecase) SetingIpv4Forward(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("seting ipv4 forward...")
	output, err := exec.Command("sudo", "sysctl", "-w", "net.ipv4.ip_forward=1").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to seting ipv4 forward: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("seting ipv4 forward success")
	return nil
}

// close swap
func (uc *CloudUsecase) CloseSwap(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("closing swap...")
	output, err := exec.Command("sudo", "swapoff", "-a").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to close swap: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("closing swap success")
	return nil
}

// CloseFirewall
func (uc *CloudUsecase) CloseFirewall(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("closing firewall...")
	output, err := exec.Command("sudo", "systemctl", "stop", "firewalld").CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to close firewall: %v, output: %s", err, string(output))
		return err
	}
	uc.log.WithContext(ctx).Info("closing firewall success")
	return nil
}
