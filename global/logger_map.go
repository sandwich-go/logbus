package global

import (
	"sync"
)

// FIXME 用什么优化掉sync map
var LoggerMap sync.Map //used in basics and stdl
