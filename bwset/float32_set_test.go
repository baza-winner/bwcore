// Code generated by "bwsetter -type=float32"; DO NOT EDIT; bwsetter: go get -type=float32 -set=Float32 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestFloat32(t *testing.T) {
	bwtesting.BwRunTests(t, Float32From, map[string]bwtesting.Case{"Float32From": {
		In: []interface{}{_Float32TestItemA, _Float32TestItemB},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32FromSlice, map[string]bwtesting.Case{"Float32FromSlice": {
		In: []interface{}{[]float32{_Float32TestItemA, _Float32TestItemB}},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32FromSet, map[string]bwtesting.Case{"Float32FromSet": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32.Copy, map[string]bwtesting.Case{"Float32.Copy": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32.ToSlice, map[string]bwtesting.Case{"Float32.ToSlice": {
		In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}},
		Out: []interface{}{[]float32{_Float32TestItemA}},
	}})
	bwtesting.BwRunTests(t, _Float32ToSliceTestHelper, map[string]bwtesting.Case{"_Float32ToSliceTestHelper": {
		In:  []interface{}{[]float32{_Float32TestItemB, _Float32TestItemA}},
		Out: []interface{}{[]float32{_Float32TestItemA, _Float32TestItemB}},
	}})
	bwtesting.BwRunTests(t, Float32.String, map[string]bwtesting.Case{"Float32.String": {
		In: []interface{}{Float32{_Float32TestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_Float32TestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, Float32.MarshalJSON, map[string]bwtesting.Case{"Float32.MarshalJSON": {
		In: []interface{}{Float32{_Float32TestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_Float32TestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Float32.ToSliceOfStrings, map[string]bwtesting.Case{"Float32.ToSliceOfStrings": {
		In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatFloat(float64(_Float32TestItemA), byte(0x66), -1, 64)}},
	}})
	bwtesting.BwRunTests(t, Float32.Has, map[string]bwtesting.Case{
		"Float32.Has: false": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, _Float32TestItemB},
			Out: []interface{}{false},
		},
		"Float32.Has: true": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, _Float32TestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32.HasAny, map[string]bwtesting.Case{
		"Float32.HasAny: empty": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"Float32.HasAny: false": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, _Float32TestItemB},
			Out: []interface{}{false},
		},
		"Float32.HasAny: true": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, _Float32TestItemA, _Float32TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32.HasAnyOfSlice, map[string]bwtesting.Case{
		"Float32.HasAnyOfSlice: empty": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, []float32{}},
			Out: []interface{}{false},
		},
		"Float32.HasAnyOfSlice: false": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, []float32{_Float32TestItemB}},
			Out: []interface{}{false},
		},
		"Float32.HasAnyOfSlice: true": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, []float32{_Float32TestItemA, _Float32TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32.HasAnyOfSet, map[string]bwtesting.Case{
		"Float32.HasAnyOfSet: empty": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{}},
			Out: []interface{}{false},
		},
		"Float32.HasAnyOfSet: false": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{_Float32TestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Float32.HasAnyOfSet: true": {
			In: []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{
				_Float32TestItemA: struct{}{},
				_Float32TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32.HasEach, map[string]bwtesting.Case{
		"Float32.HasEach: empty": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"Float32.HasEach: false": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, _Float32TestItemA, _Float32TestItemB},
			Out: []interface{}{false},
		},
		"Float32.HasEach: true": {
			In: []interface{}{Float32{
				_Float32TestItemA: struct{}{},
				_Float32TestItemB: struct{}{},
			}, _Float32TestItemA, _Float32TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32.HasEachOfSlice, map[string]bwtesting.Case{
		"Float32.HasEachOfSlice: empty": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, []float32{}},
			Out: []interface{}{true},
		},
		"Float32.HasEachOfSlice: false": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, []float32{_Float32TestItemA, _Float32TestItemB}},
			Out: []interface{}{false},
		},
		"Float32.HasEachOfSlice: true": {
			In: []interface{}{Float32{
				_Float32TestItemA: struct{}{},
				_Float32TestItemB: struct{}{},
			}, []float32{_Float32TestItemA, _Float32TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32.HasEachOfSet, map[string]bwtesting.Case{
		"Float32.HasEachOfSet: empty": {
			In:  []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{}},
			Out: []interface{}{true},
		},
		"Float32.HasEachOfSet: false": {
			In: []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{
				_Float32TestItemA: struct{}{},
				_Float32TestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Float32.HasEachOfSet: true": {
			In: []interface{}{Float32{
				_Float32TestItemA: struct{}{},
				_Float32TestItemB: struct{}{},
			}, Float32{
				_Float32TestItemA: struct{}{},
				_Float32TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float32._AddTestHelper, map[string]bwtesting.Case{"Float32.Add": {
		In: []interface{}{Float32{_Float32TestItemA: struct{}{}}, _Float32TestItemB},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32._AddSliceTestHelper, map[string]bwtesting.Case{"Float32.AddSlice": {
		In: []interface{}{Float32{_Float32TestItemA: struct{}{}}, []float32{_Float32TestItemB}},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32._AddSetTestHelper, map[string]bwtesting.Case{"Float32.AddSet": {
		In: []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{_Float32TestItemB: struct{}{}}},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32._DelTestHelper, map[string]bwtesting.Case{"Float32.Del": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}, _Float32TestItemB},
		Out: []interface{}{Float32{_Float32TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float32._DelSliceTestHelper, map[string]bwtesting.Case{"Float32.DelSlice": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}, []float32{_Float32TestItemB}},
		Out: []interface{}{Float32{_Float32TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float32._DelSetTestHelper, map[string]bwtesting.Case{"Float32.DelSet": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}, Float32{_Float32TestItemB: struct{}{}}},
		Out: []interface{}{Float32{_Float32TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float32.Union, map[string]bwtesting.Case{"Float32.Union": {
		In: []interface{}{Float32{_Float32TestItemA: struct{}{}}, Float32{_Float32TestItemB: struct{}{}}},
		Out: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float32.Intersect, map[string]bwtesting.Case{"Float32.Intersect": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}, Float32{_Float32TestItemB: struct{}{}}},
		Out: []interface{}{Float32{_Float32TestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float32.Subtract, map[string]bwtesting.Case{"Float32.Subtract": {
		In: []interface{}{Float32{
			_Float32TestItemA: struct{}{},
			_Float32TestItemB: struct{}{},
		}, Float32{_Float32TestItemB: struct{}{}}},
		Out: []interface{}{Float32{_Float32TestItemA: struct{}{}}},
	}})
}
