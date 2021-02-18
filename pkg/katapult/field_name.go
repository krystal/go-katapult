package katapult

type FieldName int

const (
	AddressField FieldName = iota
	FQDNField
	IDField
	NameField
	ObjectIDField
	PermalinkField
	SubDomainField
)
