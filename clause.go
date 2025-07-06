package qgb

import (
	"errors"
	"strings"
)

type placeholder struct {
	name string
}

func Placeholder(name string) any {
	return placeholder{name}
}

type Clause struct {
	op string

	field         string
	placeholder   string
	value         any
	needFieldLink bool

	sub []*Clause
}

func (c *Clause) toSQL(counter *counter) (string, []placeholderValue, error) {
	switch c.op {
	case "raw":
		return c.field, nil, nil
	case "eq":
		return c.field + " = " + c.getPlaceholder(counter), c.valueMap(), nil
	case "neq":
		return c.field + " <> " + c.getPlaceholder(counter), c.valueMap(), nil
	case "gt":
		return c.field + " > " + c.getPlaceholder(counter), c.valueMap(), nil
	case "gte":
		return c.field + " >= " + c.getPlaceholder(counter), c.valueMap(), nil
	case "lt":
		return c.field + " < " + c.getPlaceholder(counter), c.valueMap(), nil
	case "lte":
		return c.field + " <= " + c.getPlaceholder(counter), c.valueMap(), nil
	case "in":
		return c.field + " IN " + c.getPlaceholder(counter), c.valueMap(), nil
	case "any":
		return c.field + " = ANY(" + c.getPlaceholder(counter) + ")", c.valueMap(), nil
	case "isnull":
		return c.field + " IS NULL", nil, nil
	case "notnull":
		return c.field + " IS NOT NULL", nil, nil
	case "contains":
		return c.field + " @> " + c.getPlaceholder(counter), c.valueMap(), nil
	case "and":
		return c.buildAnd(counter)
	case "or":
		return c.buildOr(counter)
	case "not":
		return c.buildNot(counter)
	default:
		return "", nil, errors.New("unknown clause")
	}
}

func (c *Clause) getPlaceholder(counter *counter) string {
	if c.placeholder == "" {
		c.placeholder = c.field + counter.IncrementString()
		c.needFieldLink = true
	}

	return "@" + c.placeholder
}

type placeholderValue struct {
	field string
	value any
}

func (c *Clause) valueMap() []placeholderValue {
	if c.value != nil {
		return []placeholderValue{
			{field: c.placeholder, value: c.value},
		}
	}

	var name string
	if c.needFieldLink {
		name = c.field
	}

	return []placeholderValue{
		{field: c.placeholder, value: placeholder{name}},
	}
}

func (c *Clause) mergeSubs(counter *counter) ([]string, []placeholderValue, error) {
	clauses := make([]string, 0, len(c.sub))
	args := make([]placeholderValue, 0, len(c.sub))

	for _, sub := range c.sub {
		clause, arg, err := sub.toSQL(counter)
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)

		args = append(args, arg...)
	}

	return clauses, args, nil
}

func (c *Clause) buildAnd(counter *counter) (string, []placeholderValue, error) {
	clauses, args, err := c.mergeSubs(counter)
	if err != nil {
		return "", nil, err
	}

	return "(" + strings.Join(clauses, ") AND (") + ")", args, nil
}

func (c *Clause) buildOr(counter *counter) (string, []placeholderValue, error) {
	clauses, args, err := c.mergeSubs(counter)
	if err != nil {
		return "", nil, err
	}

	return "(" + strings.Join(clauses, ") OR (") + ")", args, nil
}

func (c *Clause) buildNot(counter *counter) (string, []placeholderValue, error) {
	if len(c.sub) != 1 || c.sub[0] == nil {
		return "", nil, errors.New("not clause must have one sub clause")
	}

	sql, args, err := c.sub[0].toSQL(counter)
	if err != nil {
		return "", nil, err
	}

	return "NOT (" + sql + ")", args, nil
}

func clauseInit(op string, field string, value any) *Clause {
	return clauseInitWithSub(op, field, value, nil)
}

func clauseInitWithSub(op string, field string, value any, sub []*Clause) *Clause {
	var pholder string

	p, ok := value.(placeholder)
	if ok {
		value = nil
		pholder = p.name
	}

	return &Clause{
		op:          op,
		field:       field,
		placeholder: pholder,
		value:       value,
		sub:         sub,
	}
}

func RAW(c string) *Clause {
	return clauseInit("raw", c, nil)
}

func EQ(field string) *Clause {
	return clauseInit("eq", field, nil)
}

func EQv(field string, value any) *Clause {
	return clauseInit("eq", field, value)
}

func NEQ(field string) *Clause {
	return clauseInit("neq", field, nil)
}

func NEQv(field string, value any) *Clause {
	return clauseInit("neq", field, value)
}

func LT(field string) *Clause {
	return clauseInit("lt", field, nil)
}

func LTv(field string, value any) *Clause {
	return clauseInit("lt", field, value)
}

func GT(field string) *Clause {
	return clauseInit("gt", field, nil)
}

func GTv(field string, value any) *Clause {
	return clauseInit("gt", field, value)
}

func GTE(field string) *Clause {
	return clauseInit("gte", field, nil)
}

func GTEv(field string, value any) *Clause {
	return clauseInit("gte", field, value)
}

func LTE(field string) *Clause {
	return clauseInit("lte", field, nil)
}

func LTEv(field string, value any) *Clause {
	return clauseInit("lte", field, value)
}

func IN(field string, value any) *Clause {
	return clauseInit("in", field, value)
}

func ANY(field string, value any) *Clause {
	return clauseInit("any", field, value)
}

func ISNULL(field string) *Clause {
	return clauseInit("isnull", field, nil)
}

func NOTNULL(field string) *Clause {
	return clauseInit("notnull", field, nil)
}

func CONTAINS(field string, value any) *Clause {
	return clauseInit("contains", field, value)
}

func AND(clauses ...*Clause) *Clause {
	return clauseInitWithSub("and", "", nil, clauses)
}

func OR(clauses ...*Clause) *Clause {
	return clauseInitWithSub("or", "", nil, clauses)
}

func NOT(c *Clause) *Clause {
	return clauseInitWithSub("not", "", nil, []*Clause{c})
}
