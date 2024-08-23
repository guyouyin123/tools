package genModel

import (
	"testing"
)

func TestInitDb(t *testing.T) {
	//初始化
	//if err := gorm2.InitDb(); err != nil {
	//	panic(err)
	//}
	//
	//ctx := context.Background()
	//if err := InitDb(); err != nil {
	//	panic(err)
	//}
	//db := NewCtx(ctx)
	//query := dal.Use(db)
	//nameListQuery := query.NameList
	//
	////普通查询
	//infos, err := nameListQuery.Where(nameListQuery.NameListID.Eq(123)).Find()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(infos)
	//
	////事务
	//err = query.Transaction(func(tx *dal.Query) error {
	//	nameQ := tx.NameList
	//	info, _ := nameQ.Where(nameQ.NameListID.Eq(123)).Find()
	//	fmt.Println(info)
	//	return nil
	//})
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//type User struct {
	//}
	//
	////原生sql
	//sql := "select xx,oo from name_list"
	//userList := make([]*User, 0)
	//db.Exec(sql, &userList)

	//where 拼接
}
