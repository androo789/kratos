package log

// FilterOption is filter option.
type FilterOption func(*Filter)

const fuzzyStr = "***"

// FilterLevel with filter level.
// 屏蔽一个level
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKey with filter key.
// 屏蔽指定key
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.key[v] = struct{}{}
		}
	}
}

// FilterValue with filter value.
// 屏蔽指定value
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.value[v] = struct{}{}
		}
	}
}

// FilterFunc with filter func.
func FilterFunc(f func(level Level, keyvals ...interface{}) bool) FilterOption {
	return func(o *Filter) {
		o.filter = f
	}
}

// Filter is a logger filter.
type Filter struct {
	logger Logger
	//屏蔽level
	level Level
	//屏蔽key
	key map[interface{}]struct{}
	//屏蔽value
	value map[interface{}]struct{}
	//自定义屏蔽规则
	filter func(level Level, keyvals ...interface{}) bool
}

// NewFilter new a logger filter.
func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := Filter{
		logger: logger,
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

// Log Print log by level and keyvals.
// 实现了log接口
func (f *Filter) Log(level Level, keyvals ...interface{}) error {
	if level < f.level {
		return nil
	}
	// prefixkv contains the slice of arguments defined as prefixes during the log initialization
	var prefixkv []interface{}
	l, ok := f.logger.(*logger)
	if ok && len(l.prefix) > 0 {
		prefixkv = make([]interface{}, 0, len(l.prefix))
		prefixkv = append(prefixkv, l.prefix...)
	}

	if f.filter != nil && (f.filter(level, prefixkv...) || f.filter(level, keyvals...)) {
		return nil
	}

	if len(f.key) > 0 || len(f.value) > 0 {
		for i := 0; i < len(keyvals); i += 2 {
			v := i + 1
			if v >= len(keyvals) {
				continue
			}
			if _, ok := f.key[keyvals[i]]; ok {
				keyvals[v] = fuzzyStr
			}
			if _, ok := f.value[keyvals[v]]; ok {
				keyvals[v] = fuzzyStr
			}
		}
	}
	return f.logger.Log(level, keyvals...)
}
