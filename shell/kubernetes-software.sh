#!/bin/bash
set -e

if [ -z "${1// /}" ] || [ -z "${2// /}" ]; then
  echo "Usage: $0 <RESOURCE> <CLUSTER_VERSION>"
  exit 1
fi

RESOURCE=$1
CLUSTER_VERSION=$2

if [[ ! $(readlink -f "$RESOURCE") =~ ^/ ]]; then
  echo "Error: RESOURCE must be an absolute path"
  exit 1
fi

if [[ ! $CLUSTER_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: CLUSTER_VERSION must be in the format X.Y.Z"
  exit 1
fi

if [ ! -d "$RESOURCE" ] || [ ! -r "$RESOURCE" ]; then
  echo "Error: RESOURCE directory $RESOURCE does not exist or is not readable"
  exit 1
fi

if [ ! -d "$RESOURCE/kubernetes-software/$CLUSTER_VERSION" ]; then
  echo "Error: CLUSTER_VERSION $CLUSTER_VERSION does not exist in the RESOURCE directory"
  exit 1
fi

ARCH=$(uname -m)
case $ARCH in
aarch64)
  ARCH="arm64"
  ;;
x86_64)
  ARCH="amd64"
  ;;
*)
  echo "Error: Unsupported architecture $ARCH. Supported architectures are: aarch64, x86_64"
  exit 1
  ;;
esac

kubeadmConfig=$(
  cat <<EOF
# Note: This dropin only works with kubeadm and kubelet v1.11+
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
# This is a file that "kubeadm init" and "kubeadm join" generates at runtime, populating the KUBELET_KUBEADM_ARGS variable dynamically
EnvironmentFile=-/var/lib/kubelet/kubeadm-flags.env
# This is a file that the user can use for overrides of the kubelet args as a last resort. Preferably, the user should use
# the .NodeRegistration.KubeletExtraArgs object in the configuration files instead. KUBELET_EXTRA_ARGS should be sourced from this file.
EnvironmentFile=-/etc/sysconfig/kubelet
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
EOF
)

kubeletService=$(
  cat <<EOF
[Unit]
Description=kubelet: The Kubernetes Node Agent
Documentation=https://kubernetes.io/docs/
Wants=network-online.target
After=network-online.target

[Service]
ExecStart=/usr/bin/kubelet
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
)

function install_crio() {
  crioPath="$RESOURCE/kubernetes-software/$CLUSTER_VERSION/crio/$ARCH/cri-o"

  if [ ! -d "$crioPath" ] || [ ! -r "$crioPath" ]; then
    echo "Error: Directory $crioPath does not exist or is not readable"
    exit 1
  fi

  if [ ! -f "$crioPath/install" ]; then
    echo "Error: File $crioPath/install does not exist"
    exit 1
  fi

  # 确保 install 脚本具有执行权限
  if [ ! -x "$crioPath/install" ]; then
    chmod +x "$crioPath/install"
  fi

  cd "$crioPath" && ./install && cd -

  crictlPath="$RESOURCE/kubernetes-software/$CLUSTER_VERSION/crictl/$ARCH/crictl"
  if [ ! -d "$crictlPath" ] || [ ! -r "$crictlPath" ]; then
    echo "Error: Directory $crictlPath does not exist or is not readable"
    exit 1
  fi

  if [ ! -f "$crictlPath/crictl" ]; then
    echo "Error: File $crictlPath/crictl does not exist"
    exit 1
  fi

  mv "$crictlPath/crictl" /usr/local/bin/crictl && chmod +x /usr/local/bin/crictl

}

function install_kubernetes_software() {
  kubernetesPath="$RESOURCE/kubernetes-software/$CLUSTER_VERSION/kubernetes/$ARCH"

  if [ ! -d "$kubernetesPath" ] || [ ! -r "$kubernetesPath" ]; then
    echo "Error: Directory $kubernetesPath does not exist or is not readable"
    exit 1
  fi

  if [ ! -f "$kubernetesPath/kubeadm" ]; then
    echo "Error: File $kubernetesPath/kubeadm does not exist"
    exit 1
  fi

  if [ ! -x "$kubernetesPath/kubeadm" ]; then
    chmod +x "$kubernetesPath/kubeadm"
  fi

  mv "$kubernetesPath/kubeadm" /usr/local/bin/kubeadm

  if [ ! -f "$kubernetesPath/kubectl" ]; then
    echo "Error: File $kubernetesPath/kubectl does not exist"
    exit 1
  fi

  if [ ! -x "$kubernetesPath/kubectl" ]; then
    chmod +x "$kubernetesPath/kubectl"
  fi

  mv "$kubernetesPath/kubectl" /usr/local/bin/kubectl

  if ! echo "$kubeadmConfig" | sed "s:/usr/bin:/usr/local/bin:g" | tee /usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf >/dev/null; then
    echo "Error: Failed to write to /usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"
    exit 1
  fi

  if ! echo "$kubeletService" | sed "s:/usr/bin:/usr/local/bin:g" | tee /usr/lib/systemd/system/kubelet.service >/dev/null; then
    echo "Error: Failed to write to /usr/lib/systemd/system/kubelet.service"
    exit 1
  fi

  if ! systemctl daemon-reload || ! systemctl enable --now kubelet; then
    echo "Error: Failed to start kubelet service"
    exit 1
  fi
}

function install_kubernetes_images() {
  kubernetes_images_path="$RESOURCE/kubernetes-software/$CLUSTER_VERSION/kubernetes/${ARCH}/images/kubernetes-images.tar"
}

if systemctl is-active --quiet crio; then
  echo "crio is already running, skipping installation."
else
  echo "crio is not running, proceeding with installation."
  install_crio
fi

if systemctl is-active --quiet kubelet; then
  echo "kubelet is already running, skipping installation."
else
  echo "kubelet is not running, proceeding with installation."
  install_kubernetes_software
fi

install_kubernetes_images
