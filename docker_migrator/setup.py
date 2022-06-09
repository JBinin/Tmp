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
from distutils.core import setup
setup(name='client',
      description='Tool to live-migrate docker containers',
      license='GPLv2',
      packages=['client'],
      )
setup(name='tool',
      description='Tool to live-migrate docker containers',
      license='GPLv2',
      packages=['tool'],
      )
setup(name='server',
      description='Tool to live-migrate docker containers',
      license='GPLv2',
      packages=['server'],
      )
setup(name='connect',
      description='Tool to live-migrate docker containers',
      license='GPLv2',
      packages=['connect'],
      )
