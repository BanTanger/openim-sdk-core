cd ./single_test
go build -o msg_delay.exe msg_delay_open_im.go
cd ../testv3new
go test -c pressure_test.go pressure_tester.go register_manager.go -o pressure_test.test