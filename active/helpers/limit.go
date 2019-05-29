package helpers

type Limit struct {
	limits map[string][]string
}

func NewLimit() Limit {
	return Limit{limits: make(map[string][]string, 0)}
}

func NewLimitFromMap(entryMap map[string][]string) Limit {
	return Limit{limits: entryMap }
}

func (l *Limit) AddLimit(tableName string, colNames ...string) {
	l.limits[tableName] = colNames
}

func (l *Limit) GetLimits() map[string][]string {
	return l.limits
}

func (l *Limit) Valid() bool {
	if len(l.limits) > 0 {
		return true
	}

	return false
}