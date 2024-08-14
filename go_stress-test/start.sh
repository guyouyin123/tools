#!/bin/bash
./go-stress-testing-mac -c {并发数} -n {每个并发执行请求的次数} -u http://127.0.0.1:8090/api/outer/login/user-id/get \
-H 'App-Id: bossup' \
-H 'Timestamp: 1718788806' \
-H 'Nonce: db5c88f2-098c-46f5-c941-bf16e4702f0c' \
-H 'Sign: for_woda_test' \
-H 'Sign-Ver: 1.0.0' \
-H 'Content-Type: application/json' \
-data '{"mobile": "15085141369"}'



go-stress-testing-mac -c 1 -n 10 -u http://127.0.0.1:8090/ping \
-H 'App-Id: bossup' \
-H 'Timestamp: 1718788806' \
-H 'Nonce: db5c88f2-098c-46f5-c941-bf16e4702f0c' \
-H 'Sign: for_woda_test' \
-H 'Sign-Ver: 1.0.0' \
-H 'Content-Type: application/json' \
-data '{"mobile": "15085141369"}'

