package isoline

type Idx struct {
	Idx_i int
	Idx_j int
}

func NewIdx(args [2]int) *Idx {
	idx := Idx{}
	idx.Idx_i = args[0]
	idx.Idx_j = args[1]
	return &idx
}

type LatLon struct {
	Lat float32
	Lon float32
}

type EdgeInfo struct {
	// 与该边相邻的边,最多两个
	Neighbours *[]Idx
	// 等值线穿过该边的位置
	Location *LatLon
	Passed   bool
}

// func (edgeInfo *EdgeInfo) print() {
// 	fmt.Print(*edgeInfo.Neighbours)
// 	fmt.Print(*edgeInfo.Location)
// 	fmt.Print(edgeInfo.Passed)
// 	fmt.Println()
// }

func BuildShape(x int, y int, threshold float32, edges *[][]EdgeInfo, lats *[]float32, lons *[]float32, vals *[][]float32) {
	// 构造边的信息

	idxes := [4][2]int{{2 * x, y}, {2*x + 1, y + 1}, {2*x + 2, y}, {2*x + 1, y}}
	idxes_ := [4][2]int{{x, y}, {x, y + 1}, {x + 1, y + 1}, {x + 1, y}}
	v := [4]float32{}
	for i := range idxes_ {
		v[i] = (*vals)[idxes_[i][0]][idxes_[i][1]]
	}
	judge := [4]bool{}
	total := 0
	for i := range judge {
		judge[i] = v[i] > threshold
		if judge[i] {
			total += 1
		}
	}
	lat_init, lon_init := (*lats)[x], (*lons)[y]
	lat_diff, lon_diff := (*lats)[x+1]-(*lats)[x], (*lons)[y+1]-(*lons)[y]

	rate_cal := [4][2]float32{{0, 1}, {1, 0}, {0, 1}, {1, 0}}
	bias := [4][2]float32{{0, 0}, {0, 1}, {1, 0}, {0, 0}}

	// 根据线性插值计算等值线穿过边的位置
	for i := 0; i < 4; i++ {
		a, b := v[i], v[(i+1)%4]
		if !((threshold >= a && threshold >= b) || (threshold < a && threshold < b)) {
			rate := (threshold - a) / (b - a)
			if i > 1 {
				rate = 1 - rate
			}
			// fmt.Println(lat_init+lat_diff*(rate*rate_cal[i][0]+bias[i][0]), lon_init+lon_diff*(rate*rate_cal[i][1]+bias[i][1]))
			(*edges)[idxes[i][0]][idxes[i][1]].Location.Lat = lat_init + lat_diff*(rate*rate_cal[i][0]+bias[i][0])
			(*edges)[idxes[i][0]][idxes[i][1]].Location.Lon = lon_init + lon_diff*(rate*rate_cal[i][1]+bias[i][1])

		}
	}

	appendNeighbour := func(a int, b int) {
		// fmt.Println(a, b)
		// fmt.Println(idxes[a], idxes[b])
		*(*edges)[idxes[a][0]][idxes[a][1]].Neighbours = append(*(*edges)[idxes[a][0]][idxes[a][1]].Neighbours, *NewIdx(idxes[b]))
	}

	// 根据16种情况来添加相邻的边
	switch total {
	case 1:
		for i := range judge {
			if judge[i] {
				j := i - 1
				if j < 0 {
					j += 4
				}
				appendNeighbour(i, j)
				appendNeighbour(j, i)
				break
			}
		}
	case 3:
		for i := range judge {
			if !judge[i] {
				j := i - 1
				if j < 0 {
					j += 4
				}
				appendNeighbour(i, j)
				appendNeighbour(j, i)
				break
			}
		}
	case 2:
		if (judge[0] && judge[3]) || (judge[1] && judge[2]) {
			appendNeighbour(0, 2)
			appendNeighbour(2, 0)
		}
		if (judge[0] && judge[1]) || (judge[2] && judge[3]) {
			appendNeighbour(1, 3)
			appendNeighbour(3, 1)
		}
		if judge[0] && judge[2] {
			appendNeighbour(0, 1)
			appendNeighbour(1, 0)
			appendNeighbour(2, 3)
			appendNeighbour(3, 2)
		}
		if judge[1] && judge[3] {
			appendNeighbour(0, 3)
			appendNeighbour(3, 0)
			appendNeighbour(1, 2)
			appendNeighbour(2, 1)
		}
	}
}

