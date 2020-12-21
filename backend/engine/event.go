package engine

//Event represents a single event that occured, events can be either something that happens to entity or that entity does
//movementEvent, applyPoisonEvent etc...
type Event struct {
	ID       string
	Index    int
	Data     map[string]interface{}
	Priority int
}

//EventPriorityQueue is
type eventPriorityQueue []Event

func (e eventPriorityQueue) Len() int           { return len(e) }
func (e eventPriorityQueue) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e eventPriorityQueue) Less(i, j int) bool { return e[i].Priority > e[j].Priority }

func (epq *eventPriorityQueue) Push(x interface{}) {
	no := x.(Event)
	*epq = append(*epq, no)
}

func (epq *eventPriorityQueue) Pop() interface{} {
	old := *epq
	n := len(old)
	no := old[n-1]
	*epq = old[0 : n-1]
	return no
}
