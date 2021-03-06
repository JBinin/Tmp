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
package parser

import (
	"fmt"
	"strings"
)

const (
	NoLabel = "NoLabel"
)

func LabelsFromString(str string) (map[string]string, error) {
	if len(str) == 0 {
		return nil, nil
	}
	labels := strings.Split(str, ",")
	m := make(map[string]string, len(labels))
	for _, v := range labels {
		if strings.Contains(v, "=") {
			label := strings.SplitN(v, "=", 2)
			if len(label) != 2 {
				return m, fmt.Errorf("Malformed label: %s", v)
			}
			if _, ok := m[label[0]]; ok {
				return m, fmt.Errorf("Duplicate label: %s", v)
			}
			m[label[0]] = label[1]
		} else if len(v) != 0 {
			m[v] = ""
		}
	}
	return m, nil
}

func LabelsToString(labels map[string]string) string {
	l := ""
	for k, v := range labels {
		if len(l) != 0 {
			l += ","
		}
		if len(v) != 0 {
			l += k + "=" + v
		} else if len(k) != 0 {
			l += k
		}
	}
	return l
}

func MergeLabels(old map[string]string, new map[string]string) map[string]string {
	if old == nil {
		return new
	}
	if new == nil {
		return old
	}
	m := make(map[string]string, len(old)+len(new))
	for k, v := range old {
		m[k] = v
	}
	for k, v := range new {
		m[k] = v
	}
	return m
}

func HasLabels(set map[string]string, subset map[string]string) bool {
	for k, v1 := range subset {
		if v2, ok := set[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}
