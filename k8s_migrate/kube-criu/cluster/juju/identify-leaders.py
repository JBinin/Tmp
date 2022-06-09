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
#!/usr/bin/env python

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from subprocess import check_output
import yaml


cmd = ['juju', 'run', '--application', 'kubernetes', '--format=yaml', 'is-leader']
out = check_output(cmd)
try:
    parsed_output = yaml.safe_load(out)
    for unit in parsed_output:
        standard_out = unit['Stdout'].rstrip()
        unit_id = unit['UnitId']
        if 'True' in standard_out:
            print(unit_id)
except:
    pass
