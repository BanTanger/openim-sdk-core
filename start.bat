cd /d %~p0./testv3new
go test -v -o pressure_test.test -run TestPressureTester_PressureSendMsgs -args -m=1000 -s=5338610321