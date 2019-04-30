package views

var Str_var =`ws=new WebSocket(location.href.replace('http://','ws://')+'ws');
ws.onopen=function(e){
    console.log('ws open')
}
ws.onmessage=function(e){
    console.log('ws message:'+e.data)
}
ws.onclose=function(e){
    console.log('ws closed')
}

bridge=new Object();
bridge.send=function(msg){
    ws.send(msg);
}`
