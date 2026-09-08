package _dto

type IPagable interface {
	GetPage() uint32
	GetSize() uint32
	GetOffset() int
	GetLimit() int
}

type Pagable struct {
	Text string `form:"text"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
	Sort string `form:"sort"`
}

func (p *Pagable) GetPage() uint32 {
	if p.Page <= 0 {
		return 0
	}
	return p.Page
}

func (p *Pagable) GetSize() uint32 {
	if p.Size <= 0 {
		return 20
	}
	return p.Size
}

func (p *Pagable) GetOffset() int {
	page := p.GetPage()
	size := p.GetSize()

	if page < 1 {
		return 0
	}
	return int(page) * int(size)
}

func (p *Pagable) GetLimit() int {
	if p.Size <= 0 {
		return 20
	}
	if p.Size >= 100 {
		return 100
	}
	return int(p.Size)
}

func (p *Pagable) GetLimitAdmin() int {
	if p.Size <= 0 {
		return 20
	}
	if p.Size >= 999 {
		return 999
	}
	return int(p.Size)
}

func (p *Pagable) GetSort() string {
	if p.Sort == "" {
		return "created_at desc"
	}
	return p.Sort
}

func (p *Pagable) Normalize() {
	if p.Page <= 0 {
		p.Page = 0
	}
	if p.Size <= 0 {
		p.Size = 20
	}
	// if p.Size > 100 {
	// 	p.Size = 100
	// }
}

func NewPagableFromGrpc(page *uint32, size *uint32, sort *string) *Pagable {
	p := &Pagable{}

	if page != nil {
		p.Page = *page
	}
	if size != nil {
		p.Size = *size
	}
	if sort != nil {
		p.Sort = *sort
	}

	p.Normalize()

	return p
}
