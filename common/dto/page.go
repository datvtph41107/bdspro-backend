package _dto

type Page struct {
	Page int    `form:"page"`
	Size int    `form:"size"`
	Sort string `form:"sort"`
}

func (p *Page) GetPage() int {
	if p.Page <= 0 {
		return 0
	}
	return p.Page
}

func (p *Page) GetSize() int {
	if p.Size <= 0 {
		return 20
	}
	return p.Size
}

func (p *Page) GetOffset() int {
	page := p.GetPage()
	size := p.GetSize()

	return page * size
}

func (p *Page) GetLimit() int {
	if p.Size <= 0 {
		return 20
	}
	return p.Size
}
