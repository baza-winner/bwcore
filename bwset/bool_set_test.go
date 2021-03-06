// Code generated by "bwsetter -type=bool"; DO NOT EDIT; bwsetter: go get -type=bool -set=Bool -test -nosort%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestBool(t *testing.T) {
	bwtesting.BwRunTests(t, BoolFrom, map[string]bwtesting.Case{"BoolFrom": {
		In: []interface{}{_BoolTestItemA, _BoolTestItemB},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, BoolFromSlice, map[string]bwtesting.Case{"BoolFromSlice": {
		In: []interface{}{[]bool{_BoolTestItemA, _BoolTestItemB}},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, BoolFromSet, map[string]bwtesting.Case{"BoolFromSet": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Bool.Copy, map[string]bwtesting.Case{"Bool.Copy": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Bool.ToSlice, map[string]bwtesting.Case{"Bool.ToSlice": {
		In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}},
		Out: []interface{}{[]bool{_BoolTestItemA}},
	}})
	bwtesting.BwRunTests(t, Bool.String, map[string]bwtesting.Case{"Bool.String": {
		In: []interface{}{Bool{_BoolTestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_BoolTestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, Bool.MarshalJSON, map[string]bwtesting.Case{"Bool.MarshalJSON": {
		In: []interface{}{Bool{_BoolTestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_BoolTestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Bool.ToSliceOfStrings, map[string]bwtesting.Case{"Bool.ToSliceOfStrings": {
		In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatBool(_BoolTestItemA)}},
	}})
	bwtesting.BwRunTests(t, Bool.Has, map[string]bwtesting.Case{
		"Bool.Has: false": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, _BoolTestItemB},
			Out: []interface{}{false},
		},
		"Bool.Has: true": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, _BoolTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool.HasAny, map[string]bwtesting.Case{
		"Bool.HasAny: empty": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"Bool.HasAny: false": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, _BoolTestItemB},
			Out: []interface{}{false},
		},
		"Bool.HasAny: true": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, _BoolTestItemA, _BoolTestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool.HasAnyOfSlice, map[string]bwtesting.Case{
		"Bool.HasAnyOfSlice: empty": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, []bool{}},
			Out: []interface{}{false},
		},
		"Bool.HasAnyOfSlice: false": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, []bool{_BoolTestItemB}},
			Out: []interface{}{false},
		},
		"Bool.HasAnyOfSlice: true": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, []bool{_BoolTestItemA, _BoolTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool.HasAnyOfSet, map[string]bwtesting.Case{
		"Bool.HasAnyOfSet: empty": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{}},
			Out: []interface{}{false},
		},
		"Bool.HasAnyOfSet: false": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{_BoolTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Bool.HasAnyOfSet: true": {
			In: []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{
				_BoolTestItemA: struct{}{},
				_BoolTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool.HasEach, map[string]bwtesting.Case{
		"Bool.HasEach: empty": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"Bool.HasEach: false": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, _BoolTestItemA, _BoolTestItemB},
			Out: []interface{}{false},
		},
		"Bool.HasEach: true": {
			In: []interface{}{Bool{
				_BoolTestItemA: struct{}{},
				_BoolTestItemB: struct{}{},
			}, _BoolTestItemA, _BoolTestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool.HasEachOfSlice, map[string]bwtesting.Case{
		"Bool.HasEachOfSlice: empty": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, []bool{}},
			Out: []interface{}{true},
		},
		"Bool.HasEachOfSlice: false": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, []bool{_BoolTestItemA, _BoolTestItemB}},
			Out: []interface{}{false},
		},
		"Bool.HasEachOfSlice: true": {
			In: []interface{}{Bool{
				_BoolTestItemA: struct{}{},
				_BoolTestItemB: struct{}{},
			}, []bool{_BoolTestItemA, _BoolTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool.HasEachOfSet, map[string]bwtesting.Case{
		"Bool.HasEachOfSet: empty": {
			In:  []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{}},
			Out: []interface{}{true},
		},
		"Bool.HasEachOfSet: false": {
			In: []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{
				_BoolTestItemA: struct{}{},
				_BoolTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Bool.HasEachOfSet: true": {
			In: []interface{}{Bool{
				_BoolTestItemA: struct{}{},
				_BoolTestItemB: struct{}{},
			}, Bool{
				_BoolTestItemA: struct{}{},
				_BoolTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Bool._AddTestHelper, map[string]bwtesting.Case{"Bool.Add": {
		In: []interface{}{Bool{_BoolTestItemA: struct{}{}}, _BoolTestItemB},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Bool._AddSliceTestHelper, map[string]bwtesting.Case{"Bool.AddSlice": {
		In: []interface{}{Bool{_BoolTestItemA: struct{}{}}, []bool{_BoolTestItemB}},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Bool._AddSetTestHelper, map[string]bwtesting.Case{"Bool.AddSet": {
		In: []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{_BoolTestItemB: struct{}{}}},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Bool._DelTestHelper, map[string]bwtesting.Case{"Bool.Del": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}, _BoolTestItemB},
		Out: []interface{}{Bool{_BoolTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Bool._DelSliceTestHelper, map[string]bwtesting.Case{"Bool.DelSlice": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}, []bool{_BoolTestItemB}},
		Out: []interface{}{Bool{_BoolTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Bool._DelSetTestHelper, map[string]bwtesting.Case{"Bool.DelSet": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}, Bool{_BoolTestItemB: struct{}{}}},
		Out: []interface{}{Bool{_BoolTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Bool.Union, map[string]bwtesting.Case{"Bool.Union": {
		In: []interface{}{Bool{_BoolTestItemA: struct{}{}}, Bool{_BoolTestItemB: struct{}{}}},
		Out: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Bool.Intersect, map[string]bwtesting.Case{"Bool.Intersect": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}, Bool{_BoolTestItemB: struct{}{}}},
		Out: []interface{}{Bool{_BoolTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Bool.Subtract, map[string]bwtesting.Case{"Bool.Subtract": {
		In: []interface{}{Bool{
			_BoolTestItemA: struct{}{},
			_BoolTestItemB: struct{}{},
		}, Bool{_BoolTestItemB: struct{}{}}},
		Out: []interface{}{Bool{_BoolTestItemA: struct{}{}}},
	}})
}
