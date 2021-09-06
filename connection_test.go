package FunRTC

import (
	"fmt"
	"github.com/chet-hub/FunRTC/utils"
	"testing"
)

/**
Js method for test

function PostData(url, str) {
    return new Promise((ok, error) => {
        const xhr = new XMLHttpRequest();
        xhr.onload = () => {
            if (xhr.readyState === 4 && xhr.status === 200) {
                ok(xhr.responseText);
            } else {
                error(xhr.status + ":" + xhr.statusText)
            }
        }
        xhr.open('POST', url, true);
        xhr.setRequestHeader('Content-Type', 'text/plain');
        xhr.send(str);
    })
}

 */


func TestConnection_as_Client(t *testing.T) {
	var http = utils.HTTPServer()
	var connection, _ = new(nil, nil, "", "")

	msg := <-http
	<-msg.Request
	var applyString, err = connection.apply()
	if err != nil {
		fmt.Print("apply Error" + applyString + "\n")
		t.Fail()
	}
	msg.Response <- applyString

	msg = <-http
	approveString := <-msg.Request
	err = connection.connect(approveString)
	if err != nil {
		fmt.Print(" connect Error \n")
		t.Fail()
	}
	msg.Response <- "ok"

	/**
		start this test method and run the js coding in Chrome

	    const applyPaper = await PostData('//localhost:8080/msg', "getApplyPaper")
	    const approvePaper = await FunRTC.Accept("golang client",applyPaper)
	    const connect = await PostData('//localhost:8080/msg', approvePaper)
	*/

}

func TestConnection_as_Server(t *testing.T) {
	var http = utils.HTTPServer()
	var connection, _ = new(nil, nil, "", "")
	msg := <-http
	re := <-msg.Request
	var approveString, err = connection.approve(re)
	if err != nil {
		msg.Response <- ""
		t.Fail()
	}
	msg.Response <- approveString

	/**
	  start this test method and run the js coding in Chrome

	  let applyPaper = await FunRTC.ToConnect("Golang server")
	  const approvePaper = await PostData('//localhost:8080/msg', applyPaper)
	  await FunRTC.DoConnect("Golang server", approvePaper)
	*/
}


func TestConnection_as_ServerAndClient(t *testing.T) {

	var Local, _ = new(nil, nil, "", "")
	var Remote, _ = new(nil, nil, "", "")

	localApply,e := Local.apply()
	if e != nil {
		t.Fail()
	}
	remoteApprove,e := Remote.approve(localApply)
	if e != nil {
		t.Fail()
	}
	e = Local.connect(remoteApprove)
	if e != nil {
		t.Fail()
	}

}