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
package progress

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestOutputOnPrematureClose(t *testing.T) {
	content := []byte("TESTING")
	reader := ioutil.NopCloser(bytes.NewReader(content))
	progressChan := make(chan Progress, 10)

	pr := NewProgressReader(reader, ChanOutput(progressChan), int64(len(content)), "Test", "Read")

	part := make([]byte, 4, 4)
	_, err := io.ReadFull(pr, part)
	if err != nil {
		pr.Close()
		t.Fatal(err)
	}

drainLoop:
	for {
		select {
		case <-progressChan:
		default:
			break drainLoop
		}
	}

	pr.Close()

	select {
	case <-progressChan:
	default:
		t.Fatalf("Expected some output when closing prematurely")
	}
}

func TestCompleteSilently(t *testing.T) {
	content := []byte("TESTING")
	reader := ioutil.NopCloser(bytes.NewReader(content))
	progressChan := make(chan Progress, 10)

	pr := NewProgressReader(reader, ChanOutput(progressChan), int64(len(content)), "Test", "Read")

	out, err := ioutil.ReadAll(pr)
	if err != nil {
		pr.Close()
		t.Fatal(err)
	}
	if string(out) != "TESTING" {
		pr.Close()
		t.Fatalf("Unexpected output %q from reader", string(out))
	}

drainLoop:
	for {
		select {
		case <-progressChan:
		default:
			break drainLoop
		}
	}

	pr.Close()

	select {
	case <-progressChan:
		t.Fatalf("Should have closed silently when read is complete")
	default:
	}
}
