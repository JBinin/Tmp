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
#
#Help to operate on target host
#
import os
import tool.criu_req
import tool.criu_api
import logging
import client.iters
import client.img_migrator
import client.docker_migrate_worker
import copy
import subprocess as sp

class rpc_migrate_service(object):
    def __init__(self, connection):
        self.connection = connection
        self._migrate_worker = None
        self.criu_connection = None
        self.dump_iter_index = 0
        self.img = None
        self.cacheids_remap = []
        self.__mode = client.iters.MIGRATION_MODE_LIVE
        self.restored = False
    def on_connect(self):
        logging.info("Rpc Service Connected!")

    def dis_connect(self):
        logging.info("Rpc Service Disconnected!")
        if self.criu_connection:
            self.criu_connection.close()
        if self._migrate_worker:
            if client.iters.is_live_mode(self.__mode):
                self._migrate_worker.umount()
        if self.img:
            logging.info("Closing images")
        if not self.restored:
            self.img.save_images()
        self.img.close()

    def rpc_setup(self, ct_id, mode):
        self.mode = mode
        self._migrate_worker = client.docker_migrate_worker.docker_lm_worker(ct_id)
        self._migrate_worker.init_dst()
        if client.iters.is_live_mode(self.__mode):
            self.img = client.img_migrator.lm_docker_img("rst")
            self.criu_connection = tool.criu_api.criu_conn(self.connection.fdmem)

    def rpc_set_options(self, opts):
        self._migrate_worker.set_options(opts)
        if self.criu_connection:
            self.criu_connection.set_options(opts)

        if self.img:
            self.img.set_options(opts)

    def rpc_start_accept_images(self, dir_id):
        self.img.start_accept_images(dir_id, self._migrate_worker, self.connection.fdmem)

    def rpc_stop_accept_images(self):
        self.img.stop_accept_images()

    def rpc_check_cpuinfo(self):
        logging.info("Checking cpuinfo")
        req = tool.criu_req.make_cpuinfo_check_req(self.img)
        resp = self.criu_connection.send_req(req)
        logging.info("\t`- %s", resp.success)
        return resp.success

    def rpc_check_criu_version(self, source_version):
        logging.info("Checking criu version")
        target_version = criu_api.get_criu_version()
        if not target_version:
            logging.info("\t`- Can't get criu version")
            return False
        lsource_version = distutils.version.LooseVersion(source_version)
        ltarget_version = distutils.version.LooseVersion(target_version)
        result = lsource_version <= ltarget_version
        logging.info("\t`- %s -> %s", source_version, target_version)
        logging.info("\t`- %s", result)
        return result

    def rpc_start_iter(self, need_page_server):
        self.dump_iter_index += 1
        self.img.new_image_dir()
        if need_page_server:
            self.start_page_server()

    def rpc_end_iter(self):
        pass
    def rpc_restore_time(self):
        stats = tool.criu_api.criu_get_rstats(self.img)
        return stats.restore_time
    def rpc_restore_from_images(self, ctid, ck_dir):
        logging.info("Restoring from images")
        self._migrate_worker.put_meta_images(self.img.image_dir(), ctid, ck_dir)
        ret = self._migrate_worker.final_restore(self.img, self.criu_connection, ck_dir)
        if ret != 0:
            return ret
        else:
            logging.info("Restore succeeded")
            self.restored = True
            return ret
    def rpc_get_image_dir(self):
        return self.img.image_dir()

    def rpc_mk_merged_dir(self,merged_parent_dir):
        os.chdir(merged_parent_dir)
        if not os.path.exists(os.path.join(merged_parent_dir,"merged")):
            os.mkdir("merged")
    def rpc_get_cacheids_sended(self,cache_dirs):
        for cacheid in self.cacheids_remap:
            cache_dirs.pop()
        return cache_dirs
    def rpc_rebuild_layer(self,layer_dirs):
        if self.cacheids_remap:
            for layer_dir in layer_dirs:
                parent_ids = []
                try:
                    relation_file = open(layer_dir)
                    while 1:
                        parent_id= relation_file.readline()
                        if not parent_id:
                            break
                        parent_id = parent_id.strip()
                        logging.info("parent_id:%s",parent_id)
                        parent_ids.append(parent_id+"\n")
                finally:
                    relation_file.close()
                for cache_id in self.cacheids_remap:
                    parent_ids.pop()
                parent_ids.extend(self.cacheids_remap)
                os.remove(layer_dir)
                try:
                    new_layer_file = open(layer_dir,'w')
                    new_layer_file.writelines(parent_ids)
                finally:
                    new_layer_file.close()

    def rpc_check_layers(self,chain_dirs):
        chain_dirs_sended = copy.deepcopy(chain_dirs)
        chain_dirs.reverse()
        for chain_dir in chain_dirs:
            logging.info("check chain_dirs : %s",chain_dir)
            if os.path.exists(chain_dir):
                chain_dirs_sended.pop()
                cache_ids_file = open(chain_dir+"/cache-id")
                cacheid = cache_ids_file.readline()
                self.cacheids_remap.append(cacheid+"\n")
            else:
                break
        self.cacheids_remap.reverse()
        logging.info("cacheids_remap : %s",self.cacheids_remap)
        return chain_dirs_sended
    def rpc_decompress(self, tar_path, tmp_path, cache_root):
        logf = open("/tmp/tar.log","w+")
        ret = sp.call(
                ["tar", "-xzvf", tar_path, "-C", "/"],
                stdout=logf, stderr=logf)
        ret = sp.call(
                    ["rm", "-rf", "/tmp/cache.tar"],
                    stdout= self.__logf, stderr= self.__logf)
        layers = os.listdir(tmp_path)
        os.chdir(tmp_path)
        for layer in layers:
            ret = sp.call(
                ["tar", "-xzvf", layer, "-C", "/"],
                stdout=logf, stderr=logf)
        ret = sp.call(
                    ["rm", "-rf", "/tmp/cache"],
                    stdout= self.__logf, stderr= self.__logf)
        

