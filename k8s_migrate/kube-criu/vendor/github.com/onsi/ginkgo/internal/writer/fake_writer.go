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
package writer

type FakeGinkgoWriter struct {
	EventStream []string
}

func NewFake() *FakeGinkgoWriter {
	return &FakeGinkgoWriter{
		EventStream: []string{},
	}
}

func (writer *FakeGinkgoWriter) AddEvent(event string) {
	writer.EventStream = append(writer.EventStream, event)
}

func (writer *FakeGinkgoWriter) Truncate() {
	writer.EventStream = append(writer.EventStream, "TRUNCATE")
}

func (writer *FakeGinkgoWriter) DumpOut() {
	writer.EventStream = append(writer.EventStream, "DUMP")
}

func (writer *FakeGinkgoWriter) DumpOutWithHeader(header string) {
	writer.EventStream = append(writer.EventStream, "DUMP_WITH_HEADER: "+header)
}

func (writer *FakeGinkgoWriter) Bytes() []byte {
	writer.EventStream = append(writer.EventStream, "BYTES")
	return nil
}

func (writer *FakeGinkgoWriter) Write(data []byte) (n int, err error) {
	return 0, nil
}
