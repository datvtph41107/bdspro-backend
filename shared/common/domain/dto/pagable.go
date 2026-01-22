package _dto

type IPagable interface {
	GetPage() uint32
	GetSize() uint32
	GetOffset() int
	GetLimit() int
}

type Pagable struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
	Sort string `form:"sort"`
}

func (p *Pagable) GetPage() uint32 {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *Pagable) GetSize() uint32 {
	if p.Size <= 0 {
		return 20
	}
	if p.Size > 100 {
		return 100
	}
	return p.Size
}

func (p *Pagable) GetOffset() int {
	page := p.GetPage()
	size := p.GetSize()

	return int((page - 1) * size)
}

func (p *Pagable) GetLimit() int {
	return int(p.GetSize())
}
