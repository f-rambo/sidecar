package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "cloud.c"
*/
import "C"

import (
	"context"
	"fmt"

	"github.com/f-rambo/ship/utils"
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
	if !utils.IsFileExist(crioFilePath) {
		return errors.New("crio file not found")
	}
	kubeadmPath := "/tmp/kubeadm"
	if !utils.IsFileExist(kubeadmPath) {
		return errors.New("kubeadm not found")
	}
	kubeletPath := "/tmp/kubelet"
	if !utils.IsFileExist(kubeletPath) {
		return errors.New("kubelet not found")
	}
	shipHomePath, err := utils.GetPackageStorePathByNames()
	if err != nil {
		return err
	}
	bash := utils.NewBash(uc.log)
	_, err = bash.RunCommand("tar -zxvf", crioFilePath, "-C", fmt.Sprintf("%s/crio", shipHomePath),
		"&& chmod +x", fmt.Sprintf("%s/crio/install", shipHomePath),
		"&& bash", fmt.Sprintf("%s/crio/install", shipHomePath))
	if err != nil {
		return err
	}
	_, err = bash.RunCommand("cp -r", kubeadmPath, "/usr/local/bin/kubeadm && chmod +x /usr/local/bin/kubeadm")
	if err != nil {
		return err
	}
	_, err = bash.RunCommand("cp -r", kubeletPath, "/usr/local/bin/kubelet && chmod +x /usr/local/bin/kubelet")
	if err != nil {
		return err
	}
	return nil
}

// add kubelet service and seting kubeadm config
func (uc *CloudUsecase) AddKubeletServiceAndSettingKubeadmConfig(ctx context.Context, cloud *Cloud) error {
	bash := utils.NewBash(uc.log)
	_, err := bash.RunCommand("echo", cloud.KubeadmConfig, "|", "sed", "s:/usr/bin:/usr/local/bin:g", "|", "sudo", "tee", "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf")
	if err != nil {
		return err
	}
	_, err = bash.RunCommand("echo", cloud.KubeletService, "|", "sed", "s:/usr/bin:/usr/local/bin:g", "|", "sudo", "tee", "/usr/lib/systemd/system/kubelet.service")
	if err != nil {
		return err
	}
	_, err = bash.RunCommand("systemctl daemon-reload && systemctl enable --now kubelet && systemctl status kubelet")
	if err != nil {
		return err
	}
	return nil
}

func (uc *CloudUsecase) KubeadmJoin(ctx context.Context, cloud *Cloud) error {
	joinCommand := fmt.Sprintf("kubeadm join %s --token %s --discovery-token-ca-cert-hash %s", cloud.ControlPlaneEndpoint, cloud.Token, cloud.DiscoveryTokenCACertHash)
	if cloud.JoinConfig != "" {
		joinCommand = fmt.Sprintf("kubeadm join --config %s", cloud.JoinConfig)
	}
	err := utils.NewBash(uc.log).RunCommandWithLogging(joinCommand)
	if err != nil {
		return err
	}
	return nil
}

func (uc *CloudUsecase) KubeadmInit(ctx context.Context, cloud *Cloud) error {
	initCommand := fmt.Sprintf("kubeadm init --config %s", cloud.KubeadmInitConfig)
	if cloud.KubeadmInitClusterConfig != "" {
		initCommand = fmt.Sprintf("kubeadm init --config %s --upload-config %s", cloud.KubeadmInitConfig, cloud.KubeadmInitClusterConfig)
	}
	err := utils.NewBash(uc.log).RunCommandWithLogging(initCommand)
	if err != nil {
		return err
	}
	return nil
}

// kubeadm reset
func (uc *CloudUsecase) KubeadmReset(ctx context.Context, cloud *Cloud) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("kubeadm", "reset", "--config", cloud.KubeadmResetConfig)
	if err != nil {
		return err
	}
	return nil
}

// kubeadm upgrade
func (uc *CloudUsecase) KubeadmUpgrade(ctx context.Context, cloud *Cloud) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("kubeadm", "upgrade", "apply", "--config", cloud.KubeadmUpgradeConfig)
	if err != nil {
		return err
	}
	return nil
}

// kubeadm upgrade plan
func (uc *CloudUsecase) KubeadmUpgradePlan(ctx context.Context, cloud *Cloud) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("kubeadm", "upgrade", "plan", "--config", cloud.KubeadmUpgradeConfig)
	if err != nil {
		return err
	}
	return nil
}

// seting ipv4 forward
func (uc *CloudUsecase) SetingIpv4Forward(ctx context.Context) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("sysctl", "-w", "net.ipv4.ip_forward=1")
	if err != nil {
		return err
	}
	return nil
}

// close swap
func (uc *CloudUsecase) CloseSwap(ctx context.Context) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("swapoff", "-a")
	if err != nil {
		return err
	}
	return nil
}

// CloseFirewall
func (uc *CloudUsecase) CloseFirewall(ctx context.Context) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("systemctl", "stop", "firewalld")
	if err != nil {
		return err
	}
	return nil
}