func Search(i int, j int, edges *[][]EdgeInfo, partResult *[]LatLon) {
	// 类似dfs的方法搜索相邻边组成的等值线
	// 由于每条边最多两个邻居,向两个方向搜索直到形成环路或者仅有一个邻边

	idxes := (*edges)[i][j].Neighbours
	var idx_i, idx_j, prev_i, prev_j int
	var edgeInfo *EdgeInfo
	// fmt.Println(len(*idxes), *idxes)
	if len(*idxes) == 2 {
		idx_i, idx_j = i, j
		prev_i, prev_j = (*idxes)[0].Idx_i, (*idxes)[0].Idx_j
	begin_for_1:
		for {
			edgeInfo = &(*edges)[idx_i][idx_j]
			if edgeInfo.Passed {
				if len(*partResult) > 1 && idx_i == i && idx_j == j {
					*partResult = append(*partResult, *edgeInfo.Location)
				}
				break begin_for_1
			}
			edgeInfo.Passed = true
			*partResult = append(*partResult, *edgeInfo.Location)
			if len(*edgeInfo.Neighbours) < 2 && (idx_i != i || idx_j != j) {
				break begin_for_1
			}
			// fmt.Println(len(*edgeInfo.Neighbours))
			// fmt.Println(*edgeInfo.Neighbours)
			if prev_i == (*edgeInfo.Neighbours)[0].Idx_i && prev_j == (*edgeInfo.Neighbours)[0].Idx_j {
				prev_i, prev_j = idx_i, idx_j
				idx_i, idx_j = (*edgeInfo.Neighbours)[1].Idx_i, (*edgeInfo.Neighbours)[1].Idx_j
			} else {
				prev_i, prev_j = idx_i, idx_j
				idx_i, idx_j = (*edgeInfo.Neighbours)[0].Idx_i, (*edgeInfo.Neighbours)[0].Idx_j
			}
		}
		if len(*partResult) <= 1 || (*partResult)[0].Lat != (*partResult)[len(*partResult)-1].Lat || (*partResult)[0].Lon != (*partResult)[len(*partResult)-1].Lon {
			idx_i, idx_j = (*idxes)[0].Idx_i, (*idxes)[0].Idx_j
			prev_i, prev_j = i, j
		begin_for_2:
			for {
				edgeInfo = &(*edges)[idx_i][idx_j]
				if edgeInfo.Passed {
					if len(*partResult) > 1 && idx_i == i && idx_j == j {
						temp := []LatLon{*edgeInfo.Location}
						*partResult = append(temp, *partResult...)
					}
					break begin_for_2
				}
				edgeInfo.Passed = true
				temp := []LatLon{*edgeInfo.Location}
				*partResult = append(temp, *partResult...)
				if len(*edgeInfo.Neighbours) < 2 && (idx_i != i || idx_j != j) {
					break begin_for_2
				}
				if prev_i == (*edgeInfo.Neighbours)[0].Idx_i && prev_j == (*edgeInfo.Neighbours)[0].Idx_j {
					prev_i, prev_j = idx_i, idx_j
					idx_i, idx_j = (*edgeInfo.Neighbours)[1].Idx_i, (*edgeInfo.Neighbours)[1].Idx_j
				} else {
					prev_i, prev_j = idx_i, idx_j
					idx_i, idx_j = (*edgeInfo.Neighbours)[0].Idx_i, (*edgeInfo.Neighbours)[0].Idx_j
				}
			}
		}
	} else {
		idx_i, idx_j = i, j
		prev_i, prev_j = -1, -1
	begin_for_3:
		for {
			edgeInfo = &(*edges)[idx_i][idx_j]
			if edgeInfo.Passed {
				break begin_for_3
			}
			edgeInfo.Passed = true
			// fmt.Println(*partResult, edgeInfo.Location)
			*partResult = append(*partResult, *edgeInfo.Location)
			if len(*edgeInfo.Neighbours) < 2 && (idx_i != i || idx_j != j) {
				break begin_for_3
			}
			// fmt.Println(*edgeInfo.Neighbours)
			// fmt.Println(len(*edgeInfo.Neighbours))
			if prev_i == (*edgeInfo.Neighbours)[0].Idx_i && prev_j == (*edgeInfo.Neighbours)[0].Idx_j {
				prev_i, prev_j = idx_i, idx_j
				idx_i, idx_j = (*edgeInfo.Neighbours)[1].Idx_i, (*edgeInfo.Neighbours)[1].Idx_j
			} else {
				prev_i, prev_j = idx_i, idx_j
				idx_i, idx_j = (*edgeInfo.Neighbours)[0].Idx_i, (*edgeInfo.Neighbours)[0].Idx_j
			}
		}
	}
}

func IsolineWithThreshold(lats *[]float32, lons *[]float32, vals *[][]float32, threshold float32) (*[][]EdgeInfo, *[][]LatLon) {
	// 根据指定的阈值计算等值线

	n := len(*lats)
	m := len(*lons)

	edges := make([][]EdgeInfo, 2*n-1)
	for i := range edges {
		edges[i] = make([]EdgeInfo, m)
		for j := range edges[i] {
			neighbours := make([]Idx, 0, 2)
			location := LatLon{}
			edges[i][j] = EdgeInfo{Neighbours: &neighbours, Location: &location}
		}
	}

	results := make([][]LatLon, 0)
	for i := 0; i < n-1; i++ {
		for j := 0; j < m-1; j++ {
			BuildShape(i, j, threshold, &edges, lats, lons, vals)
		}
	}

	// for i := range edges {
	// 	for j := range (edges)[i] {
	// 		edges[i][j].print()
	// 	}
	// }

	for i := 0; i < 2*n-1; i++ {
		for j := 0; j < m; j++ {
			// fmt.Printf("%d,%d\n", i, j)
			if len(*edges[i][j].Neighbours) > 0 && !edges[i][j].Passed {
				partResult := make([]LatLon, 0, 10)
				Search(i, j, &edges, &partResult)
				results = append(results, partResult)
			}
		}
	}

	return &edges, &results
}
