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
#define SECRET_WITH_UNSTABLE 1
#define SECRET_API_SUBJECT_TO_CHANGE 1
#include <libsecret/secret.h>

const SecretSchema *docker_get_schema(void) G_GNUC_CONST;

#define DOCKER_SCHEMA docker_get_schema()

GError *add(char *server, char *username, char *secret);
GError *delete(char *server);
GError *get(char *server, char **username, char **secret);
GError *list(char *** paths, char *** accts, unsigned int *list_l);
void freeListData(char *** data, unsigned int length);
