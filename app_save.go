package e365_gateway

//func (app App[T]) CreateNew() AppItemForSave[T] {
//	return AppItemForSave[T]{
//		Data:  new(T),
//		isNew: true,
//		gen:   &app,
//	}
//}
//
//func (ais AppItemForSave[T]) Save() (*T, error) {
//	if ais.isNew {
//		return ais.gen.Create(ais.Data)
//	}
//	return ais.gen.Update(ais.Data, ais.id)
//	//return nil, nil
//}
//
//type AppItemForSave[T interface{}] struct {
//	Data  *T
//	isNew bool
//	gen   *App[T]
//	id    string
//}
//
//func NewAppItem[T interface{}]() AppItemForSave[T] {
//	return AppItemForSave[T]{Data: nil}
//}
