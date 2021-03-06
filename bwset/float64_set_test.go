// Code generated by "bwsetter -type=float64"; DO NOT EDIT; bwsetter: go get -type=float64 -set=Float64 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestFloat64(t *testing.T) {
	bwtesting.BwRunTests(t, Float64From, map[string]bwtesting.Case{"Float64From": {
		In: []interface{}{_Float64TestItemA, _Float64TestItemB},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64FromSlice, map[string]bwtesting.Case{"Float64FromSlice": {
		In: []interface{}{[]float64{_Float64TestItemA, _Float64TestItemB}},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64FromSet, map[string]bwtesting.Case{"Float64FromSet": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64.Copy, map[string]bwtesting.Case{"Float64.Copy": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64.ToSlice, map[string]bwtesting.Case{"Float64.ToSlice": {
		In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}},
		Out: []interface{}{[]float64{_Float64TestItemA}},
	}})
	bwtesting.BwRunTests(t, _Float64ToSliceTestHelper, map[string]bwtesting.Case{"_Float64ToSliceTestHelper": {
		In:  []interface{}{[]float64{_Float64TestItemB, _Float64TestItemA}},
		Out: []interface{}{[]float64{_Float64TestItemA, _Float64TestItemB}},
	}})
	bwtesting.BwRunTests(t, Float64.String, map[string]bwtesting.Case{"Float64.String": {
		In: []interface{}{Float64{_Float64TestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_Float64TestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, Float64.MarshalJSON, map[string]bwtesting.Case{"Float64.MarshalJSON": {
		In: []interface{}{Float64{_Float64TestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_Float64TestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Float64.ToSliceOfStrings, map[string]bwtesting.Case{"Float64.ToSliceOfStrings": {
		In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatFloat(float64(_Float64TestItemA), byte(0x66), -1, 64)}},
	}})
	bwtesting.BwRunTests(t, Float64.Has, map[string]bwtesting.Case{
		"Float64.Has: false": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, _Float64TestItemB},
			Out: []interface{}{false},
		},
		"Float64.Has: true": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, _Float64TestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64.HasAny, map[string]bwtesting.Case{
		"Float64.HasAny: empty": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"Float64.HasAny: false": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, _Float64TestItemB},
			Out: []interface{}{false},
		},
		"Float64.HasAny: true": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, _Float64TestItemA, _Float64TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64.HasAnyOfSlice, map[string]bwtesting.Case{
		"Float64.HasAnyOfSlice: empty": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, []float64{}},
			Out: []interface{}{false},
		},
		"Float64.HasAnyOfSlice: false": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, []float64{_Float64TestItemB}},
			Out: []interface{}{false},
		},
		"Float64.HasAnyOfSlice: true": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, []float64{_Float64TestItemA, _Float64TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64.HasAnyOfSet, map[string]bwtesting.Case{
		"Float64.HasAnyOfSet: empty": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{}},
			Out: []interface{}{false},
		},
		"Float64.HasAnyOfSet: false": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{_Float64TestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Float64.HasAnyOfSet: true": {
			In: []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{
				_Float64TestItemA: struct{}{},
				_Float64TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64.HasEach, map[string]bwtesting.Case{
		"Float64.HasEach: empty": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"Float64.HasEach: false": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, _Float64TestItemA, _Float64TestItemB},
			Out: []interface{}{false},
		},
		"Float64.HasEach: true": {
			In: []interface{}{Float64{
				_Float64TestItemA: struct{}{},
				_Float64TestItemB: struct{}{},
			}, _Float64TestItemA, _Float64TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64.HasEachOfSlice, map[string]bwtesting.Case{
		"Float64.HasEachOfSlice: empty": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, []float64{}},
			Out: []interface{}{true},
		},
		"Float64.HasEachOfSlice: false": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, []float64{_Float64TestItemA, _Float64TestItemB}},
			Out: []interface{}{false},
		},
		"Float64.HasEachOfSlice: true": {
			In: []interface{}{Float64{
				_Float64TestItemA: struct{}{},
				_Float64TestItemB: struct{}{},
			}, []float64{_Float64TestItemA, _Float64TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64.HasEachOfSet, map[string]bwtesting.Case{
		"Float64.HasEachOfSet: empty": {
			In:  []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{}},
			Out: []interface{}{true},
		},
		"Float64.HasEachOfSet: false": {
			In: []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{
				_Float64TestItemA: struct{}{},
				_Float64TestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Float64.HasEachOfSet: true": {
			In: []interface{}{Float64{
				_Float64TestItemA: struct{}{},
				_Float64TestItemB: struct{}{},
			}, Float64{
				_Float64TestItemA: struct{}{},
				_Float64TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64._AddTestHelper, map[string]bwtesting.Case{"Float64.Add": {
		In: []interface{}{Float64{_Float64TestItemA: struct{}{}}, _Float64TestItemB},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64._AddSliceTestHelper, map[string]bwtesting.Case{"Float64.AddSlice": {
		In: []interface{}{Float64{_Float64TestItemA: struct{}{}}, []float64{_Float64TestItemB}},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64._AddSetTestHelper, map[string]bwtesting.Case{"Float64.AddSet": {
		In: []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{_Float64TestItemB: struct{}{}}},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64._DelTestHelper, map[string]bwtesting.Case{"Float64.Del": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}, _Float64TestItemB},
		Out: []interface{}{Float64{_Float64TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64._DelSliceTestHelper, map[string]bwtesting.Case{"Float64.DelSlice": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}, []float64{_Float64TestItemB}},
		Out: []interface{}{Float64{_Float64TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64._DelSetTestHelper, map[string]bwtesting.Case{"Float64.DelSet": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}, Float64{_Float64TestItemB: struct{}{}}},
		Out: []interface{}{Float64{_Float64TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64.Union, map[string]bwtesting.Case{"Float64.Union": {
		In: []interface{}{Float64{_Float64TestItemA: struct{}{}}, Float64{_Float64TestItemB: struct{}{}}},
		Out: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64.Intersect, map[string]bwtesting.Case{"Float64.Intersect": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}, Float64{_Float64TestItemB: struct{}{}}},
		Out: []interface{}{Float64{_Float64TestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64.Subtract, map[string]bwtesting.Case{"Float64.Subtract": {
		In: []interface{}{Float64{
			_Float64TestItemA: struct{}{},
			_Float64TestItemB: struct{}{},
		}, Float64{_Float64TestItemB: struct{}{}}},
		Out: []interface{}{Float64{_Float64TestItemA: struct{}{}}},
	}})
}
