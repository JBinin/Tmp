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
// Package rest provides RESTful serialization of AWS requests and responses.
package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
)

// RFC822 returns an RFC822 formatted timestamp for AWS protocols
const RFC822 = "Mon, 2 Jan 2006 15:04:05 GMT"

// Whether the byte value can be sent without escaping in AWS URLs
var noEscape [256]bool

var errValueNotSet = fmt.Errorf("value not set")

func init() {
	for i := 0; i < len(noEscape); i++ {
		// AWS expects every character except these to be escaped
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '-' ||
			i == '.' ||
			i == '_' ||
			i == '~'
	}
}

// BuildHandler is a named request handler for building rest protocol requests
var BuildHandler = request.NamedHandler{Name: "awssdk.rest.Build", Fn: Build}

// Build builds the REST component of a service request.
func Build(r *request.Request) {
	if r.ParamsFilled() {
		v := reflect.ValueOf(r.Params).Elem()
		buildLocationElements(r, v)
		buildBody(r, v)
	}
}

func buildLocationElements(r *request.Request, v reflect.Value) {
	query := r.HTTPRequest.URL.Query()

	for i := 0; i < v.NumField(); i++ {
		m := v.Field(i)
		if n := v.Type().Field(i).Name; n[0:1] == strings.ToLower(n[0:1]) {
			continue
		}

		if m.IsValid() {
			field := v.Type().Field(i)
			name := field.Tag.Get("locationName")
			if name == "" {
				name = field.Name
			}
			if m.Kind() == reflect.Ptr {
				m = m.Elem()
			}
			if !m.IsValid() {
				continue
			}

			var err error
			switch field.Tag.Get("location") {
			case "headers": // header maps
				err = buildHeaderMap(&r.HTTPRequest.Header, m, field.Tag.Get("locationName"))
			case "header":
				err = buildHeader(&r.HTTPRequest.Header, m, name)
			case "uri":
				err = buildURI(r.HTTPRequest.URL, m, name)
			case "querystring":
				err = buildQueryString(query, m, name)
			}
			r.Error = err
		}
		if r.Error != nil {
			return
		}
	}

	r.HTTPRequest.URL.RawQuery = query.Encode()
	updatePath(r.HTTPRequest.URL, r.HTTPRequest.URL.Path)
}

func buildBody(r *request.Request, v reflect.Value) {
	if field, ok := v.Type().FieldByName("_"); ok {
		if payloadName := field.Tag.Get("payload"); payloadName != "" {
			pfield, _ := v.Type().FieldByName(payloadName)
			if ptag := pfield.Tag.Get("type"); ptag != "" && ptag != "structure" {
				payload := reflect.Indirect(v.FieldByName(payloadName))
				if payload.IsValid() && payload.Interface() != nil {
					switch reader := payload.Interface().(type) {
					case io.ReadSeeker:
						r.SetReaderBody(reader)
					case []byte:
						r.SetBufferBody(reader)
					case string:
						r.SetStringBody(reader)
					default:
						r.Error = awserr.New("SerializationError",
							"failed to encode REST request",
							fmt.Errorf("unknown payload type %s", payload.Type()))
					}
				}
			}
		}
	}
}

func buildHeader(header *http.Header, v reflect.Value, name string) error {
	str, err := convertType(v)
	if err == errValueNotSet {
		return nil
	} else if err != nil {
		return awserr.New("SerializationError", "failed to encode REST request", err)
	}

	header.Add(name, str)

	return nil
}

func buildHeaderMap(header *http.Header, v reflect.Value, prefix string) error {
	for _, key := range v.MapKeys() {
		str, err := convertType(v.MapIndex(key))
		if err == errValueNotSet {
			continue
		} else if err != nil {
			return awserr.New("SerializationError", "failed to encode REST request", err)

		}

		header.Add(prefix+key.String(), str)
	}
	return nil
}

func buildURI(u *url.URL, v reflect.Value, name string) error {
	value, err := convertType(v)
	if err == errValueNotSet {
		return nil
	} else if err != nil {
		return awserr.New("SerializationError", "failed to encode REST request", err)
	}

	uri := u.Path
	uri = strings.Replace(uri, "{"+name+"}", EscapePath(value, true), -1)
	uri = strings.Replace(uri, "{"+name+"+}", EscapePath(value, false), -1)
	u.Path = uri

	return nil
}

func buildQueryString(query url.Values, v reflect.Value, name string) error {
	switch value := v.Interface().(type) {
	case []*string:
		for _, item := range value {
			query.Add(name, *item)
		}
	case map[string]*string:
		for key, item := range value {
			query.Add(key, *item)
		}
	case map[string][]*string:
		for key, items := range value {
			for _, item := range items {
				query.Add(key, *item)
			}
		}
	default:
		str, err := convertType(v)
		if err == errValueNotSet {
			return nil
		} else if err != nil {
			return awserr.New("SerializationError", "failed to encode REST request", err)
		}
		query.Set(name, str)
	}

	return nil
}

func updatePath(url *url.URL, urlPath string) {
	scheme, query := url.Scheme, url.RawQuery

	hasSlash := strings.HasSuffix(urlPath, "/")

	// clean up path
	urlPath = path.Clean(urlPath)
	if hasSlash && !strings.HasSuffix(urlPath, "/") {
		urlPath += "/"
	}

	// get formatted URL minus scheme so we can build this into Opaque
	url.Scheme, url.Path, url.RawQuery = "", "", ""
	s := url.String()
	url.Scheme = scheme
	url.RawQuery = query

	// build opaque URI
	url.Opaque = s + urlPath
}

// EscapePath escapes part of a URL path in Amazon style
func EscapePath(path string, encodeSep bool) string {
	var buf bytes.Buffer
	for i := 0; i < len(path); i++ {
		c := path[i]
		if noEscape[c] || (c == '/' && !encodeSep) {
			buf.WriteByte(c)
		} else {
			fmt.Fprintf(&buf, "%%%02X", c)
		}
	}
	return buf.String()
}

func convertType(v reflect.Value) (string, error) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return "", errValueNotSet
	}

	var str string
	switch value := v.Interface().(type) {
	case string:
		str = value
	case []byte:
		str = base64.StdEncoding.EncodeToString(value)
	case bool:
		str = strconv.FormatBool(value)
	case int64:
		str = strconv.FormatInt(value, 10)
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case time.Time:
		str = value.UTC().Format(RFC822)
	default:
		err := fmt.Errorf("Unsupported value for param %v (%s)", v.Interface(), v.Type())
		return "", err
	}
	return str, nil
}
