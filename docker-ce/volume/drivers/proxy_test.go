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
package volumedrivers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/go-connections/tlsconfig"
)

func TestVolumeRequestError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/VolumeDriver.Create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Cannot create volume"}`)
	})

	mux.HandleFunc("/VolumeDriver.Remove", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Cannot remove volume"}`)
	})

	mux.HandleFunc("/VolumeDriver.Mount", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Cannot mount volume"}`)
	})

	mux.HandleFunc("/VolumeDriver.Unmount", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Cannot unmount volume"}`)
	})

	mux.HandleFunc("/VolumeDriver.Path", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Unknown volume"}`)
	})

	mux.HandleFunc("/VolumeDriver.List", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Cannot list volumes"}`)
	})

	mux.HandleFunc("/VolumeDriver.Get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `{"Err": "Cannot get volume"}`)
	})

	mux.HandleFunc("/VolumeDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		http.Error(w, "error", 500)
	})

	u, _ := url.Parse(server.URL)
	client, err := plugins.NewClient("tcp://"+u.Host, &tlsconfig.Options{InsecureSkipVerify: true})
	if err != nil {
		t.Fatal(err)
	}

	driver := volumeDriverProxy{client}

	if err = driver.Create("volume", nil); err == nil {
		t.Fatal("Expected error, was nil")
	}

	if !strings.Contains(err.Error(), "Cannot create volume") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	_, err = driver.Mount("volume", "123")
	if err == nil {
		t.Fatal("Expected error, was nil")
	}

	if !strings.Contains(err.Error(), "Cannot mount volume") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	err = driver.Unmount("volume", "123")
	if err == nil {
		t.Fatal("Expected error, was nil")
	}

	if !strings.Contains(err.Error(), "Cannot unmount volume") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	err = driver.Remove("volume")
	if err == nil {
		t.Fatal("Expected error, was nil")
	}

	if !strings.Contains(err.Error(), "Cannot remove volume") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	_, err = driver.Path("volume")
	if err == nil {
		t.Fatal("Expected error, was nil")
	}

	if !strings.Contains(err.Error(), "Unknown volume") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	_, err = driver.List()
	if err == nil {
		t.Fatal("Expected error, was nil")
	}
	if !strings.Contains(err.Error(), "Cannot list volumes") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	_, err = driver.Get("volume")
	if err == nil {
		t.Fatal("Expected error, was nil")
	}
	if !strings.Contains(err.Error(), "Cannot get volume") {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	_, err = driver.Capabilities()
	if err == nil {
		t.Fatal(err)
	}
}
