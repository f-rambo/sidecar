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

	"github.com/f-rambo/ship/internal/conf"
	"github.com/f-rambo/ship/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type Cloud struct {
	ID                       int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	NodeID                   int64  `json:"node_id" gorm:"column:node_id"`
	NodeName                 string `json:"node_name" gorm:"column:node_name"`
	ClusterVersion           string `json:"cluster_version" gorm:"column:cluster_version"`
	Token                    string `json:"token" gorm:"column:token"`
	DiscoveryTokenCACertHash string `json:"discovery_token_ca_cert_hash" gorm:"column:discovery_token_ca_cert_hash"`
	ControlPlaneEndpoint     string `json:"control_plane_endpoint" gorm:"column:control_plane_endpoint"`
	JoinConfig               string `json:"join_config" gorm:"column:join_config"`
}

var (
	NodeInitShell           string = "init.sh"
	KubernetesSowfwareShell string = "kubernetes-software.sh"
)

// kubeadm-init.conf
var KubeadmInitConfig = fmt.Sprintf(`apiVersion: kubeadm.k8s.io/v1beta4
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: %s
  bindPort: 6443`, "192.168.1.1")

// kubeadm-cluster.conf
var KubeadmClusterConfig = `apiVersion: kubeadm.k8s.io/v1beta4
kind: ClusterConfiguration
kubernetesVersion: v1.30.0
imageRepository: registry.aliyuncs.com/google_containers
controlPlaneEndpoint: "your-control-plane-endpoint:6443"
networking:
  podSubnet: "10.244.0.0/16"`

// kubeadm-join.conf
var KubeadmJoinConfig = `apiVersion: kubeadm.k8s.io/v1beta4
kind: JoinConfiguration
nodeRegistration:
  kubeletExtraArgs:
    node-labels: "node-role.kubernetes.io/master"`

var KubeadmResetConfig = `apiVersion: kubeadm.k8s.io/v1beta4
kind: ResetConfiguration`

var KubeadmUpgradeConfig = `apiVersion: kubeadm.k8s.io/v1beta4
kind: UpgradeConfiguration`

var KubeProxyConfig = `apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration`

var KubeletConfig = `apiVersion: kubelet.config.k8s.io/v1
kind: KubeletConfiguration`

type CloudRepo interface{}

type CloudUsecase struct {
	cloudRepo CloudRepo
	log       *log.Helper
	c         *conf.Server
}

func NewCloudUseCase(conf *conf.Server, cloudRepo CloudRepo, logger log.Logger) *CloudUsecase {
	c := &CloudUsecase{
		cloudRepo: cloudRepo,
		log:       log.NewHelper(logger),
		c:         conf,
	}
	return c
}

func (uc *CloudUsecase) NodeInit(ctx context.Context, cloud *Cloud) error {
	if !utils.IsFileExist(utils.MergePath(uc.c.Shell, NodeInitShell)) {
		return errors.New("init.sh not found")
	}
	return utils.NewBash(uc.log).RunCommandWithLogging("sudo bash", utils.MergePath(uc.c.Shell, NodeInitShell), cloud.NodeName)
}

func (uc *CloudUsecase) InstallKubeadmKubeletCriO(ctx context.Context, cloud *Cloud) error {
	if !utils.IsFileExist(uc.c.Resource) {
		return errors.New("resource not found")
	}
	if !utils.IsFileExist(utils.MergePath(uc.c.Shell, KubernetesSowfwareShell)) {
		return errors.New("kubernetes-software.sh not found")
	}
	return utils.NewBash(uc.log).RunCommandWithLogging("sudo bash", utils.MergePath(uc.c.Shell, KubernetesSowfwareShell), uc.c.Resource, cloud.ClusterVersion)
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
	// need write to file
	err := utils.NewBash(uc.log).RunCommandWithLogging(fmt.Sprintf("kubeadm init --config %s --upload-config %s", KubeadmInitConfig, KubeadmClusterConfig))
	if err != nil {
		return err
	}
	return nil
}

// kubeadm reset
func (uc *CloudUsecase) KubeadmReset(ctx context.Context, cloud *Cloud) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("kubeadm", "reset", "--config", KubeadmResetConfig)
	if err != nil {
		return err
	}
	return nil
}

// kubeadm upgrade
func (uc *CloudUsecase) KubeadmUpgrade(ctx context.Context, cloud *Cloud) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("kubeadm", "upgrade", "apply", "--config", KubeadmUpgradeConfig)
	if err != nil {
		return err
	}
	return nil
}

// kubeadm upgrade plan
func (uc *CloudUsecase) KubeadmUpgradePlan(ctx context.Context, cloud *Cloud) error {
	err := utils.NewBash(uc.log).RunCommandWithLogging("kubeadm", "upgrade", "plan", "--config", KubeadmUpgradeConfig)
	if err != nil {
		return err
	}
	return nil
}
