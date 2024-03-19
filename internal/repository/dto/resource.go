package dto

type ResourceDto struct {
	Id   string `db:"id"`
	Data []byte `db:"data"`
}
