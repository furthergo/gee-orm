package clause

import "strings"

type Clause struct {
	sql map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause)Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
	}
	if c.sqlVars == nil {
		c.sqlVars = make(map[Type][]interface{})
	}
	q, vs := generators[name](vars...)
	c.sql[name] = q
	c.sqlVars[name] = vs
}

func (c *Clause)Build(names ...Type) (string, []interface{})  {
	qs := make([]string, 0)
	vars := make([]interface{}, 0)

	for _, n := range names {
		if q, ok := c.sql[n]; ok {
			vs := c.sqlVars[n]
			qs = append(qs, q)
			vars = append(vars, vs...)
		}
	}
	return strings.Join(qs, " "), vars
}