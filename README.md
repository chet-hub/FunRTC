![loading gif](https://raw.githubusercontent.com/chet-hub/chet-hub.github.io/main/FunRTC/FunRTC.gif)

# What

`Three lines codes to communicate with people in the way of P2P - Having fun to Play with WebRTC`

```javascript
//Tom wants to communicate with Lily, so he calls ToConnect to generate a applying letter
let applying_letter = await FunRTC.ToConnect("Lily")
                        //Lily know Tom want to communicate with her, and get the applying 
                        //letter from Tom in her Messenger or Email. She call the Accept  
                        //method to approve Tom's applying and get the approve letter
                        let approval_letter = await FunRTC.Accept("Tom",applying_letter)
//Tom gets the approval letter from Lily and calls DoConnect to establish a connection to Lily
let success = await FunRTC.DoConnect("Lily",approval_letter)

//Now, they can text each other directly without any middle service
//Tom send a message to Lily
FunRTC.Send("Lily","Hi Lily, how are you.")
                        //Lily got the message and reply 
                        FunRTC.Send("Tom","Hi Tom.")

//....
//Tom think another frend named Olivia
let applying_letter_for_Olivia = await FunRTC.ToConnect("Olivia")
//That's another story...
```

# Why

The most amazon feature of WebRTC is P2P - users can connect each other directly.
However, the p2p connection required the help of signaling server to set up - the p2p
users need to exchange their SDP, icecandidate and other information. This project, based
the webrtc specification, abstracts the interfaces to achieve pure p2p.

# How

* Five methods for Javascript FunRTC

```javascript
FunRTC.ToConnect(name, onCreateDataChannelFun, onConnectionstatechangeFun, configuration)
FunRTC.Accept(name, onCreateDataChannelFun, onConnectionstatechangeFun, configuration)
FunRTC.DoConnect(name,remoteApproveString)
FunRTC.Send(name,data)
FunRTC.Close(name)
```

* These are all methods for golang FunRTC

```golang
func ToConnect(name,applyStr,OnCreateDataChannelFunc,OnConnectionStateChangeFunc,webrtcConfiguration,dataChannelOptions) (string, error)
func Accept(name,applyStr,OnCreateDataChannelFunc,OnConnectionStateChangeFunc,webrtcConfiguration,dataChannelOptions) (string, error)
func DoConnect(name, remoteApprove) error
func SendText(name, text) error
func Send(name string, data []byte) error
func Close(name string) error
```

# Play the examples
`check out the html - https://chet-hub.github.io/FunRTC/`

* Example - Js connect js in one page
```javascript
// Copy and paste the js codes in your Chrome console
let applyPaper = await FunRTC.ToConnect("B")
let approvePaper = await FunRTC.Accept("A",applyPaper)
await FunRTC.DoConnect("B",approvePaper)
applyPaper = await FunRTC.ToConnect("C")
approvePaper = await FunRTC.Accept("B",applyPaper)
await FunRTC.DoConnect("C",approvePaper)
//after the connect established
FunRTC.Send("B","Hi B")
FunRTC.Send("A","Hi A")
FunRTC.Send("C","Hi C")
```

* Example - Js connect golang
```shell
git clone https://github.com/chet-hub/FunRTC
cd FunRTC
go test -run TestConnection_as_Client
```
Then
```javascript
// Copy and paste the js codes in your Chrome console
let applyPaper = await FunRTC.ToConnect("Golang server")
const approvePaper = await PostData('//localhost:8080/msg', applyPaper)
await FunRTC.DoConnect("Golang server", approvePaper)
```

* Example - golang connect js

```shell
git clone https://github.com/chet-hub/FunRTC
cd FunRTC
go test -run TestConnection_as_Server
```
Then
```javascript
// Copy and paste the js codes in your Chrome console
const applyPaper = await PostData('//localhost:8080/msg', "getApplyPaper")
const approvePaper = await FunRTC.Accept("golang client",applyPaper)
const connect = await PostData('//localhost:8080/msg', approvePaper)
```


# Todo
- More test to ensure connection stability and error handling
- Add other language binding for FunRTC
- Add more examples

