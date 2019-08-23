server='{{.Addr}}';

function connectWs(){
    var ws=new WebSocket('ws://'+server+'/fp/ws');
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

{{ range .Funcs}}
function {{.Name}}({{range .Args}}{{.}},{{end}} callback){
    var xhr=new XMLHttpRequest();
    xhr.onreadystatechange=function(e){
        if (this.readyState==4){
            if (this.status==200){
                var obj=JSON.parse(this.responseText);
                callback(obj);
                return;
            }

            console.error(this.responseText);
        }
    };
    xhr.open('POST','/fp/call/{{.Name}}',true);
    var body=[
        {{range .Args}}
        {{.}},
        {{end}}
    ];
    xhr.send(body);
}
{{ end }}