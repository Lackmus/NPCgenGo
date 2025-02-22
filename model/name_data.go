package model

// NameData : A name has a name, a list of forenames and a list of surnames.
type NameData struct {
	Name      string
	Forenames []string
	Surnames  []string
}

// Name : Return the name of the NameData
func (n NameData) GetName() string {
	return n.Name
}
