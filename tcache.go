package slog

type targetCache struct {
	targets       map[Level]func(map[string]interface{})
	defaultTarget func(map[string]interface{})
	parent        *targetCache
}

func (tc targetCache) dispatch(level Level, line map[string]interface{}) {
	if t, ok := tc.targets[level]; ok {
		t(line)
	} else {
		tc.defaultTarget(line)
	}
}
