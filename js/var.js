server='{{.Addr}}';

function connectWs(){
    var ws=new WebSocket('ws://'+server+'/ws');
    ws.onopen=function(e){
        console.log('ws open');
    };
    ws.onmessage=function(e){
        setTimeout('window.close()',1000);
    };
    ws.onclose=function(e){
        console.log('ws closed');
    };
}


