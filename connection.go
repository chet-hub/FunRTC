package FunRTC

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chet-hub/FunRTC/utils"
	"github.com/pion/webrtc/v3"
)

type connection struct {
	IsServer                interface{}
	Configuration           webrtc.Configuration
	DataChanelInit          webrtc.DataChannelInit
	PeerConnection          *webrtc.PeerConnection
	DataChannel             *webrtc.DataChannel
	ApplyString             string
	OnConnectionStateChange func(webrtc.PeerConnectionState)
	OnCreateDataChannel     func(*webrtc.DataChannel)
}

func new(OnCreateDataChannelFunc func(*webrtc.DataChannel),
	OnConnectionStateChangeFunc func(webrtc.PeerConnectionState),
	webrtcConfiguration string,
	dataChannelOptions string) (*connection, error) {

	WebRTCConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
					"stun:stun3.l.google.com:19302",
					"stun:stun4.l.google.com:19302",
				},
			},
		},
	}
	if webrtcConfiguration != "" {
		if e := json.Unmarshal([]byte(webrtcConfiguration), &WebRTCConfig); e != nil {
			return nil, e
		}
	}
	DataChanelConfig := webrtc.DataChannelInit{}
	if dataChannelOptions != "" {
		if e := json.Unmarshal([]byte(dataChannelOptions), &DataChanelConfig); e != nil {
			return nil, e
		}
	}
	var Con = connection{
		IsServer:       nil,
		Configuration:  WebRTCConfig,
		DataChanelInit: DataChanelConfig,
		OnConnectionStateChange: func(s webrtc.PeerConnectionState) {
			fmt.Printf("PeerConnectionState changed'%s' \n", s.String())
		},
		OnCreateDataChannel: func(d *webrtc.DataChannel) {
			fmt.Printf("DataChannel New '%s' \n", d.Label())
			d.OnOpen(func() {
				fmt.Printf("DataChannel[%d] OnOpen '%s'\n", *d.ID(), d.Label())
			})
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				fmt.Printf("DataChannel[%d] Message '%s': '%s'\n", *d.ID(), d.Label(), string(msg.Data))
			})
			d.OnClose(func() {
				fmt.Printf("DataChannel[%d] OnClose '%s'\n", *d.ID(), d.Label())
			})
			d.OnError(func(err error) {
				fmt.Printf("DataChannel[%d] OnError '%s' '%s' \n", *d.ID(), d.Label(), err)
			})
		},
	}
	if OnConnectionStateChangeFunc != nil {
		Con.OnConnectionStateChange = OnConnectionStateChangeFunc
	}
	if OnCreateDataChannelFunc != nil {
		Con.OnCreateDataChannel = OnCreateDataChannelFunc
	}
	peerConnection, err := webrtc.NewPeerConnection(Con.Configuration)
	if err != nil {
		return nil, err
	}
	Con.PeerConnection = peerConnection
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		Con.OnConnectionStateChange(s)
	})
	return &Con, nil
}

func (Con *connection) apply() (string, error) {
	if Con.IsServer == true {
		return "", errors.New("can not call apply after calling approve")
	} else if Con.IsServer == false {
		return "", errors.New("can not call apply two times")
	} else {
		Con.IsServer = false
		channel, err := Con.PeerConnection.CreateDataChannel("Client", &Con.DataChanelInit)
		if err != nil {
			return "", err
		}
		Con.DataChannel = channel
		Con.OnCreateDataChannel(channel)
	}
	offer, e := Con.PeerConnection.CreateOffer(nil)
	if e != nil {
		return "", e
	}
	gatherComplete := webrtc.GatheringCompletePromise(Con.PeerConnection)
	if e = Con.PeerConnection.SetLocalDescription(offer); e != nil {
		return "", e
	}
	<-gatherComplete
	var applyString = utils.Encode(Con.PeerConnection.LocalDescription())
	return applyString, nil
}

func (Con *connection) connect(approveString string) error {
	offer := webrtc.SessionDescription{}
	utils.Decode(approveString, &offer)
	if sdpErr := Con.PeerConnection.SetRemoteDescription(offer); sdpErr != nil {
		return sdpErr
	}
	return nil
}

func (Con *connection) approve(applyString string) (string, error) {
	if Con.IsServer == true {
		return "", errors.New("can not call approve two times")
	} else if Con.IsServer == false {
		return "", errors.New("can not call approve after calling apply")
	} else {
		Con.IsServer = true
		Con.PeerConnection.OnDataChannel(func(channel *webrtc.DataChannel) {
			Con.DataChannel = channel
			Con.OnCreateDataChannel(channel)
		})
	}
	offer := webrtc.SessionDescription{}
	utils.Decode(applyString, &offer)
	err := Con.PeerConnection.SetRemoteDescription(offer)
	if err != nil {
		return "", err
	}
	answer, err := Con.PeerConnection.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	gatherComplete := webrtc.GatheringCompletePromise(Con.PeerConnection)
	err = Con.PeerConnection.SetLocalDescription(answer)
	if err != nil {
		return "", err
	}
	<-gatherComplete
	var approveString = utils.Encode(Con.PeerConnection.LocalDescription())
	return approveString, nil
}
