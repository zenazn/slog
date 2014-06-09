package slog

import (
	"sort"
	"sync"
)

type logger struct {
	lock          sync.RWMutex
	parent        *logger
	context       map[string]interface{}
	rules         map[string]Level
	targets       map[Level]func(map[string]interface{})
	defaultTarget func(map[string]interface{})

	lcache *levelCache
	tcache *targetCache
}

// Get the cache corresponding to this logger. The logic is a bit subtle, but
// roughly speaking, your cache is valid if:
// - You are a root node (we define a modification to have occurred precisely
//   when the atomic set finishes, so whatever cache you happen to atomic get is
//   by definition valid)
// - Your cache is identical to a valid cache from your parent
// - Your cache's parent is identical to a valid cache from your parent
// In all other cases, you need to grab a lock and generate yourself a new
// cache.
func (l *logger) getLCache() *levelCache {
	if l.parent == nil {
		return l.atomicGetLCache()
	}

	pcache := l.parent.getLCache()
	cache := l.atomicGetLCache()

	if cache == pcache || cache.parent == pcache {
		return cache
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	if cache != l.atomicGetLCache() {
		// This means we raced with someone else on cache invalidation,
		// and is fine so long as the cache is invalidated infrequently
		// compared to how fast we can chew through our stack.
		return l.getLCache()
	}

	return l.genLCache(pcache)
}

// The caller must hold l's mutex.
func (l *logger) genLCache(pcache *levelCache) *levelCache {
	if len(l.context) == 0 {
		l.atomicSetLCache(pcache)
		return pcache
	}
	lc := &levelCache{
		iCache: make(map[uintptr]Level),
		parent: pcache,
		logger: l,
		rules:  l.generateRules(pcache),
	}
	l.atomicSetLCache(lc)
	return lc
}

// This is exactly the same mechanism as getLCache but for a different struct.
func (l *logger) getTCache() *targetCache {
	if l.parent == nil {
		return l.atomicGetTCache()
	}

	pcache := l.parent.getTCache()
	cache := l.atomicGetTCache()

	if cache == pcache || cache.parent == pcache {
		return cache
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	if cache != l.atomicGetTCache() {
		return l.getTCache()
	}

	return l.genTCache(pcache)
}

// The caller must hold l's mutex.
func (l *logger) genTCache(pcache *targetCache) *targetCache {
	if len(l.targets) == 0 && l.defaultTarget == nil {
		l.atomicSetTCache(pcache)
		return pcache
	}

	targets := make(map[Level]func(map[string]interface{}))
	if pcache != nil {
		for k, v := range pcache.targets {
			targets[k] = v
		}
	}
	for k, v := range l.targets {
		targets[k] = v
	}

	var defaultTarget func(map[string]interface{})
	if pcache != nil {
		defaultTarget = pcache.defaultTarget
	}
	if defaultTarget == nil {
		defaultTarget = l.defaultTarget
	}

	tc := &targetCache{
		targets:       targets,
		defaultTarget: defaultTarget,
		parent:        pcache,
	}
	l.atomicSetTCache(tc)
	return tc
}

func (l *logger) generateRules(pcache *levelCache) rules {
	rulesSet := make(map[string]Level, len(l.rules))
	if pcache != nil {
		for _, rule := range pcache.rules {
			rulesSet[rule.selector] = rule.level
		}
	}
	for k, v := range l.rules {
		rulesSet[k] = v
	}

	rules := make(rules, 0, len(rulesSet))
	for k, v := range rulesSet {
		rules = append(rules, rule{k, v})
	}
	sort.Sort(rules)

	return rules
}

func (l *logger) log(level Level, lines ...map[string]interface{}) bool {
	cache := l.getLCache()
	if ok := cache.shouldLogAt(level); !ok {
		return false
	}

	for _, line := range lines {
		m := make(map[string]interface{}, len(line)+len(l.context)+1)
		m["$level"] = level
		for k, v := range l.context {
			m[k] = v
		}
		for k, v := range line {
			m[k] = v
		}
		l.getTCache().dispatch(level, m)
	}

	return true
}

func (l *logger) Debug(lines ...map[string]interface{}) bool {
	return l.log(LDebug, lines...)
}
func (l *logger) Log(lines ...map[string]interface{}) bool {
	return l.log(LInfo, lines...)
}
func (l *logger) Warn(lines ...map[string]interface{}) bool {
	return l.log(LWarn, lines...)
}
func (l *logger) Error(lines ...map[string]interface{}) bool {
	return l.log(LError, lines...)
}

func (l *logger) Bind(context map[string]interface{}) Logger {
	l.lock.RLock()
	defer l.lock.RUnlock()
	ctx := make(map[string]interface{}, len(context)+len(l.context))
	for k, v := range l.context {
		ctx[k] = v
	}
	for k, v := range context {
		ctx[k] = v
	}
	return &logger{
		parent:  l,
		context: ctx,
		lcache:  l.atomicGetLCache(),
		tcache:  l.atomicGetTCache(),
	}
}

func (l *logger) SetLevel(selector string, level Level) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.rules == nil {
		l.rules = map[string]Level{selector: level}
	} else {
		l.rules[selector] = level
	}

	var pcache *levelCache
	if l.parent != nil {
		pcache = l.parent.getLCache()
	}
	l.genLCache(pcache)
}

func (l *logger) LogTo(ch chan<- string, levels ...Level) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.setTarget(func(line map[string]interface{}) {
		ch <- Format(line)
	}, levels)
}

func (l *logger) SlogTo(ch chan<- map[string]interface{}, levels ...Level) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.setTarget(func(line map[string]interface{}) {
		ch <- line
	}, levels)
}

func (l *logger) setTarget(fn func(map[string]interface{}), levels []Level) {
	if len(levels) == 0 {
		l.defaultTarget = fn
	}
	for _, level := range levels {
		l.targets[level] = fn
	}

	var pcache *targetCache
	if l.parent != nil {
		pcache = l.parent.getTCache()
	}
	l.genTCache(pcache)
}
