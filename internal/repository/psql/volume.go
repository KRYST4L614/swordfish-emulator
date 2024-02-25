package psql

import provider "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"

type VolumeRepository struct {
	db *provider.DbProvider
}

func NewVolumeRepository(db *provider.DbProvider) VolumeRepository {
	return VolumeRepository{
		db: db,
	}
}

func (r VolumeRepository) Create() {
	panic("not implemented") // TODO: Implement
}
func (r VolumeRepository) Delete() {
	panic("not implemented") // TODO: Implement
}

func (r VolumeRepository) Modify() {
	panic("not implemented") // TODO: Implement
}
