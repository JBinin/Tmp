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
import subprocess as sp

rsync_log_file = "rsync.log"


class lm_docker_fs(object):
	def __init__(self, subtree_paths):
		self.__roots = []
		for path in subtree_paths:
			logging.info("Initialized subtree FS hauler (%s)", path)
			self.__roots.append(path)

		self.__thost = None

	def set_options(self, opts):
		self.__thost = opts["to"]

	def set_work_dir(self, wdir):
		self.__wdir = wdir

	def __run_rsync(self):
		logf = open(os.path.join(self.__wdir, rsync_log_file), "w+")

		for dir_name in self.__roots:

			dst = "%s:%s" % (self.__thost, os.path.dirname(dir_name))

			# First rsync might be very long. Wait for it not
			# to produce big pause between the 1st pre-dump and
			# .stop_migration

			ret = sp.call(
				["rsync", "-a", dir_name, dst],
				stdout=logf, stderr=logf)
			if ret != 0:
				raise Exception("Rsync failed")

	def start_migration(self):
		logging.info("Starting FS migration")
		self.__run_rsync()
		return None

	def next_iteration(self):
		return None

	def stop_migration(self):
		logging.info("Doing final FS sync")
		self.__run_rsync()
		return None

	# When rsync-ing FS inodes number will change
	def persistent_inodes(self):
		return False