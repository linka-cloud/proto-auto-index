// Copyright 2021 Linka Cloud  All rights reserved.
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

package protoautoindex

import (
	"bytes"
	"errors"
	"os"
	"strconv"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

type Indexer interface {
	Parse(path string) error
	SetIndexes() error
	Write(path string) error
}

func New() Indexer {
	return &indexer{}
}

type indexer struct {
	node *ast.FileNode
}

func (i *indexer) Parse(path string) error {
	parser := protoparse.Parser{}
	nodes, err := parser.ParseToAST(path)
	if err != nil {
		return err
	}
	if len(nodes) != 1 {
		return errors.New("unexpected file count")
	}
	i.node = nodes[0]
	return nil
}

func (i *indexer) SetIndexes() error {
	for _, v := range i.node.Decls {
		m, ok := v.(*ast.MessageNode)
		if !ok {
			continue
		}
		index := 0
		for _, vv := range m.MessageBody.Decls {
			setIndex(vv, &index)
		}
	}
	return nil
}

func (i *indexer) Write(path string) error {
	b, err := i.print()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, b, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (i *indexer) print() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ast.Print(buf, i.node); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func setIndex(v ast.Node, index *int) {
	switch n := v.(type) {
	case interface{ FieldTag() ast.Node }:
		t, ok := n.FieldTag().(*ast.UintLiteralNode)
		if !ok {
			return
		}
		*index++
		t.Val = uint64(*index)
		l := ast.NewUintLiteralNode(uint64(*index), ast.TokenInfo{
			PosRange:          ast.PosRange{Start: *t.Start(), End: *t.End()},
			RawText:           strconv.Itoa(*index),
			LeadingComments:   t.LeadingComments(),
			TrailingComments:  t.TrailingComments(),
			LeadingWhitespace: t.LeadingWhitespace(),
		})
		*t = *l
	case *ast.OneOfNode:
		for _, e := range n.Decls {
			setIndex(e, index)
		}
	case *ast.MessageNode:
		index2 := 0
		for _, v := range n.MessageBody.Decls {
			setIndex(v, &index2)
		}
	case *ast.ReservedNode:
		for _, v := range n.Ranges {
			uiln, ok := v.StartVal.(*ast.UintLiteralNode)
			if !ok {
				continue
			}
			*index++
			oldStart, _ := v.StartVal.AsInt64()
			sv := v.StartVal
			v.StartVal = ast.NewUintLiteralNode(uint64(*index), ast.TokenInfo{
				PosRange:          ast.PosRange{Start: *uiln.Start(), End: *uiln.End()},
				RawText:           strconv.Itoa(*index),
				LeadingComments:   uiln.LeadingComments(),
				TrailingComments:  uiln.TrailingComments(),
				LeadingWhitespace: uiln.LeadingWhitespace(),
			})
			for ii := range v.Children() {
				if v.Children()[ii] == sv {
					v.Children()[ii] = v.StartVal
				}
			}
			v.StartVal.(*ast.UintLiteralNode).Val = uint64(*index)
			if v.EndVal == nil {
				continue
			}
			uiln, ok = v.EndVal.(*ast.UintLiteralNode)
			if !ok {
				continue
			}
			e, _ := v.EndVal.AsInt64()
			*index += int(e - oldStart)
			ev := v.EndVal
			v.EndVal = ast.NewUintLiteralNode(uint64(*index), ast.TokenInfo{
				PosRange:          ast.PosRange{Start: *uiln.Start(), End: *uiln.End()},
				RawText:           strconv.Itoa(*index),
				LeadingComments:   uiln.LeadingComments(),
				TrailingComments:  uiln.TrailingComments(),
				LeadingWhitespace: uiln.LeadingWhitespace(),
			})
			for ii := range v.Children() {
				if v.Children()[ii] == ev {
					v.Children()[ii] = v.EndVal
				}
			}
		}
	}
}
