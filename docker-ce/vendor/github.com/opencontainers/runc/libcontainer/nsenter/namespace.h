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
#ifndef NSENTER_NAMESPACE_H
#define NSENTER_NAMESPACE_H

#ifndef _GNU_SOURCE
#	define _GNU_SOURCE
#endif
#include <sched.h>

/* All of these are taken from include/uapi/linux/sched.h */
#ifndef CLONE_NEWNS
#	define CLONE_NEWNS 0x00020000 /* New mount namespace group */
#endif
#ifndef CLONE_NEWCGROUP
#	define CLONE_NEWCGROUP 0x02000000 /* New cgroup namespace */
#endif
#ifndef CLONE_NEWUTS
#	define CLONE_NEWUTS 0x04000000 /* New utsname namespace */
#endif
#ifndef CLONE_NEWIPC
#	define CLONE_NEWIPC 0x08000000 /* New ipc namespace */
#endif
#ifndef CLONE_NEWUSER
#	define CLONE_NEWUSER 0x10000000 /* New user namespace */
#endif
#ifndef CLONE_NEWPID
#	define CLONE_NEWPID 0x20000000 /* New pid namespace */
#endif
#ifndef CLONE_NEWNET
#	define CLONE_NEWNET 0x40000000 /* New network namespace */
#endif

#endif /* NSENTER_NAMESPACE_H */
