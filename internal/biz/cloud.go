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
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type Cloud struct {
	ID                       int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	KubeProxyConfig          string `json:"kube_proxy_config" gorm:"column:kube_proxy_config"`
	KubeletConfig            string `json:"kubelet_config" gorm:"column:kubelet_config"`
	KubeadmInitClusterConfig string `json:"kubeadm_init_cluster_config" gorm:"column:kubeadm_init_cluster_config"`
	KubeadmInitConfig        string `json:"kubeadm_init_config" gorm:"column:kubeadm_init_config"`
	KubeadmResetConfig       string `json:"kubeadm_reset_config" gorm:"column:kubeadm_reset_config"`
	KubeadmUpgradeConfig     string `json:"kubeadm_upgrade_config" gorm:"column:kubeadm_upgrade_config"`
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
	output, err := uc.execCommand("sudo", "tar", "-zxvf", crioFilePath, "-C", "$HOME/crio")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to install crio: %v, output: %s", err, output)
		return err
	}
	output, err = uc.execCommand("sudo", "chmod", "+x", "$HOME/crio/install")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to chmod crio: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("chmod crio success")
	output, err = uc.execCommand("sudo", "bash", "$HOME/crio/install")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to install crio: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("install crio success......")
	output, err = uc.execCommand("sudo", "mv", "/tmp/kubeadm", "/usr/local/bin/kubeadm")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to move kubeadm: %v, output: %s", err, output)
		return err
	}
	output, err = uc.execCommand("sudo", "chmod", "+x", "/usr/local/bin/kubeadm")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to chmod kubeadm: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("chmod kubeadm success")
	output, err = uc.execCommand("sudo", "mv", "/tmp/kubelet", "/usr/local/bin/kubelet")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to move kubelet: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("move kubelet success")
	output, err = uc.execCommand("sudo", "chmod", "+x", "/usr/local/bin/kubelet")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to chmod kubelet: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("chmod kubelet success")
	return nil
}

// add kubelet service and seting kubeadm config
func (uc *CloudUsecase) AddKubeletServiceAndSettingKubeadmConfig(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("add kubelet service and setting kubeadm config")
	output, err := uc.execCommand("echo", cloud.KubeadmConfig, "|", "sed", "s:/usr/bin:/usr/local/bin:g", "|", "sudo", "tee", "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to add kubelet service and setting kubeadm config: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("add kubelet service and setting kubeadm config success")
	output, err = uc.execCommand("echo", cloud.KubeletService, "|", "sed", "s:/usr/bin:/usr/local/bin:g", "|", "sudo", "tee", "/usr/lib/systemd/system/kubelet.service")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to add kubelet service: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("add kubelet service success")
	output, err = uc.execCommand("sudo", "systemctl", "daemon-reload")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to reload daemon: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("reload daemon success")
	output, err = uc.execCommand("sudo", "systemctl", "enable", "--now", "kubelet")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to enable kubelet: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("enable kubelet success")
	output, err = uc.execCommand("sudo", "systemctl", "status", "kubelet")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to get kubelet status: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("get kubelet status success")
	return nil
}

func (uc *CloudUsecase) KubeadmJoin(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("joining node to Kubernetes cluster using join config file...")

	joinCommand := fmt.Sprintf("sudo kubeadm join %s --token %s --discovery-token-ca-cert-hash %s", cloud.ControlPlaneEndpoint, cloud.Token, cloud.DiscoveryTokenCACertHash)
	if cloud.JoinConfig != "" {
		joinCommand = fmt.Sprintf("sudo kubeadm join --config %s", cloud.JoinConfig)
	}

	err := uc.runCommandWithLogging("bash", "-c", joinCommand)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to join node to Kubernetes cluster: %v", err)
		return err
	}

	uc.log.WithContext(ctx).Info("node joined to Kubernetes cluster successfully")
	return nil
}

func (uc *CloudUsecase) KubeadmInit(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("initializing Kubernetes cluster...")

	initCommand := fmt.Sprintf("sudo kubeadm init --config %s", cloud.KubeadmInitConfig)
	if cloud.KubeadmInitClusterConfig != "" {
		initCommand = fmt.Sprintf("sudo kubeadm init --config %s --upload-config %s", cloud.KubeadmInitConfig, cloud.KubeadmInitClusterConfig)
	}

	err := uc.runCommandWithLogging("bash", "-c", initCommand)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to initialize Kubernetes cluster: %v", err)
		return err
	}

	uc.log.WithContext(ctx).Info("Kubernetes cluster initialized successfully")
	return nil
}

// kubeadm reset
func (uc *CloudUsecase) KubeadmReset(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("resetting Kubernetes cluster...")
	output, err := uc.execCommand("sudo", "kubeadm", "reset", "--config", cloud.KubeadmResetConfig)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to reset Kubernetes cluster: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("resetting Kubernetes cluster success")
	return nil
}

// kubeadm upgrade
func (uc *CloudUsecase) KubeadmUpgrade(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("upgrading Kubernetes cluster...")
	output, err := uc.execCommand("sudo", "kubeadm", "upgrade", "apply", "--config", cloud.KubeadmUpgradeConfig)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to upgrade Kubernetes cluster: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("upgrading Kubernetes cluster success")
	return nil
}

// kubeadm upgrade plan
func (uc *CloudUsecase) KubeadmUpgradePlan(ctx context.Context, cloud *Cloud) error {
	uc.log.WithContext(ctx).Info("getting Kubernetes cluster upgrade plan...")
	output, err := uc.execCommand("sudo", "kubeadm", "upgrade", "plan", "--config", cloud.KubeadmUpgradeConfig)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to get Kubernetes cluster upgrade plan: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("getting Kubernetes cluster upgrade plan success")
	return nil
}

// seting ipv4 forward
func (uc *CloudUsecase) SetingIpv4Forward(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("seting ipv4 forward...")
	output, err := uc.execCommand("sudo", "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to seting ipv4 forward: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("seting ipv4 forward success")
	return nil
}

// close swap
func (uc *CloudUsecase) CloseSwap(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("closing swap...")
	output, err := uc.execCommand("sudo", "swapoff", "-a")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to close swap: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("closing swap success")
	return nil
}

// CloseFirewall
func (uc *CloudUsecase) CloseFirewall(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("closing firewall...")
	output, err := uc.execCommand("sudo", "systemctl", "stop", "firewalld")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to close firewall: %v, output: %s", err, output)
		return err
	}
	uc.log.WithContext(ctx).Info("closing firewall success")
	return nil
}

func (uc *CloudUsecase) runCommandWithLogging(command string, args ...string) error {
	uc.log.Info("exec command: ", fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
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
			uc.log.Info(scanner.Text())
		}
	}()
	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "command failed")
	}
	return nil
}

// exec command
func (uc *CloudUsecase) execCommand(command string, args ...string) (output string, err error) {
	uc.log.Info("exec command: ", fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	outputBytes, err := exec.Command(command, args...).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(outputBytes))
	}
	return string(outputBytes), err
}
