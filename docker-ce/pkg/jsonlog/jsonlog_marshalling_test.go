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
package jsonlog

import (
	"regexp"
	"testing"
)

func TestJSONLogMarshalJSON(t *testing.T) {
	logs := map[*JSONLog]string{
		&JSONLog{Log: `"A log line with \\"`}:           `^{\"log\":\"\\\"A log line with \\\\\\\\\\\"\",\"time\":\".{20,}\"}$`,
		&JSONLog{Log: "A log line"}:                     `^{\"log\":\"A log line\",\"time\":\".{20,}\"}$`,
		&JSONLog{Log: "A log line with \r"}:             `^{\"log\":\"A log line with \\r\",\"time\":\".{20,}\"}$`,
		&JSONLog{Log: "A log line with & < >"}:          `^{\"log\":\"A log line with \\u0026 \\u003c \\u003e\",\"time\":\".{20,}\"}$`,
		&JSONLog{Log: "A log line with utf8 : 🚀 ψ ω β"}: `^{\"log\":\"A log line with utf8 : 🚀 ψ ω β\",\"time\":\".{20,}\"}$`,
		&JSONLog{Stream: "stdout"}:                      `^{\"stream\":\"stdout\",\"time\":\".{20,}\"}$`,
		&JSONLog{}:                                      `^{\"time\":\".{20,}\"}$`,
		// These ones are a little weird
		&JSONLog{Log: "\u2028 \u2029"}:      `^{\"log\":\"\\u2028 \\u2029\",\"time\":\".{20,}\"}$`,
		&JSONLog{Log: string([]byte{0xaF})}: `^{\"log\":\"\\ufffd\",\"time\":\".{20,}\"}$`,
		&JSONLog{Log: string([]byte{0x7F})}: `^{\"log\":\"\x7f\",\"time\":\".{20,}\"}$`,
	}
	for jsonLog, expression := range logs {
		data, err := jsonLog.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		res := string(data)
		t.Logf("Result of WriteLog: %q", res)
		logRe := regexp.MustCompile(expression)
		if !logRe.MatchString(res) {
			t.Fatalf("Log line not in expected format [%v]: %q", expression, res)
		}
	}
}
