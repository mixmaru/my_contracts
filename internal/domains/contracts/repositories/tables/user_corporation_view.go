package tables

type UserCorporationView struct {
	UserRecord
	ContactPersonName string `db:"contact_person_name"`
	PresidentName     string `db:"president_name"`
}
