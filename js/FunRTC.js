
function Connection(config, onCreateDataChannelFunction, onConnectionstatechangeFunction) {
    let Configuration = (config === undefined) ? {
        iceServers: [
            {urls: "stun:stun.l.google.com:19302"},
            {urls: "stun:stun1.l.google.com:19302"},
            {urls: "stun:stun2.l.google.com:19302"},
            {urls: "stun:stun3.l.google.com:19302"},
            {urls: "stun:stun4.l.google.com:19302"},
        ]
    } : config
    let PeerConnection = null
    let ApplyStringResolve = null
    let ApplyString = null
    let DataChannel = null
    let IsServer = null
    let OnConnectionstatechange = function (event) {
        if (onConnectionstatechangeFunction instanceof Function) {
            onConnectionstatechangeFunction(event)
        } else {
            console.log("connection state : " + event.currentTarget.connectionState)
        }
    }
    let OnCreateDataChannel = function (connection, channel) {
        DataChannel = channel
        if (onCreateDataChannelFunction instanceof Function) {
            onCreateDataChannelFunction(connection, channel)
        } else {
            channel.onclose = () => console.log('channel has closed')
            channel.onopen = () => {
                console.log('channel opened')
                setInterval(() => {
                    channel.send(Math.random())
                }, 2000);
            }
            channel.onmessage = (e) => {
                console.log(`Message from '${channel.id}' : '${e.data}'`)
            }
        }
    }

    PeerConnection = new RTCPeerConnection(Configuration);
    PeerConnection.onconnectionstatechange = OnConnectionstatechange
    PeerConnection.onnegotiationneeded = event => {
        PeerConnection.createOffer().then(d => PeerConnection.setLocalDescription(d)).catch(console.error)
    }
    PeerConnection.onicecandidate = event => {
        if (event.candidate === null) {
            ApplyString = btoa(JSON.stringify(PeerConnection.localDescription));
            if (ApplyStringResolve !== null) {
                ApplyStringResolve(ApplyString)
            }
        }
    }
    PeerConnection.ondatachannel = (event) => {
        DataChannel = event.channel
        OnCreateDataChannel(PeerConnection, event.channel)
    }

    return {
        Apply: () => {
            if (IsServer === null) {
                IsServer = false
            } else if (IsServer === true) {
                return Promise.reject('If you have invoked approve method, you can not invoke approve method to as a client')
            }
            DataChannel = PeerConnection.createDataChannel('client')
            OnCreateDataChannel(PeerConnection, DataChannel)
            if (ApplyString !== null) {
                return Promise.resolve(ApplyString)
            } else {
                return new Promise((resolve) => {
                    ApplyStringResolve = resolve
                })
            }
        },
        Approve: (RemoteApplyString) => {
            if (IsServer == null) {
                IsServer = true
            } else if (IsServer === false) {
                return Promise.reject('If you have invoked apply method, you can not invoke approve method to as a server')
            }
            return PeerConnection.setRemoteDescription(JSON.parse(atob(RemoteApplyString))).then(() => {
                return PeerConnection.createAnswer()
            }).then((answer) => {
                return PeerConnection.setLocalDescription(answer)
            }).then(() => {
                return Promise.resolve(btoa(JSON.stringify(PeerConnection.localDescription)))
            })
        },
        Connect: (RemoteApproveString) => {
            return PeerConnection.setRemoteDescription(JSON.parse(atob(RemoteApproveString)))
        },
        Send: (data) => {
            DataChannel.send(data)
        },
        Close: () => {
            DataChannel.close()
        },
    }
}

const connections = {}

const newConnection = (name, configuration, onCreateDataChannelFun, onConnectionstatechangeFun) => {
    if (name === undefined || name.toString().trim() === "") {
        throw "please specify name for this connection"
    }
    const con = new Connection(configuration, onCreateDataChannelFun, onConnectionstatechangeFun)
    return con
}

const FunRTC = {
    ToConnect: (name, configuration, onCreateDataChannelFun, onConnectionstatechangeFun) => {
        const con = newConnection(name, configuration, onCreateDataChannelFun, onConnectionstatechangeFun)
        connections[name] = con
        return con.Apply()
    },
    Accept: (name, remoteApplyString, configuration, onCreateDataChannelFun, onConnectionstatechangeFun) => {
        let con = newConnection(name, configuration, onCreateDataChannelFun, onConnectionstatechangeFun)
        connections[name] = con
        return con.Approve(remoteApplyString)
    },
    DoConnect: (name,remoteApproveString) => {
        return connections[name].Connect(remoteApproveString)
    },
    /**
     send(data: string): void;
     send(data: Blob): void;
     send(data: ArrayBuffer): void;
     send(data: ArrayBufferView): void;
     * @param data
     */
    Send: (name,data) => {
        connections[name].Send(data)
    },
    Close: (name) => {
        connections[name].Close()
    }
}





