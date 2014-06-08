package slog

import (
	"runtime"
	"strings"
	"sync"
)

type levelCache struct {
	sync.RWMutex
	iCache map[uintptr]Level
	parent *levelCache
	logger *logger
	rules  rules
}

func (lc *levelCache) shouldLogAt(level Level) bool {
	pc, _, _, ok := runtime.Caller(3)

	// Unclear when this would happen, but let's fail open instead of
	// closed.
	if !ok {
		return true
	}

	lc.RLock()
	l, ok := lc.iCache[pc]
	lc.RUnlock()

	if ok {
		return l <= level
	}

	f := runtime.FuncForPC(pc)
	l = lc.levelForFunc(f.Name())

	lc.Lock()
	lc.iCache[pc] = l
	lc.Unlock()

	return l <= level
}

func (lc *levelCache) levelForFunc(fname string) Level {
	for _, rule := range lc.rules {
		if !strings.HasPrefix(fname, rule.selector) {
			continue
		}
		ftail := fname[len(rule.selector):]
		if strings.HasSuffix(rule.selector, ".") {
			// For pattern "foo.", don't match anything in a
			// hypothetical oddly-named package "foo.bar".
			if strings.ContainsRune(ftail, '/') {
				continue
			}
		} else if !strings.HasSuffix(rule.selector, "/") {
			// For pattern "foo", match "foo.bar" and "foo/bar.baz",
			// but not "foobar".
			if len(ftail) > 0 && ftail[0] != '/' && ftail[0] != '.' {
				continue
			}
		}

		return rule.level
	}

	return DefaultLevel
}
