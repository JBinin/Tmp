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
Name: cri-tools
Version: OVERRIDE_THIS
Release: 00
License: ASL 2.0
Summary: Container Runtime Interface tools

URL: https://kubernetes.io

%description
Binaries to interface with the container runtime.

%prep
# This has to be hard coded because bazel does a path substitution before rpm's %{version} is substituted.
tar -xzf {crictl-v1.11.1-linux-amd64.tar.gz}

%install
install -m 755 -d %{buildroot}%{_bindir}
install -p -m 755 -t %{buildroot}%{_bindir} crictl

%files
%{_bindir}/crictl
