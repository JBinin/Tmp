#
#Copyright (c) 2014-2020 CGCL Labs
#Container_Migrate is licensed under Mulan PSL v2.
#You can use this software according to the terms and conditions of the Mulan PSL v2.
#You may obtain a copy of Mulan PSL v2 at:
#         http://license.coscl.org.cn/MulanPSL2
#THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
#EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
#MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
#See the Mulan PSL v2 for more details.
#
Name: kubeadm
Version: OVERRIDE_THIS
Release: 00
License: ASL 2.0
Summary: Container Cluster Manager - Kubernetes Cluster Bootstrapping Tool
Requires: kubelet >= 1.8.0
Requires: kubectl >= 1.8.0
Requires: kubernetes-cni >= 0.5.1
Requires: cri-tools >= 1.11.0

URL: https://kubernetes.io

%description
Command-line utility for deploying a Kubernetes cluster.

%install
install -m 755 -d %{buildroot}%{_bindir}
install -m 755 -d %{buildroot}%{_sysconfdir}/systemd/system/
install -m 755 -d %{buildroot}%{_sysconfdir}/systemd/system/kubelet.service.d/
install -m 755 -d %{buildroot}%{_sysconfdir}/sysconfig/
install -p -m 755 -t %{buildroot}%{_bindir} {kubeadm}
install -p -m 755 -t %{buildroot}%{_sysconfdir}/systemd/system/kubelet.service.d/ {10-kubeadm.conf}
install -p -m 755 -T {kubelet.env} %{buildroot}%{_sysconfdir}/sysconfig/kubelet

%files
%{_bindir}/kubeadm
%{_sysconfdir}/systemd/system/kubelet.service.d/10-kubeadm.conf
%{_sysconfdir}/sysconfig/kubelet
