package bwmap_test

// func TestCropMap(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		func(m interface{}, crop ...interface{}) (result interface{}) {
// 			result = deepcopy.Copy(m)
// 			bwmap.CropMap(result, crop...)
// 			return
// 		},
// 		map[string]bwtesting.Case{
// 			"string": {
// 				In: []interface{}{
// 					map[string]interface{}{
// 						"some": "thing",
// 						"good": "is not bad",
// 					},
// 					// []interface{}{
// 					`some`,
// 					// },
// 				},
// 				Out: []interface{}{
// 					map[string]interface{}{
// 						"some": "thing",
// 					},
// 				},
// 			},
// 			"[]string": {
// 				In: []interface{}{
// 					map[string]interface{}{
// 						"A": 1,
// 						"B": 2,
// 						"C": 3,
// 						"D": 4,
// 					},
// 					// []interface{}{},
// 					[]string{"B", "C"},
// 				},
// 				Out: []interface{}{
// 					map[string]interface{}{
// 						"B": 2,
// 						"C": 3,
// 					},
// 				},
// 			},
// 			"map[string]interface{}": {
// 				In: []interface{}{
// 					map[string]int{
// 						"A": 1,
// 						"B": 2,
// 						"C": 3,
// 						"D": 4,
// 					},
// 					// []interface{}{
// 					map[string]interface{}{
// 						"A": struct{}{},
// 						"D": struct{}{},
// 					},
// 					// },
// 				},
// 				Out: []interface{}{
// 					map[string]int{
// 						"A": 1,
// 						"D": 4,
// 					},
// 				},
// 			},
// 			"mixed": {
// 				In: []interface{}{
// 					map[string]int{
// 						"A": 1,
// 						"B": 2,
// 						"C": 3,
// 						"D": 4,
// 						"E": 5,
// 						"F": 6,
// 						"G": 7,
// 						"H": 8,
// 					},
// 					// []interface{}{
// 					`A`,
// 					[]string{"C", "D"},
// 					map[string]struct{}{
// 						"F": struct{}{},
// 						"G": struct{}{},
// 					},
// 					// },
// 				},
// 				Out: []interface{}{
// 					map[string]int{
// 						"A": 1,
// 						"C": 3,
// 						"D": 4,
// 						"F": 6,
// 						"G": 7,
// 					},
// 				},
// 			},
// 		},
// 	)
// }
