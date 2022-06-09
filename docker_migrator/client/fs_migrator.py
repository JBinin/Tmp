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
import os
import json
import subprocess as sp
from threading import Thread

rsync_log_file = "rsync.log"


class lm_docker_fs(object):
    def __init__(self, subtree_paths, cache_paths):
        self.__roots = []
        self.__cache = []
        for path in subtree_paths:
            logging.info("Initialized subtree FS hauler (%s)", path)
            self.__roots.append(path)
        for path in cache_paths:
            logging.info("Initialized subtree FS Cache (%s)", path)
            self.__cache.append(path)

        self.__thost = None

    def set_options(self, opts):
        self.__thost = opts["to"]

    def set_work_dir(self, wdir):
        self.__wdir = wdir
        self.__logf = open(os.path.join(self.__wdir, rsync_log_file), "w+")
    def __run_rsync(self, caller):

        #rsync_cmd = ["rsync","-a","-H"]
        for dir_name in self.__roots:
            dst = "%s:%s" % (self.__thost, os.path.dirname(dir_name))
            
			# First rsync might be very long. Wait for it not
			# to produce big pause between the 1st pre-dump and
			# .stop_migration

            ret = sp.call(
                ["scp", "-r", dir_name, dst],
                stdout=self.__logf, stderr=self.__logf)
            logging.info("Migrating "+dir_name+" To "+dst+" result ret :%d",ret)
            if ret != 0 and ret != 24:
                raise Exception("Migrating failed")
        if self.__cache：
            os.mkdir("/tmp/cache")
            threads = []
            i = 1
            for dir_name in self.__cache:
                thread = Thread(target=self.__tar_cache, args=(dir_name,i))
                thread.start()
                threads.append(thread)
                i = i + 1
            for thread in threads:
                thread.join()
            ret = sp.call(
                    ["tar", "-cvf", "/tmp/cache.tar", "/tmp/cache"],
                    stdout= self.__logf, stderr= self.__logf)
            logging.info("Complete Pacakge Image Layers!")
            ret = sp.call(
                    ["rm", "-rf", "/tmp/cache"],
                    stdout= self.__logf, stderr= self.__logf)
            dst = "%s:%s" % (self.__thost, "/tmp")
            ret = sp.call(
                    ["scp", "-r", "/tmp/cache.tar", dst],
                    stdout= self.__logf, stderr= self.__logf)
            logging.info("Send Image Layer to Target!")
            ret = sp.call(
                    ["rm", "-rf", "/tmp/cache.tar"],
                    stdout= self.__logf, stderr= self.__logf)
            cache_root = os.path.dirname(self.__cache[0])
            caller.decompress("/tmp/cache.tar", "/tmp/cache", cache_root)
            logging.info("Decompress Image Layer on Target!")
    def __tar_cache(self, dir_name, i):
        ret = sp.call(
                ["tar", "-cvf", "/tmp/cache/layer" + repr(i) + ".tar", dir_name],
                stdout=self.__logf, stderr=self.__logf)
        logging.info("Package Image Layer: dir_name!")
    def __run_mnt_rsync(self,worker):
        #mnt_dir = worker._ct_rootfs
        top_diff_dir = worker._topdiff_dir
        rsync_flag = True
        while rsync_flag:

            #mnt_dst = "%s:%s" % (self.__thost, os.path.dirname(mnt_dir))
            topdiff_dst = "%s:%s" % (self.__thost, os.path.dirname(top_diff_dir))

			# First rsync might be very long. Wait for it not
			# to produce big pause between the 1st pre-dump and
			# .stop_migration
            #ret = sp.call(
            #    ["rsync", "-a", mnt_dir, mnt_dst],
            #    stdout=logf, stderr=logf)
            #logging.info("rsync -a "+mnt_dir+" "+mnt_dst+" result ret :%d", ret)
            ret = sp.call(
                ["rsync", "-a", top_diff_dir, topdiff_dst],
                stdout=self.__logf, stderr=self.__logf)
            logging.info("rsync -a "+top_diff_dir+" "+topdiff_dst+" result ret :%d", ret)
            
            if ret == 0:
                rsync_flag = False
            if ret != 0 and ret != 24:
                raise Exception("Rsync failed")

    def __run_upper_dir_sync(self,worker):
        upper_dir = worker._upper_dir
        rsync_flag = True
        while rsync_flag:

            upper_dst = "%s:%s" % (self.__thost, os.path.dirname(upper_dir))

			# First rsync might be very long. Wait for it not
			# to produce big pause between the 1st pre-dump and
			# .stop_migration
            ret = sp.call(
                ["rsync", "-a", upper_dir, upper_dst],
                stdout=self.__logf, stderr=self.__logf)
            logging.info("rsync -a "+upper_dir+" "+upper_dst+" result ret :%d", ret)
            
            if ret == 0:
                rsync_flag = False
            if ret != 0 and ret != 24:
                raise Exception("Rsync failed")

    def __run_last_rsync(self,worker):
        dir_name = worker._ct_config_dir
        dump_done = False
        # Wait for dump process done!
        while not dump_done:
            config_file = open(os.path.join(dir_name,"config.v2.json"))
            try:
                config_json_str = config_file.read()
                config_json = json.loads(config_json_str)
                run_state = config_json['State']['Running']
                logging.info("Running State:"+repr(run_state))
                if not run_state:
                    dump_done = True
            finally:
                config_file.close()
        # Start last rsync process
        rsync_flag = True
        while rsync_flag:

            dst = "%s:%s" % (self.__thost, os.path.dirname(dir_name))
            ret = sp.call(
                ["rsync", "-a", dir_name, dst],
                stdout= self.__logf, stderr= self.__logf)
            logging.info("rsync -a "+dir_name+" "+dst+" result ret :%d", ret)
            if ret == 0:
                rsync_flag = False
            if ret != 0 and ret != 24:
                raise Exception("Rsync failed")
    def start_migration(self,fs_driver,caller,worker):
        logging.info("Starting FS migration")
        self.__run_rsync(caller)
        if fs_driver == "aufs":
            caller.rebuild_cache(worker.layers_dirs)
            caller.mk_mnt_dir(worker._ct_rootfs)
            caller.mk_mnt_dir(worker._ct_rootfs+"-init")
        if fs_driver == "overlay" or fs_driver == "overlay2":
            logging.info("mk merged dir!")
            caller.mk_merged_dir(worker._ct_rootfs)
        return None

    def next_iteration(self):
        return None

    def stop_migration(self,worker):
        logging.info("Doing final container config sync")
        self.__run_last_rsync(worker)
        return None
    def mnt_diff_sync(self,worker):
        logging.info("Doing mnt and top diff sync")
        self.__run_mnt_rsync(worker)
        return None
    def upper_dir_sync(self,worker):
        logging.info("Doing upper_dir sync")
        self.__run_upper_dir_sync(worker)
        return None
    def rwlayer_sync(self,worker,fs_driver):
        if fs_driver == "aufs":
            self.mnt_diff_sync(worker)
        elif fs_driver == "overlay" or fs_driver == "overlay2":
            self.upper_dir_sync(worker)
	# When rsync-ing FS inodes number will change
    def persistent_inodes(self):
        return False
