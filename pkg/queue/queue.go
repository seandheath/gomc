package queue

type Queue struct {
	queue []func()
}

func NewQueue() *Queue {
	q := &Queue{}
	return q
}

func (q *Queue) Prepend(f func()) {
	q.queue = append([]func(){f}, q.queue...)
}
func (q *Queue) Append(f func()) {
	q.queue = append(q.queue, f)
}
func (q *Queue) Do() {
	for _, f := range q.queue {
		f()
	}
	q.Clear()
}

func (q *Queue) Clear() {
	q.queue = []func(){}
}
