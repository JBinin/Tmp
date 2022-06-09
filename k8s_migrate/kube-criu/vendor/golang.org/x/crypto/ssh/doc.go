/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package ssh implements an SSH client and server.

SSH is a transport security protocol, an authentication protocol and a
family of application protocols. The most typical application level
protocol is a remote shell and this is specifically implemented.  However,
the multiplexed nature of SSH is exposed to users that wish to support
others.

References:
  [PROTOCOL.certkeys]: http://cvsweb.openbsd.org/cgi-bin/cvsweb/src/usr.bin/ssh/PROTOCOL.certkeys?rev=HEAD
  [SSH-PARAMETERS]:    http://www.iana.org/assignments/ssh-parameters/ssh-parameters.xml#ssh-parameters-1

This package does not fall under the stability promise of the Go language itself,
so its API may be changed when pressing needs arise.
*/
package ssh
