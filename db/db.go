package db

type DBOperation interface {
	Init()
	Select()
	Insert()
	Update()
	Delete()
}
