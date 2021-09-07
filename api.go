package FunRTC

import (
	"errors"
	"github.com/pion/webrtc/v3"
	"strings"
)

var connectMap = make(map[string]*connection)

func newConnection(name string, OnCreateDataChannelFunc func(*webrtc.DataChannel),
	OnConnectionStateChangeFunc func(webrtc.PeerConnectionState),
	webrtcConfiguration string,
	dataChannelOptions string) (*connection, error) {

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("please specify name for this connection")
	}
	var connection, e = new(OnCreateDataChannelFunc, OnConnectionStateChangeFunc, webrtcConfiguration, dataChannelOptions)
	if e != nil {
		return nil, e
	}

	return connection, nil
}

func ToConnect(name string, OnCreateDataChannelFunc func(*webrtc.DataChannel),
	OnConnectionStateChangeFunc func(webrtc.PeerConnectionState),
	webrtcConfiguration string,
	dataChannelOptions string) (string, error) {

	if strings.TrimSpace(name) == "" {
		return "", errors.New("please specify name that you want to connect")
	}
	if connectMap[name] != nil {
		return "", errors.New("please use a different name")
	}
	var connection, e = new(OnCreateDataChannelFunc, OnConnectionStateChangeFunc, webrtcConfiguration, dataChannelOptions)
	if e != nil {
		return "", e
	}
	connectMap[name] = connection
	return connection.apply()
}

func Accept(name string, remoteApplyString string, OnCreateDataChannelFunc func(*webrtc.DataChannel),
	OnConnectionStateChangeFunc func(webrtc.PeerConnectionState),
	webrtcConfiguration string,
	dataChannelOptions string) (string, error) {

	if strings.TrimSpace(name) == "" {
		return "", errors.New("please specify name who you want to Accept")
	}
	if connectMap[name] != nil {
		return "", errors.New("please use a different name")
	}
	var connection, e = new(OnCreateDataChannelFunc, OnConnectionStateChangeFunc, webrtcConfiguration, dataChannelOptions)
	if e != nil {
		return "", e
	}
	connectMap[name] = connection
	return connection.approve(remoteApplyString)
}

func DoConnect(name string, remoteApprove string) error {
	if connectMap[name] == nil {
		return errors.New("please use a the name that used in the apply")
	}
	return connectMap[name].connect(remoteApprove)
}

func SendText(name string, text string) error {
	if connectMap[name] != nil {
		return errors.New("please use a different name")
	}
	return connectMap[name].DataChannel.SendText(text)
}

func Send(name string, data []byte) error {
	if connectMap[name] != nil {
		return errors.New("please use a different name")
	}
	return connectMap[name].DataChannel.Send(data)
}

func Close(name string) error {
	if connectMap[name] != nil {
		return errors.New("please use a different name")
	}
	e := connectMap[name].DataChannel.Close()
	if e != nil {
		return e
	}
	e = connectMap[name].PeerConnection.Close()
	if e != nil {
		return e
	}
	return nil
}
