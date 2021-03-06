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
import fcntl
import socket
import errno
import os
import tarfile

class tarfile_fileobj_wrap(object):
	"""Helper class provides read/write interface for socket object

	Current helper class wrap recv/send socket methods in read/write interface.
	This functionality needed to workaround some problems of socket.makefile
	method for sockets constructed from numerical file descriptors passed
	through command line arguments.
	"""

	def __init__(self, sk):
		self.__sk = sk
		self.__nread = 0

	def read(self, size=0x10000):
		data = self.__sk.recv(size)
		self.__nread += len(data)
		return data

	def write(self, data):
		self.__sk.sendall(data)
		return len(data)

	def discard_unread_input(self):
		"""Cleanup socket after tarfile

		tarfile module always align data on source side according to RECORDSIZE
		constant, but it don't read aligning bytes on target side in some cases
		depending on received buffer size. Read aligning manually and discard.
		"""

		remainder = self.__nread % tarfile.RECORDSIZE
		if remainder > 0:
			self.__sk.recv(tarfile.RECORDSIZE - remainder, socket.MSG_WAITALL)
			self.__nread += tarfile.RECORDSIZE - remainder

def log_uncaught_exception(type, value, traceback):
	logging.error(value, exc_info=(type, value, traceback))


def log_header():
	OFFSET_LINES_COUNT = 3
	for i in range(OFFSET_LINES_COUNT):
		logging.info("")

def set_cloexec(sk):
	flags = fcntl.fcntl(sk, fcntl.F_GETFD)
	fcntl.fcntl(sk, fcntl.F_SETFD, flags | fcntl.FD_CLOEXEC)
def makedirs(dirpath):
	try:
		os.makedirs(dirpath)
	except OSError as er:
		if er.errno == errno.EEXIST and os.path.isdir(dirpath):
			pass
		else:
			raise
