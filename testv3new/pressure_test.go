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
	sendUserID    sliceValue
	recvUserID    sliceValue
	groupID       string
	timeInterval  int64
)

func init() {
	flag.IntVar(&messageNumber, "m", 0, "messageNumber for single sender")
	flag.Var(&sendUserID, "s", "sender id list")
	flag.Var(&recvUserID, "r", "recv id list")
	flag.StringVar(&groupID, "g", "", "groupID for testing")
	flag.Int64Var(&timeInterval, "t", 0, "timeInterVal during sending message")
	flag.Parse()
	if err := log.InitFromConfig("sdk.log", "sdk", 3,
		true, false, "", 2); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	fmt.Println(messageNumber, sendUserID, recvUserID, groupID, timeInterval)
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendMsgs2)(sendUserID, recvUserID, messageNumber, time.Duration(timeInterval)*time.Millisecond)
		time.Sleep(time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	sendUserID := "register_test_4334"
	groupID := "3411007805"

	pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	pressureSendGroupMsgsWithTime := pressureTester.WithTimer(pressureTester.PressureSendGroupMsgs)
	pressureSendGroupMsgsWithTime([]string{sendUserID}, groupID, 100, time.Duration(100))
}

func TestPressureTester_PressureSendGroupMsgs2(t *testing.T) {
	start := 850
	count := 900
	step := 10
	for j := start; j <= count; j += step {
		var sendUserIDs []string
		startTime := time.Now().UnixNano()
		for i := j; i < j+step; i++ {
			sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
		}
		// groupID := "3411007805"
		// groupID := "2347514573"
		groupID := "3167736657"
		pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
		pressureTester.PressureSendGroupMsgs(sendUserIDs, groupID, 1, 0)
		endTime := time.Now().UnixNano()
		nanoSeconds := float64(endTime - startTime)
		t.Log("", nanoSeconds)
		fmt.Println()
	}
}

func TestPressureTester_Conversation(t *testing.T) {
	sendUserID := "5338610321"
	var recvUserIDs []string
	for i := 1; i <= 1000; i++ {
		recvUserIDs = append(recvUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	p.WithTimer(p.PressureSendMsgs)(sendUserID, recvUserIDs, 1, 100*time.Millisecond)
}

func TestPressureTester_PressureSendMsgs2(t *testing.T) {
	recvUserID := "5338610321"
	var sendUserIDs []string
	for i := 1; i <= 1000; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, []string{recvUserID}, 1, 100*time.Millisecond)
}
