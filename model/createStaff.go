package model

type CreateStaffInput struct {
	Name               string `json:"name" validate:"required"`
	PhoneNumber        string `json:"phoneNumber" gorm:"uniqueIndex;not null" validate:"required"`
	Email              string `json:"email" gorm:"uniqueIndex;not null" `
	IdentificationCard string `gorm:"not null;uniqueIndex;require" validate:"required,min=12,max=12" json:"identificationCard"`
	Password           string `json:"password" validate:"required,min=6,max=50"`
	Position           string `json:"position"`
	StatusWorking      string `json:"statusWorking"`
	Gender             string `json:"gender"`
	SearchKeyStr       string `json:"searchKey"`
	Act                bool   `json:"active"`
}
type CreateStaffInputs []CreateStaffInput

func (s CreateStaffInput) SearchKey() string {
	return s.SearchKeyStr

}
func (s CreateStaffInput) Active() bool {
	return s.Act

}
