package lunyu

func Minus(old []string, curr []string) (removed []string) {
	removed = []string{}
	for _, v := range old {
		ok := true

		for _, k := range curr {
			if v == k {
				ok = false
				break
			}
		}

		if ok {
			removed = append(removed, v)
		}
	}

	return removed
}

func Intersect(old []string, curr []string) (added []string, removed []string) {
	added = []string{}
	for _, v := range curr {
		ok := true

		for _, k := range old {
			if v == k {
				ok = false
				break
			}
		}

		if ok {
			added = append(added, v)
		}
	}

	removed = []string{}
	for _, v := range old {
		ok := true

		for _, k := range curr {
			if v == k {
				ok = false
				break
			}
		}

		if ok {
			removed = append(removed, v)
		}
	}

	return added, removed
}
