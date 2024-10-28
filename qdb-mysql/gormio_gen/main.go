package main

import (
	"context"
	"fmt"
	"github.com/guyouyin123/tools/qdb-mysql/gormio_gen/dal"
	"github.com/guyouyin123/tools/qdb-mysql/gormio_gen/gen"
	"github.com/guyouyin123/tools/qdb-mysql/gormio_gen/gormio"
)

func main() {
	basePath := "./qdb-mysql/gormio_gen/dal"
	gen.Run(basePath)

	//TestDB()
}

func TestDB() {
	gormio.InitDb("")

	ctx := context.Background()
	db := gormio.NewCtx(ctx)
	query := dal.Use(db)
	xq := query.User

	infos, err := xq.GetInfosMapByIDs([]int32{1, 2, 3, 4, 5, 6})
	if err != nil {
		panic(err)
	}
	fmt.Println(infos)
}
