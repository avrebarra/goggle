package ctxboard

import "sync"

func countmap(m *sync.Map) (out int) {
	m.Range(func(_, _ any) bool { out++; return true })
	return
}
