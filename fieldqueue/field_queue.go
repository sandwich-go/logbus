package fieldqueue

import (
	"sync"

	"go.uber.org/zap"
)

type FieldQueue struct {
	buff []zap.Field
}

var queuePool = sync.Pool{
	New: func() interface{} {
		return new(FieldQueue)
	},
}

func NewQueue() *FieldQueue {
	q := queuePool.Get().(*FieldQueue)
	q.Reset()
	return q
}

func (fq *FieldQueue) Reset() {
	fq.buff = fq.buff[:0]
}

func (fq *FieldQueue) Push(field zap.Field) {
	fq.buff = append(fq.buff, field)
}

// 每个FieldQueue应该只被调用一次
func (fq *FieldQueue) Retrieve() (data []zap.Field) {
	data = fq.buff
	queuePool.Put(fq)
	return
}
