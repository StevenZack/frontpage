server=location.href.split('/')[2];

function connectWs(){
    var ws=new WebSocket('ws://'+server+'/fp/ws');
    ws.onopen=function(e){
        console.log('ws open');
    };
    ws.onmessage=function(e){
        var s=e.data;
        eval(s);
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
        throw xhr.responseText;
    }
    var obj=JSON.parse(xhr.responseText);
    return obj;
}
{{ end }}


connectWs();

