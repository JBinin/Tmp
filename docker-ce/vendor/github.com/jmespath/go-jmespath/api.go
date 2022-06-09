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
package jmespath

import "strconv"

// JmesPath is the epresentation of a compiled JMES path query. A JmesPath is
// safe for concurrent use by multiple goroutines.
type JMESPath struct {
	ast  ASTNode
	intr *treeInterpreter
}

// Compile parses a JMESPath expression and returns, if successful, a JMESPath
// object that can be used to match against data.
func Compile(expression string) (*JMESPath, error) {
	parser := NewParser()
	ast, err := parser.Parse(expression)
	if err != nil {
		return nil, err
	}
	jmespath := &JMESPath{ast: ast, intr: newInterpreter()}
	return jmespath, nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled
// JMESPaths.
func MustCompile(expression string) *JMESPath {
	jmespath, err := Compile(expression)
	if err != nil {
		panic(`jmespath: Compile(` + strconv.Quote(expression) + `): ` + err.Error())
	}
	return jmespath
}

// Search evaluates a JMESPath expression against input data and returns the result.
func (jp *JMESPath) Search(data interface{}) (interface{}, error) {
	return jp.intr.Execute(jp.ast, data)
}

// Search evaluates a JMESPath expression against input data and returns the result.
func Search(expression string, data interface{}) (interface{}, error) {
	intr := newInterpreter()
	parser := NewParser()
	ast, err := parser.Parse(expression)
	if err != nil {
		return nil, err
	}
	return intr.Execute(ast, data)
}
