// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"bytes"
	"fmt"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/chunk"
	"sort"
	"strings"
)

// ExplainInfo implements the Expression interface.
func (expr *ScalarFunction) ExplainInfo() string {
	return expr.explainInfo(false)
}

func (expr *ScalarFunction) explainInfo(normalized bool) string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%s(", expr.FuncName.L)
	switch expr.FuncName.L {
	case ast.Cast:
		for _, arg := range expr.GetArgs() {
			if normalized {
				buffer.WriteString(arg.ExplainNormalizedInfo())
			} else {
				buffer.WriteString(arg.ExplainInfo())
			}
			buffer.WriteString(", ")
			buffer.WriteString(expr.RetType.String())
		}
	default:
		for i, arg := range expr.GetArgs() {
			if normalized {
				buffer.WriteString(arg.ExplainNormalizedInfo())
			} else {
				buffer.WriteString(arg.ExplainInfo())
			}
			if i+1 < len(expr.GetArgs()) {
				buffer.WriteString(", ")
			}
		}
	}
	buffer.WriteString(")")
	return buffer.String()
}

func (expr *ScalarFunction) Restore(t *model.TableInfo) string {
	return expr.restore(t)
}

func (expr *ScalarFunction) restore(t *model.TableInfo) string {
	var buffer bytes.Buffer
	switch expr.FuncName.L {
	case ast.LT:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " < ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.GT:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " > ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.LE:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " <= ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.GE:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " >= ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.EQ:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " = ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.NE:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " != ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.LogicOr:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " or ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.Plus:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " + ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.Minus:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " - ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.Mul:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " * ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	case ast.Div:
		fmt.Fprint(&buffer, "(")
		buffer.WriteString(expr.GetArgs()[0].Restore(t))
		fmt.Fprint(&buffer, " / ")
		buffer.WriteString(expr.GetArgs()[1].Restore(t))
		fmt.Fprint(&buffer, ")")
	}
	return buffer.String()
}

func (col *Column) Restore(t *model.TableInfo) string {
	for _, c := range t.Columns {
		if c.ID == col.ID {
			return fmt.Sprintf(`"%v"`, c.Name.L)
		}
	}
	return ""
}

func (expr *Constant) Restore(t *model.TableInfo) string {
	return expr.restore4S3()
}

// ExplainNormalizedInfo implements the Expression interface.
func (expr *ScalarFunction) ExplainNormalizedInfo() string {
	return expr.explainInfo(true)
}

// ExplainInfo implements the Expression interface.
func (col *Column) ExplainInfo() string {
	return col.String()
}

// ExplainNormalizedInfo implements the Expression interface.
func (col *Column) ExplainNormalizedInfo() string {
	if col.OrigName != "" {
		return col.OrigName
	}
	return "?"
}

func (expr *Constant) restore4S3() string {
	dt, err := expr.Eval(chunk.Row{})
	if err != nil {
		return "not recognized const vanue"
	}
	return expr.format4S3(dt)
}

func (expr *Constant) format4S3(dt types.Datum) string {
	switch dt.Kind() {
	case types.KindNull:
		return "NULL"
	case types.KindString, types.KindBytes, types.KindMysqlEnum, types.KindMysqlSet,
		types.KindMysqlJSON, types.KindBinaryLiteral, types.KindMysqlBit:
		return fmt.Sprintf("'%v'", dt.GetValue())
	case types.KindMysqlTime:
		t := dt.Clone().GetMysqlTime()
		t.SetFsp(0)
		return fmt.Sprintf("'%v'", t)
	}
	return fmt.Sprintf("%v", dt.GetValue())
}

// ExplainInfo implements the Expression interface.
func (expr *Constant) ExplainInfo() string {
	dt, err := expr.Eval(chunk.Row{})
	if err != nil {
		return "not recognized const vanue"
	}
	return expr.format(dt)
}

// ExplainNormalizedInfo implements the Expression interface.
func (expr *Constant) ExplainNormalizedInfo() string {
	return "?"
}

func (expr *Constant) format(dt types.Datum) string {
	switch dt.Kind() {
	case types.KindNull:
		return "NULL"
	case types.KindString, types.KindBytes, types.KindMysqlEnum, types.KindMysqlSet,
		types.KindMysqlJSON, types.KindBinaryLiteral, types.KindMysqlBit:
		return fmt.Sprintf("\"%v\"", dt.GetValue())
	}
	return fmt.Sprintf("%v", dt.GetValue())
}

// ExplainExpressionList generates explain information for a list of expressions.
func ExplainExpressionList(exprs []Expression, schema *Schema) string {
	builder := &strings.Builder{}
	for i, expr := range exprs {
		switch expr.(type) {
		case *Column, *CorrelatedColumn:
			builder.WriteString(expr.String())
		case *Constant:
			v := expr.String()
			length := 64
			if len(v) < length {
				builder.WriteString(v)
			} else {
				builder.WriteString(fmt.Sprintf("%s(len:%d)", v[:length], len(v)))
			}
			builder.WriteString("->")
			builder.WriteString(schema.Columns[i].String())
		default:
			builder.WriteString(expr.String())
			builder.WriteString("->")
			builder.WriteString(schema.Columns[i].String())
		}
		if i+1 < len(exprs) {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}

// SortedExplainExpressionList generates explain information for a list of expressions in order.
// In some scenarios, the expr's order may not be stable when executing multiple times.
// So we add a sort to make its explain result stable.
func SortedExplainExpressionList(exprs []Expression) []byte {
	return sortedExplainExpressionList(exprs, false)
}

func sortedExplainExpressionList(exprs []Expression, normalized bool) []byte {
	buffer := bytes.NewBufferString("")
	exprInfos := make([]string, 0, len(exprs))
	for _, expr := range exprs {
		if normalized {
			exprInfos = append(exprInfos, expr.ExplainNormalizedInfo())
		} else {
			exprInfos = append(exprInfos, expr.ExplainInfo())
		}
	}
	sort.Strings(exprInfos)
	for i, info := range exprInfos {
		buffer.WriteString(info)
		if i+1 < len(exprInfos) {
			buffer.WriteString(", ")
		}
	}
	return buffer.Bytes()
}

// SortedExplainNormalizedExpressionList is same like SortedExplainExpressionList, but use for generating normalized information.
func SortedExplainNormalizedExpressionList(exprs []Expression) []byte {
	return sortedExplainExpressionList(exprs, true)
}

// SortedExplainNormalizedScalarFuncList is same like SortedExplainExpressionList, but use for generating normalized information.
func SortedExplainNormalizedScalarFuncList(exprs []*ScalarFunction) []byte {
	expressions := make([]Expression, len(exprs))
	for i := range exprs {
		expressions[i] = exprs[i]
	}
	return sortedExplainExpressionList(expressions, true)
}

// ExplainColumnList generates explain information for a list of columns.
func ExplainColumnList(cols []*Column) []byte {
	buffer := bytes.NewBufferString("")
	for i, col := range cols {
		buffer.WriteString(col.ExplainInfo())
		if i+1 < len(cols) {
			buffer.WriteString(", ")
		}
	}
	return buffer.Bytes()
}
