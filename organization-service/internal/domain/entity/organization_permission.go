package entity

type OrganizationPermission struct {
	Id             uint32
	Name           string
	Key            string
	Description    string
	OrganizationId uint32
}
