<!DOCTYPE html>
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>SD GM</title>
    <link rel="stylesheet" type="text/css" href="{{.resource_domain}}/static/easyui/jquery-easyui/themes/default/easyui.css" />
    <link rel="stylesheet" type="text/css" href="{{.resource_domain}}/static/easyui/jquery-easyui/themes/icon.css" />
    <script type="text/javascript" src="{{.resource_domain}}/static/easyui/jquery-easyui/jquery.min.js"></script>
    <script type="text/javascript" src="{{.resource_domain}}/static/easyui/jquery-easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="{{.resource_domain}}/static/easyui/jquery-easyui/common.js"></script>
    <script type="text/javascript" src="{{.resource_domain}}/static/easyui/jquery-easyui/easyui_expand.js"></script>
    <script type="text/javascript" src="{{.resource_domain}}/static/easyui/jquery-easyui/phpjs-min.js"></script>
    <script>
        var PATH = window.location.pathname
        function saveHome() {
            if (confirm("是否将本模块做为登录后模板进入的页面?")){
                $.post('/user/user/save_home', {"path": PATH}, function (data) {
                    alert(data.info);
                })
            }
        }
        window.onload = function () {
            if (PATH == "/public/login" || PATH == "/rbac/servers/terminal"){
                $('#btn-group-fix').hide()
            }
        }
    </script>
</head>
<div id="btn-group-fix" style="margin: 10px 0">
    <div class="btn-group" role="group">
        <button type="button" onclick="saveHome()" class="btn btn-default">存为Home</button>
    </div>
</div>