package pagination

import (
	"strconv"
)

type Pagination struct {
	Params map[string]interface{}
}

func (p *Pagination) Size() uint {
	if value, exists := p.Params["size"]; exists {
		size := value.([]string)[0]
		i, _ := strconv.Atoi(size)
		return uint(i)
	}
	return 1
}

func (p *Pagination) Offset() uint {
	page := 1
	if value, exists := p.Params["page"]; exists {
		size := value.([]string)[0]
		page, _ = strconv.Atoi(size)
	}

	return (p.Size() * uint(page-1))
}
