package views

var Str_var =`server='{{.Addr}}';

function connectWs(){
    var ws=new WebSocket('ws://'+server+'/fp/ws');
    ws.onopen=function(e){
        console.log('ws open');
    };
    ws.onmessage=function(e){
    };
    ws.onclose=function(e){
        console.log('ws closed');
        window.close();
    };
}

{{ range .Funcs}}
function {{.Name}}({{range .Args}}{{.}},{{end}}){
    var xhr=new XMLHttpRequest();
    xhr.open('POST','/fp/call/{{.Name}}',false);
    var body=[
        {{range .Args}}
        {{.}},
        {{end}}
    ];
    xhr.send(JSON.stringify(body));
    if (xhr.status!=200){
        throw new Exception(xhr.responseText);
    }
    var obj=JSON.parse(xhr.responseText);
    return obj;
}
{{ end }}


connectWs();

`
