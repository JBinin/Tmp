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
#include <Security/Security.h>

struct Server {
  SecProtocolType proto;
  char *host;
  char *path;
  unsigned int port;
};

char *keychain_add(struct Server *server, char *username, char *secret);
char *keychain_get(struct Server *server, unsigned int *username_l, char **username, unsigned int *secret_l, char **secret);
char *keychain_delete(struct Server *server);
char *keychain_list(char *** data, char *** accts, unsigned int *list_l);
void freeListData(char *** data, unsigned int length);