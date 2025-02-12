package qgb

import "bytes"

func DoNothing(fields ...string) *onConflict {
	return &onConflict{
		action: "DO NOTHING",
		fields: fields,
	}
}

func DoUpdate(set string, fields ...string) *onConflict {
	return &onConflict{
		action: "DO UPDATE ",
		set:    set,
		fields: fields,
	}
}

type onConflict struct {
	action string
	set    string
	fields []string
}

func (c *onConflict) build() string {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	buf.WriteString(" ON CONFLICT ")

	if len(c.fields) > 0 {
		buf.WriteString("(")

		for i := range c.fields {
			buf.WriteString(c.fields[i])

			if i < len(c.fields)-1 {
				buf.WriteString(", ")
			}
		}

		buf.WriteString(") ")
	}

	buf.WriteString(c.action)

	if c.set != "" {
		buf.WriteString(c.set)
	}

	return buf.String()
}
