package testv3new

import (
	"flag"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/testv3new/testcore"
	"strings"
	"testing"
	"time"
)

type sliceValue []string

func (s *sliceValue) String() string {
	return ""
}

func (s *sliceValue) Set(val string) error {
	*s = sliceValue(strings.Split(val, ","))
	return nil
}

var (
	messageNumber int
	sendUserIDs   sliceValue
	recvUserIDs   sliceValue
	groupID       string
	timeInterval  int64
)

func init() {
	flag.IntVar(&messageNumber, "m", 0, "messageNumber for single sender")
	flag.Var(&sendUserIDs, "s", "sender id list")
	flag.Var(&sendUserIDs, "r", "recv id list")
	flag.StringVar(&groupID, "g", "", "groupID for testing")
	flag.Int64Var(&timeInterval, "t", 0, "timeInterVal during sending message")
	flag.Parse()
	if err := log.InitFromConfig("sdk.log", "sdk", 3,
		true, false, "", 2); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	fmt.Println(messageNumber, sendUserIDs, recvUserIDs, groupID, timeInterval)
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, recvUserIDs, messageNumber, time.Duration(timeInterval)*time.Millisecond)
		time.Sleep(time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendGroupMsgs2)(sendUserIDs, groupID, messageNumber, time.Duration(timeInterval)*time.Millisecond)
		time.Sleep(time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func Test_WithTimer(t *testing.T) {
	p := WithTimer(NewPressureTester)(testcore.APIADDR, testcore.WSADDR)
	fmt.Println("test test", p)
	tester := p[0].Interface().(*PressureTester)
	add := tester.add(1, 2)
	fmt.Println("test test", add)
}
