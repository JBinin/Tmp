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
import argparse
import client.iters

def docker_migrate_service_parse():
    """ Parse docker migrate service Command args"""
    parser = argparse.ArgumentParser("Open docker migrate service!")
    parser.add_argument("--fdrpc",type=int,required=True,help="RPC Socket File Descriptor")
    parser.add_argument("--fdmem",type=int,required=True,help="Memory Socket File Descriptor")
    parser.add_argument("--fdfs",help="Module specific definition of fs channel")
    parser.add_argument("--log-file",help="The logfile path to write log")

    return parser.parse_args()

def docker_migrate_client_parse():
    """Parse docker migrate client cmd args"""

    parser = argparse.ArgumentParser("Client to implement docker live migration")
    parser.add_argument("ct_id",help = "ID of which container to migrate!")
    parser.add_argument("--dest_ip",help = "IP of where to migrate!")
    parser.add_argument("--fdmem",type=int,required=True,help="socket fd to transmit memory data!")
    parser.add_argument("--fdrpc",type=int,required=True,help="socket fd to send rpc require data!")
    parser.add_argument("--fdfs",help="socket fd to send fs data!")
    parser.add_argument("--mode",choices=iters.MIGRATION_MODES, default=iters.MIGRATION_MODE_LIVE,help="Mode of migration")
    return parser.parse_args()