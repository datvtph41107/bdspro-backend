package _dto

type Pagable struct {
	Page int    `form:"page"`
	Size int    `form:"size"`
	Sort string `form:"sort"`
}

func (p *Pagable) GetPage() int {
	if p.Page <= 0 {
		return 0
	}
	return p.Page
}

func (p *Pagable) GetSize() int {
	if p.Size <= 0 {
		return 20
	}
	return p.Size
}

func (p *Pagable) GetOffset() int {
	page := p.GetPage()
	size := p.GetSize()

	return page * size
}

func (p *Pagable) GetLimit() int {
	if p.Size <= 0 {
		return 20
	}
	if p.Size >= 100 {
		return 100
	}
	return p.Size
}
