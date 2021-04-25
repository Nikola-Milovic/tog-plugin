package engine

// SpatialNode keeps over all count of ids in sets for quick lookup
// Its a main building piece of hasher, as we do not expect big amounts of entities in a single node
// it does not use maps to store ids witch is not that elegant but fatser
type SpatialNode struct {
	Count int
	Sets  []Set
}

// Set is an id set that also has a group important for node
type Set struct {
	PlayerTag int
	IDs       []int
}

// Insert ...
func (n *SpatialNode) Insert(id int, tag int) {
	n.Count++
	for i := range n.Sets {
		s := &n.Sets[i]
		if s.PlayerTag == tag {
			s.IDs = append(s.IDs, id)
			return
		}
	}

	l := len(n.Sets)
	if cap(n.Sets) != l {
		n.Sets = n.Sets[:l+1]
		s := &n.Sets[l]
		s.PlayerTag = tag
		s.IDs[0] = id
		return
	}
	n.Sets = append(n.Sets, Set{tag, []int{id}})

}

// Remove panics if id does not exist within the node, you always have to make sure
// you are removing correctly as leaving dead ids in a hasher is leaking of memory
//
// method panics if object you tried to remove is not present to remove
func (n *SpatialNode) Remove(id int, group int) bool {
	n.Count--
	ll := len(n.Sets)
	var nil int // because this is a template
	for i := range n.Sets {
		s := &n.Sets[i]
		if s.PlayerTag == group {
			l := len(s.IDs)

			if l == 1 {
				if s.IDs[0] != id {
					return false
				}
				s.IDs[0] = nil
				n.Sets[i], n.Sets[ll-1] = n.Sets[ll-1], *s
				n.Sets = n.Sets[:ll-1]
				return true
			}

			for j := 0; j < l; j++ {
				if id == s.IDs[j] {
					l--
					s.IDs[j] = s.IDs[l]
					s.IDs[l] = nil
					s.IDs = s.IDs[:l]
					return true
				}
			}
		}
	}
	return false
}

// Collect retrieve ids from a node to coll, if include is true only ids of given group
// will get collected, otherwise ewerithing but specified group is returned
func (n *SpatialNode) Collect(group int, include bool, coll []int) []int {
	if include {
		for _, s := range n.Sets {
			if s.PlayerTag == group {
				coll = append(coll, s.IDs...)
				return coll
			}
		}
	} else {
		for _, s := range n.Sets {
			if s.PlayerTag != group {
				coll = append(coll, s.IDs...)
			}
		}
	}
	return coll
}

// CollectAll colects all objects withoud differentiating a group
func (n *SpatialNode) CollectAll(coll *[]int) {
	for _, s := range n.Sets {
		*coll = append(*coll, s.IDs...)
	}
}
