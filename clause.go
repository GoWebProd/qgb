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

type clause struct {
	op string

	field         string
	placeholder   string
	value         any
	needFieldLink bool

	sub []*clause
}

func (c *clause) toSQL(counter *counter) (string, []placeholderValue, error) {
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
	case "isnull":
		return c.field + " IS NULL", nil, nil
	case "notnull":
		return c.field + " IS NOT NULL", nil, nil
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

func (c *clause) getPlaceholder(counter *counter) string {
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

func (c *clause) valueMap() []placeholderValue {
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

func (c *clause) mergeSubs(counter *counter) ([]string, []placeholderValue, error) {
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

func (c *clause) buildAnd(counter *counter) (string, []placeholderValue, error) {
	clauses, args, err := c.mergeSubs(counter)
	if err != nil {
		return "", nil, err
	}

	return "(" + strings.Join(clauses, ") AND (") + ")", args, nil
}

func (c *clause) buildOr(counter *counter) (string, []placeholderValue, error) {
	clauses, args, err := c.mergeSubs(counter)
	if err != nil {
		return "", nil, err
	}

	return "(" + strings.Join(clauses, ") OR (") + ")", args, nil
}

func (c *clause) buildNot(counter *counter) (string, []placeholderValue, error) {
	if len(c.sub) != 1 || c.sub[0] == nil {
		return "", nil, errors.New("not clause must have one sub clause")
	}

	sql, args, err := c.sub[0].toSQL(counter)
	if err != nil {
		return "", nil, err
	}

	return "NOT (" + sql + ")", args, nil
}

func clauseInit(op string, field string, value any) *clause {
	return clauseInitWithSub(op, field, value, nil)
}

func clauseInitWithSub(op string, field string, value any, sub []*clause) *clause {
	var pholder string

	p, ok := value.(placeholder)
	if ok {
		value = nil
		pholder = p.name
	}

	return &clause{
		op:          op,
		field:       field,
		placeholder: pholder,
		value:       value,
		sub:         sub,
	}
}

func RAW(c string) *clause {
	return clauseInit("raw", c, nil)
}

func EQ(field string) *clause {
	return clauseInit("eq", field, nil)
}

func EQv(field string, value any) *clause {
	return clauseInit("eq", field, value)
}

func NEQ(field string) *clause {
	return clauseInit("neq", field, nil)
}

func NEQv(field string, value any) *clause {
	return clauseInit("neq", field, value)
}

func LT(field string) *clause {
	return clauseInit("lt", field, nil)
}

func LTv(field string, value any) *clause {
	return clauseInit("lt", field, value)
}

func GT(field string) *clause {
	return clauseInit("gt", field, nil)
}

func GTv(field string, value any) *clause {
	return clauseInit("gt", field, value)
}

func GTE(field string) *clause {
	return clauseInit("gte", field, nil)
}

func GTEv(field string, value any) *clause {
	return clauseInit("gte", field, value)
}

func LTE(field string) *clause {
	return clauseInit("lte", field, nil)
}

func LTEv(field string, value any) *clause {
	return clauseInit("lte", field, value)
}

func IN(field string, value any) *clause {
	return clauseInit("in", field, value)
}

func ISNULL(field string) *clause {
	return clauseInit("isnull", field, nil)
}

func NOTNULL(field string) *clause {
	return clauseInit("notnull", field, nil)
}

func AND(clauses ...*clause) *clause {
	return clauseInitWithSub("and", "", nil, clauses)
}

func OR(clauses ...*clause) *clause {
	return clauseInitWithSub("or", "", nil, clauses)
}

func NOT(c *clause) *clause {
	return clauseInitWithSub("not", "", nil, []*clause{c})
}
