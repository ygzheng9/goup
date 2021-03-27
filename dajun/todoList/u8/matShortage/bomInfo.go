package matShortage

import ()

// 加载 单层BOM 数据，根据母件，展开一层（穿透虚拟，停止在：外购，委外，领用层级）

// BOM 一层结构
type BOM struct {
	ChildInv   string  `db:"childInv" json:"childInv"`
	ChildName  string  `db:"childName" json:"childName"`
	MatType    int     `db:"matType" json:"matType"`
	ParentInv  string  `db:"parentInv" json:"parentInv"`
	ParentName string  `db:"parentName" json:"parentName"`
	BaseQty    float64 `db:"baseQty" json:"baseQty"`
}

// loadBOM 读取单层 BOM 数据
func loadBOM() ([]BOM, error) {
	sqlCmd := `
			select isnull(m1.cInvCode,'') childInv, isnull(m1.cInvName,'') childName,
			case when ((m1.bProxyForeign = 1) or (m1.bPurchase = 1) or (o.WIPType = 3 )) then 1 else 0 end 'matType',
			isnull(m2.cInvCode, '') parentInv, isnull(m2.cInvName, '') parentName, isnull(a.BaseQtyN, 0) baseQty
		from bom_opcomponent a
		inner join bom_opcomponentopt o on o.OptionsId = a.OpComponentId
		inner join bom_parent p on p.BomId = a.BomId
		inner join  bom_bom c on c.BomId = a.BomId and c.Status = 3
		inner join bas_part a1 on a1.PartId = a.ComponentId
		inner join Inventory m1 on m1.cInvCode = a1.InvCode
		inner join bas_part p1 on p1.PartId = p.ParentId
		inner join Inventory m2 on m2.cInvCode = p1.InvCode;
	`
	// fmt.Println(sqlCmd)

	items := []BOM{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// replacePart 使用替代料替换
func replacePart(inter InterchangeMatSlice, boms []BOM) []BOM {
	for idx, i := range boms {
		for _, a := range inter {
			// 子件在替代料清单中
			if i.ChildInv == a.SrcMat {
				// 替换子件的料号
				boms[idx].ChildInv = a.DestMat
				break
			}
		}
	}

	return boms
}

// 根据母件号，取得下一级子件
func getNextLevel(parentInv string, boms []BOM) []BOM {
	// 默认返回值
	result := []BOM{}
	for _, v := range boms {
		if v.ParentInv == parentInv {
			result = append(result, v)
		}
	}

	return result
}

// 根据母件，找到所有子件，穿透虚拟，停止在：自制领用，委外，外购
func findAllSubs(parentInv string, boms []BOM) []OneLevel {
	// 默认返回值
	results := []OneLevel{}

	// 母件压栈
	stack := NewStack()
	a := OneLevel{
		InvCode: parentInv,
		BaseQty: 1,
	}
	stack.push(a)

	// 设置最大单层子件数
	maxCount := 1000
	idx := 0

	for true {
		idx = idx + 1

		// stack 为空，表示已经处理完毕
		if (idx >= maxCount) || stack.isEmpty() {
			break
		}

		// 取出第一个
		current := stack.pop()
		nextLevels := getNextLevel(current.InvCode, boms)

		for _, s := range nextLevels {
			t := OneLevel{
				InvCode: s.ChildInv,
				BaseQty: current.BaseQty * s.BaseQty,
			}

			// 1 代表：采购或委外，自制领用，不需要再展开
			// 其余：还需要展开下一级
			if s.MatType == 1 {
				results = append(results, t)
			} else {
				stack.push(t)
			}
		}
	}

	if len(results) == 0 {
		// 如果没有下级，那么把自身加入到结果中
		a := OneLevel{
			InvCode: parentInv,
			BaseQty: 1,
		}
		results = append(results, a)
	}
	return results
}
