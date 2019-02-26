package repository

type Repository interface {
	GetOne(int) (interface{}, error)
	GetAll() []interface{}
	Create(interface{}) (int, error)
	CheckIfExists(int) bool
}
