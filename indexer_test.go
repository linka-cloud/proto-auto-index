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
	"testing"

	assert2 "github.com/stretchr/testify/assert"
	require2 "github.com/stretchr/testify/require"
)

var w = `syntax = "proto3";

package protoautoindex.test;


import "file/not/found.proto";

message Test {
  // comment here
  string one = 1;
  int32 two = 2;
  // comment here too...
  map<string,string> three = 3;
  repeated uint32 four = 4;
  // we could use a better name...
  oneof oneof {
    // there too
    string five = 5;
    string six = 6;
  }
  message Sub {
    string one = 1;
    int32 two = 2;
    map<string,string> three = 3;
    repeated uint32 four = 4;
    oneof oneof {
      // there too
      string five = 5;
      string six = 6;
    }
  }
  Sub seven = 7;
}
`

func TestIndexer(t *testing.T) {
	require := require2.New(t)
	assert := assert2.New(t)
	i := indexer{}
	require.NoError(i.Parse("test.proto"))
	require.NoError(i.SetIndexes())
	b, err := i.print()
	require.NoError(err)
	assert.Equal(w, string(b))
}
