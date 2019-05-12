package views

var Str_var =`ws = new WebSocket('ws://'+location.href.split('/')[2] + '/ws');
ws.onopen = function (e) {
    console.log('ws open')
}
ws.onmessage = function (e) {
    console.log('ws message:' + e.data)
    var slice = JSON.parse(e.data);
    switch (slice[0]) {
        case 'eval':
            eval(slice[1]);
            break;
        case 'invoke':
            tmp_args=slice;
            var strToInvoke = slice[1] + '(';
            for (var i = 2; i < slice.length; i++) {
                strToInvoke+='tmp_args['+i.toString()+'],'
            }
            strToInvoke += ')';
            console.log('invoking:'+strToInvoke)
            eval(strToInvoke);
            break;
    }
}
ws.onclose = function (e) {
    console.log('ws closed')
}

{{range $key, $value := .}}
{{with $value}}
{{.FnName }}=function({{ range .Ins }} {{.}}, {{ end }}) {
    var args = [
        '{{.FnName}}',
        {{ range .Ins }}
{{.}}, {{end }}
    ];
for (var i = 0; i < args.length; i++) {
    if (args[i] == null) {
        throw 'invoke function : "{{.FnName}}" failed: invalid args';
    }
}
var str = JSON.stringify(args);
ws.send(str);
}
{{end }}
{{end }}`
