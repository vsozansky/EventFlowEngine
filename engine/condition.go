package engine

// ConditionRule — интерфейс условия выполнения события.
//
// В событийной онтологии каждое событие может иметь условие (Causals),
// определяющее, при каких обстоятельствах событие может быть создано.
// ConditionRule — это дерево, где листья — проверки на существование
// события или равенство значения, а узлы — логические И/ИЛИ.
//
// Если ConditionRule == nil, условие считается всегда истинным.
type ConditionRule interface {
	// IsMet проверяет, выполнено ли условие в данном контексте.
	IsMet(ctx *EventContext) bool
}

// EventConditionRule — условие существования события.
//
// Проверяет, что событие с указанным ID существует в графе.
// Это самый простой тип условия — аналог ссылки на обуславливающее
// событие (Causals) в спецификации EventFlow.
//
// Пример: атрибут "Offer" доступен только если существует событие "Subject"
type EventConditionRule struct {
	// EventID — ID события, существование которого проверяется.
	EventID int
}

// IsMet возвращает true, если событие с EventID существует в контексте.
// Контекст должен содержать информацию о существующих событиях.
func (r *EventConditionRule) IsMet(ctx *EventContext) bool {
	return false // TODO: будет реализовано через Storage.GetEvent()
}

// PropertyEqualityRule — условие равенства значения свойства.
//
// Проверяет, что значение указанного свойства в контексте строго равно
// заданному значению. Используется для ветвления бизнес-логики.
//
// Пример: "Solution == Accept" — решение менеджера должно быть "Accept"
type PropertyEqualityRule struct {
	// PropertyID — ID свойства, значение которого проверяется.
	PropertyID int

	// Value — ожидаемое значение.
	Value string
}

// IsMet возвращает true, если значение свойства равно ожидаемому.
func (r *PropertyEqualityRule) IsMet(ctx *EventContext) bool {
	if ctx == nil || ctx.Properties == nil {
		return false
	}
	val, ok := ctx.Properties[r.PropertyID]
	return ok && val == r.Value
}

// PropertyInequalityRule — условие неравенства значения свойства.
//
// Проверяет, что значение указанного свойства не равно заданному.
// Используется для запрета действий при определённых значениях.
//
// Пример: "Status != closed" — статус не должен быть "closed"
type PropertyInequalityRule struct {
	// PropertyID — ID свойства, значение которого проверяется.
	PropertyID int

	// Value — значение, с которым сравнивается (не должно совпадать).
	Value string
}

// IsMet возвращает true, если значение свойства не равно ожидаемому.
func (r *PropertyInequalityRule) IsMet(ctx *EventContext) bool {
	if ctx == nil || ctx.Properties == nil {
		return false
	}
	val, ok := ctx.Properties[r.PropertyID]
	return ok && val != r.Value
}

// ConjunctionRule — логическое И (AND).
//
// Возвращает true только если ВСЕ дочерние правила истинны.
// Соответствует логическому умножению в языках программирования (&&).
//
// Пример: (Subject != "" && Solution == undefined)
type ConjunctionRule struct {
	// Values — список дочерних правил, все должны быть истинны.
	Values []ConditionRule
}

// IsMet возвращает true, если все дочерние правила истинны.
// Пустой список считается истинным (тривиальная конъюнкция).
func (r *ConjunctionRule) IsMet(ctx *EventContext) bool {
	if len(r.Values) == 0 {
		return true
	}
	for _, v := range r.Values {
		if v == nil || !v.IsMet(ctx) {
			return false
		}
	}
	return true
}

// DisjunctionRule — логическое ИЛИ (OR).
//
// Возвращает true если ХОТЯ БЫ ОДНО дочернее правило истинно.
// Соответствует логическому сложению в языках программирования (||).
//
// Пример: (Solution == "Reject" || Confirmation == "No")
type DisjunctionRule struct {
	// Values — список дочерних правил, хотя бы одно должно быть истинным.
	Values []ConditionRule
}

// IsMet возвращает true, если хотя бы одно дочернее правило истинно.
// Пустой список считается ложным (нет правил — нет истины).
func (r *DisjunctionRule) IsMet(ctx *EventContext) bool {
	for _, v := range r.Values {
		if v != nil && v.IsMet(ctx) {
			return true
		}
	}
	return false
}

// compile-time проверка: все типы реализуют ConditionRule
var _ ConditionRule = (*EventConditionRule)(nil)
var _ ConditionRule = (*PropertyEqualityRule)(nil)
var _ ConditionRule = (*PropertyInequalityRule)(nil)
var _ ConditionRule = (*ConjunctionRule)(nil)
var _ ConditionRule = (*DisjunctionRule)(nil)

// NilCondition — тривиальная реализация ConditionRule.
// Всегда возвращает true. Используется когда условие не задано.
type NilCondition struct{}

func (n *NilCondition) IsMet(_ *EventContext) bool {
	return true
}
