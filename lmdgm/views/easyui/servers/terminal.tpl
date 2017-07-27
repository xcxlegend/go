 {{template "../public/header.tpl"}}
 <script type="text/javascript" src="/static/js/term.js"></script>

 <style>body, #terminal_1 {position: absolute; height: 90%; width: 100%; margin: 0px;}</style>
 <!-- <button class="dlog_reload_1"  value="1">刷新</button> -->
 <div id="terminal_1">

 </div>
<script>
$(document).ready(function(){
	 $(".dlog_reload_1").click(function(){ 
        var id= $(this).val();
		close(id); 
		openDilog(id)
     });
	var openWs = function() {
	   var  ws = new WebSocket("ws://{{.host}}/rbac/ssh/ws?id={{.id}}");
       // map.put("1",ws);
	  var term = new Terminal({  
        rows: 21,
        useStyle: true,
        screenKeys: true,
        cursorBlink: false
      });
	 term.on('data', function(data) {
      ws.send(data);
     });
	 term.open(document.getElementById("terminal_1"));
   ws.onopen = function (event) {
      console.log(event);
      
   }
	  ws.onmessage = function(event) {
         term.write(window.atob(event.data));
      };
	   ws.onclose = function(event) {
            if (term) {
                term.destroy(); 
            }
            console.log("close", event);
        };
	}
  
      openWs();
 
});
</script>