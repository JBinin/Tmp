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
/*
Package toml provides facilities for decoding and encoding TOML configuration
files via reflection. There is also support for delaying decoding with
the Primitive type, and querying the set of keys in a TOML document with the
MetaData type.

The specification implemented: https://github.com/mojombo/toml

The sub-command github.com/BurntSushi/toml/cmd/tomlv can be used to verify
whether a file is a valid TOML document. It can also be used to print the
type of each key in a TOML document.

Testing

There are two important types of tests used for this package. The first is
contained inside '*_test.go' files and uses the standard Go unit testing
framework. These tests are primarily devoted to holistically testing the
decoder and encoder.

The second type of testing is used to verify the implementation's adherence
to the TOML specification. These tests have been factored into their own
project: https://github.com/BurntSushi/toml-test

The reason the tests are in a separate project is so that they can be used by
any implementation of TOML. Namely, it is language agnostic.
*/
package toml
