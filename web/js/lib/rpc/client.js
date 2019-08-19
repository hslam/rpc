var rpc=function (){};
rpc.Dial = function(address) {
    return new rpc.Client(address)
}
rpc.Client = function(address) {
    this.address=address
    this.version=1.0
    this.client_id=0
    this.url="http://"+address+"/"
    this.Call=function (name, args) {
        var args_bytes=rpc.ArgsEncode(args)
        var req_bytes=rpc.RequestEncode(0,name,args_bytes,false)
        var msg_bytes=rpc.MsgEncode(this.version, this.client_id,proto.pb.MsgType.REQ,false,proto.pb.CodecType.JSON,req_bytes)
        var result_bytes=this.RemoteCall(msg_bytes)
        var msg=rpc.MsgDecode(result_bytes)
        var res=rpc.ResponseDecode(msg.getData())
        var reply=rpc.ReplyDecode(res.getData())
        return reply
    }
    this.RemoteCall=function (msg_bytes){
        var result_bytes;
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange=function() {
            if(xhr.readyState==4&&xhr.status==200) {
                result_bytes=rpc.StringToUint8Array(xhr.responseText)
            }
        }
        xhr.open("POST", this.url, false);
        xhr.send(msg_bytes);
        return result_bytes
    }
    this.setClientId=function (client_id){
        this.client_id=client_id
    }
}

rpc.ArgsEncode= function(args){
    var str=JSON.stringify(args);
    var args_bytes=rpc.StringToUint8Array(str);
    return args_bytes
}

rpc.ReplyDecode= function (reply_bytes) {
    var str = rpc.Uint8ArrayToString(reply_bytes)
    var reply = JSON.parse(str);
    return reply
}

rpc.RequestEncode= function (id,method,args_bytes,noResponse){
    var req = new proto.pb.Request();
    req.setId(id)
    req.setMethod(method)
    req.setData(args_bytes)
    req.setNoresponse(noResponse)
    return req.serializeBinary()
}

rpc.ResponseDecode= function (res_bytes) {
    var res = new proto.pb.Response.deserializeBinary(res_bytes);
    return res
}

rpc.MsgEncode= function (version,id,msgType,batch, codecType,req_bytes) {
    var msg = new proto.pb.Msg();
    msg.setVersion(version)
    msg.setId(id)
    msg.setMsgtype(msgType)
    msg.setBatch(batch)
    msg.setCodectype(codecType)
    msg.setData(req_bytes)
    return msg.serializeBinary()
}

rpc.MsgDecode= function (msg_bytes) {
    var msg = new proto.pb.Msg.deserializeBinary(msg_bytes);
    return msg
}

rpc.Uint8ArrayToString= function (bytes){
    var dataString = "";
    for (var i = 0; i < bytes.length; i++) {
        dataString += String.fromCharCode(bytes[i]);
    }
    return dataString
}

rpc.StringToUint8Array= function (str){
    var arr = [];
    for (var i = 0, j = str.length; i < j; ++i) {
        arr.push(str.charCodeAt(i));
    }
    var tmpUint8Array = new Uint8Array(arr);
    return tmpUint8Array
}
