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
import logging
import server.rpc_server_thread

class _rpc_client_caller(object):
    def __init__(self,sk,rpc_type,rpc_funcname):
        self.rpc_sk = sk
        self.rpc_type = rpc_type
        self.rpc_funcname = rpc_funcname
    
    def __call__(self,*args):
        call = (self.rpc_type,self.rpc_funcname,args)
        raw_rpc = repr(call)
        self.rpc_sk.send(raw_rpc)
        raw_rsp = self.rpc_sk.recv(server.rpc_server_thread.rpc_sk_buf)
        resp = eval(raw_rsp)

        if resp[0] == server.rpc_server_thread.RPC_RESP:
            return resp[1]
        elif resp[0] == server.rpc_server_thread.RPC_EXC:
            logging.info("Remote rpc call failed!")
            raise Exception(resp[1])
        else:
             raise Exception("Response Error!")

class _rpc_client(object):
    def __init__(self,sk,*args):
        self._rpc_sk = sk
        c = _rpc_client_caller(self._rpc_sk,server.rpc_server_thread.RPC_CMD,"init_rpc")
        c(args)
    def __getattr__(self,funcname):
         return _rpc_client_caller(self._rpc_sk,server.rpc_server_thread.RPC_CALL,funcname)
