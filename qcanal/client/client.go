package client

import (
	"fmt"
	"time"

	protocolEntry "github.com/withlin/canal-go/protocol/entry"

	"github.com/withlin/canal-go/client"

	"google.golang.org/protobuf/proto"
)

// canal-go client demo
var (
	Debug = true
)

func clientRun() {
	// 连接canal-server
	// 请修改为你的 canal-server 配置
	connector := client.NewSimpleCanalConnector(
		"127.0.0.1", 11111, "", "", "example", 60000, 60*60*1000)
	err := connector.Connect()
	if err != nil {
		panic(err)
	}

	// mysql 数据解析关注的表，Perl正则表达式.
	//err = connector.Subscribe(".*\\..*")
	//if err != nil {
	//	fmt.Printf("connector.Subscribe failed, err:%v\n", err)
	//	panic(err)
	//}

	// 消费消息
	for {
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			fmt.Printf("connector.Get failed, err:%v\n", err)
			continue
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(time.Second)
			continue
		}
		printEntry2(message.Entries)
	}
}

func printEntry2(entries []protocolEntry.Entry) {
	for _, entry := range entries {
		// 忽略事务开启和事务关闭类型
		if entry.GetEntryType() == protocolEntry.EntryType_TRANSACTIONBEGIN ||
			entry.GetEntryType() == protocolEntry.EntryType_TRANSACTIONEND {
			continue
		}
		// RowChange对象，包含了一行数据变化的所有特征
		rowChange := new(protocolEntry.RowChange)
		// protobuf解析
		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		if err != nil {
			fmt.Printf("proto.Unmarshal failed, err:%v\n", err)
		}
		if rowChange == nil {
			continue
		}

		// 获取并打印Header信息
		header := entry.GetHeader()

		//过滤库名和表名
		//if !FilterMap(header) {
		//	continue
		//}

		fmt.Printf("binlog[%s : %d], name[%s,%s], eventType: %s\n",
			header.GetLogfileName(),
			header.GetLogfileOffset(),
			header.GetSchemaName(),
			header.GetTableName(),
			header.GetEventType(),
		)
		//判断是否为DDL语句 --理论业务没有这样的语句
		if rowChange.GetIsDdl() {
			fmt.Printf("isDdl:true, sql:%v\n", rowChange.GetSql())
		}

		// 获取操作类型：insert/update/delete等
		eventType := rowChange.GetEventType()
		fmt.Println(fmt.Sprintf("==> binlog[%s:%d], name[%s.%s], eventType[%s]", header.GetLogfileName(),
			header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))

		switch eventType {
		case protocolEntry.EventType_CREATE, // create table
			protocolEntry.EventType_ALTER,    // alter table
			protocolEntry.EventType_ERASE,    // drop table
			protocolEntry.EventType_QUERY,    // create/drop database
			protocolEntry.EventType_TRUNCATE, // truncate table
			protocolEntry.EventType_RENAME,   // rename table
			protocolEntry.EventType_CINDEX,   // create index
			protocolEntry.EventType_DINDEX:   // drop index
			sqlStr := rowChange.GetSql()
			fmt.Println(sqlStr)
		case protocolEntry.EventType_INSERT:
			for _, rowData := range rowChange.GetRowDatas() {
				sql, err := insertSql(header, rowData)
				if err != nil {
					fmt.Println("insertSql error:", err)
					continue
				}
				fmt.Println("======INSERT======")
				fmt.Println(sql)
			}
		case protocolEntry.EventType_UPDATE:
			for _, rowData := range rowChange.GetRowDatas() {
				sql, err := updateSql(header, rowData)
				if err != nil {
					fmt.Println("insertSql error:", err)
					continue
				}
				fmt.Println("======UPDATE======")
				fmt.Println(sql)
			}
		case protocolEntry.EventType_DELETE:
			for _, rowData := range rowChange.GetRowDatas() {
				sql, err := deleteSql(header, rowData)
				if err != nil {
					fmt.Println("insertSql error:", err)
					continue
				}
				fmt.Println("======DELETE======")
				fmt.Println(sql)
			}
		default:
			logWhenOthers(entry)
		}
	}
}
