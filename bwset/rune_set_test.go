// Code generated by "bwsetter -type=rune"; DO NOT EDIT; bwsetter: go get -type=rune -set=Rune -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"testing"
)

func TestRune(t *testing.T) {
	bwtesting.BwRunTests(t, RuneFrom, map[string]bwtesting.Case{"RuneFrom": {
		In: []interface{}{_RuneTestItemA, _RuneTestItemB},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, RuneFromSlice, map[string]bwtesting.Case{"RuneFromSlice": {
		In: []interface{}{[]rune{_RuneTestItemA, _RuneTestItemB}},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, RuneFromSet, map[string]bwtesting.Case{"RuneFromSet": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Rune.Copy, map[string]bwtesting.Case{"Rune.Copy": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Rune.ToSlice, map[string]bwtesting.Case{"Rune.ToSlice": {
		In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}},
		Out: []interface{}{[]rune{_RuneTestItemA}},
	}})
	bwtesting.BwRunTests(t, _RuneToSliceTestHelper, map[string]bwtesting.Case{"_RuneToSliceTestHelper": {
		In:  []interface{}{[]rune{_RuneTestItemB, _RuneTestItemA}},
		Out: []interface{}{[]rune{_RuneTestItemA, _RuneTestItemB}},
	}})
	bwtesting.BwRunTests(t, Rune.String, map[string]bwtesting.Case{"Rune.String": {
		In: []interface{}{Rune{_RuneTestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_RuneTestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, Rune.MarshalJSON, map[string]bwtesting.Case{"Rune.MarshalJSON": {
		In: []interface{}{Rune{_RuneTestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_RuneTestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Rune.ToSliceOfStrings, map[string]bwtesting.Case{"Rune.ToSliceOfStrings": {
		In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}},
		Out: []interface{}{[]string{string(_RuneTestItemA)}},
	}})
	bwtesting.BwRunTests(t, Rune.Has, map[string]bwtesting.Case{
		"Rune.Has: false": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, _RuneTestItemB},
			Out: []interface{}{false},
		},
		"Rune.Has: true": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, _RuneTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune.HasAny, map[string]bwtesting.Case{
		"Rune.HasAny: empty": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"Rune.HasAny: false": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, _RuneTestItemB},
			Out: []interface{}{false},
		},
		"Rune.HasAny: true": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, _RuneTestItemA, _RuneTestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune.HasAnyOfSlice, map[string]bwtesting.Case{
		"Rune.HasAnyOfSlice: empty": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, []rune{}},
			Out: []interface{}{false},
		},
		"Rune.HasAnyOfSlice: false": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, []rune{_RuneTestItemB}},
			Out: []interface{}{false},
		},
		"Rune.HasAnyOfSlice: true": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, []rune{_RuneTestItemA, _RuneTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune.HasAnyOfSet, map[string]bwtesting.Case{
		"Rune.HasAnyOfSet: empty": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{}},
			Out: []interface{}{false},
		},
		"Rune.HasAnyOfSet: false": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{_RuneTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Rune.HasAnyOfSet: true": {
			In: []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{
				_RuneTestItemA: struct{}{},
				_RuneTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune.HasEach, map[string]bwtesting.Case{
		"Rune.HasEach: empty": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"Rune.HasEach: false": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, _RuneTestItemA, _RuneTestItemB},
			Out: []interface{}{false},
		},
		"Rune.HasEach: true": {
			In: []interface{}{Rune{
				_RuneTestItemA: struct{}{},
				_RuneTestItemB: struct{}{},
			}, _RuneTestItemA, _RuneTestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune.HasEachOfSlice, map[string]bwtesting.Case{
		"Rune.HasEachOfSlice: empty": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, []rune{}},
			Out: []interface{}{true},
		},
		"Rune.HasEachOfSlice: false": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, []rune{_RuneTestItemA, _RuneTestItemB}},
			Out: []interface{}{false},
		},
		"Rune.HasEachOfSlice: true": {
			In: []interface{}{Rune{
				_RuneTestItemA: struct{}{},
				_RuneTestItemB: struct{}{},
			}, []rune{_RuneTestItemA, _RuneTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune.HasEachOfSet, map[string]bwtesting.Case{
		"Rune.HasEachOfSet: empty": {
			In:  []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{}},
			Out: []interface{}{true},
		},
		"Rune.HasEachOfSet: false": {
			In: []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{
				_RuneTestItemA: struct{}{},
				_RuneTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Rune.HasEachOfSet: true": {
			In: []interface{}{Rune{
				_RuneTestItemA: struct{}{},
				_RuneTestItemB: struct{}{},
			}, Rune{
				_RuneTestItemA: struct{}{},
				_RuneTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Rune._AddTestHelper, map[string]bwtesting.Case{"Rune.Add": {
		In: []interface{}{Rune{_RuneTestItemA: struct{}{}}, _RuneTestItemB},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Rune._AddSliceTestHelper, map[string]bwtesting.Case{"Rune.AddSlice": {
		In: []interface{}{Rune{_RuneTestItemA: struct{}{}}, []rune{_RuneTestItemB}},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Rune._AddSetTestHelper, map[string]bwtesting.Case{"Rune.AddSet": {
		In: []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{_RuneTestItemB: struct{}{}}},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Rune._DelTestHelper, map[string]bwtesting.Case{"Rune.Del": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}, _RuneTestItemB},
		Out: []interface{}{Rune{_RuneTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Rune._DelSliceTestHelper, map[string]bwtesting.Case{"Rune.DelSlice": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}, []rune{_RuneTestItemB}},
		Out: []interface{}{Rune{_RuneTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Rune._DelSetTestHelper, map[string]bwtesting.Case{"Rune.DelSet": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}, Rune{_RuneTestItemB: struct{}{}}},
		Out: []interface{}{Rune{_RuneTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Rune.Union, map[string]bwtesting.Case{"Rune.Union": {
		In: []interface{}{Rune{_RuneTestItemA: struct{}{}}, Rune{_RuneTestItemB: struct{}{}}},
		Out: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Rune.Intersect, map[string]bwtesting.Case{"Rune.Intersect": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}, Rune{_RuneTestItemB: struct{}{}}},
		Out: []interface{}{Rune{_RuneTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Rune.Subtract, map[string]bwtesting.Case{"Rune.Subtract": {
		In: []interface{}{Rune{
			_RuneTestItemA: struct{}{},
			_RuneTestItemB: struct{}{},
		}, Rune{_RuneTestItemB: struct{}{}}},
		Out: []interface{}{Rune{_RuneTestItemA: struct{}{}}},
	}})
}
