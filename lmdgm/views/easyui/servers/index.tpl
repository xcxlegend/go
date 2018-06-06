{{template "../public/header.tpl"}}
<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous" />
        <!-- 可选的 Bootstrap 主题文件（一般不用引入） -->
<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous" />
        <!-- 最新的 Bootstrap 核心 JavaScript 文件 -->
<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>

<link rel="stylesheet" type="text/css" href="/static/css/bootstrap/font/font-awesome.css" />
<!--<link rel="stylesheet" type="text/css" href="/static/css/bootstrap/custom-styles.css" />-->
<style>
    *{
        webkit-box-sizing: inherit;
        box-sizing: inherit;
    }
    #datagrid_run li {
        width: auto;
        border: solid 1px #fff;
        height: auto;
        float: left;
        margin-right: 8px;
        margin-bottom: 10px;
        list-style: none;
    }

    .conf-item .panel-body{
        border-color: #fff;
    }

    /*运行中*/
    .conf-run {
        border-color: #5cb85c;
    }

    .conf-run .icon{
        color: #5cb85c;
    }

    .conf-warn {
        border-color: red;
    }

    .conf-warn .icon{
        color: red;
    }

    /*就绪*/
    .conf-ready {
        border-color: rgb(154, 155, 156);
    }

    .conf-ready .icon{
        color: rgb(154, 155, 156);
    }

</style>
<script type="text/javascript">
    var editIndex = undefined;
    var statuslist = [
        {statusid:'1',name:'禁用'},
        {statusid:'2',name:'启用'}
    ];

    var confsList = [
        {key: "all",  name: "All"},
        {key: "logic",  name: "Logic"},
        {key: "auth",   name: "Auth"},
        {key: "pvpmatch",  name: "Match"},
        {key: "pvp",    name: "PVP"},
    ]

    var URL = "/rbac/servers";
    $(function() {

        //服务器列表
        $("#datagrid").datagrid({
            title: '服务器列表',
            url: URL + '/index',
            method: 'POST',
            fitColumns: true,
            striped: true,
            rownumbers: true,
            singleSelect: true,
            idField: 'Id',
            pagination: true,
            pageSize: 20,
            pageList: [10, 20, 30, 50, 100],
            columns: [
                [
                    /*  {
                         field: 'Id',
                         title: 'ID',
                         width: 50,
                         sortable: true
                     },  */
                    {
                        field: 'ServerName',
                        title: '服务器名',
                        width: 100,
                        sortable: true,
                        align: 'center',
                        editor: 'text'
                    }, {
                        field: 'OutHost',
                        title: '服务器地址',
                        width: 100,
                        align: 'center',
                        editor: 'text'
                    }, {
                        field: 'Host',
                        title: '服务器内网传输地址',
                        width: 100,
                        align: 'center',
                        editor: 'text'
                    }, {
                        field: 'Port',
                        title: '端口',
                        width: 30,
                        align: 'center',
                        editor: 'numberbox'
                    }, {
                        field: 'LoginUserName',
                        title: '登录名',
                        width: 100,
                        align: 'center',
                        editor: 'text'
                    }, {
                        field: 'LoginPassword',
                        title: '密码',
                        width: 100,
                        align: 'center',
                        editor: 'text',
                        formatter: function(value, row, index) {
                            value = "******"
                            return "******"
                        }
                    }, {
                        field: 'IsMount',
                        title: '是否挂载',
                        width: 100,
                        align: 'center',
                    }, {
                        field: 'Confs',
                        title: '配置',
                        width: 100,
                        align: 'center',
                        formatter:function(value){
                            var values = value.split(",")
                            var names = []
                            $.each(values, function (i, e) {
                                for(var i=0; i<confsList.length; i++){
                                    if (confsList[i].key == e) {
                                        names.push(confsList[i].name);
                                    }

                                }
                            })
                            return names.join(",")
//                            return value;
                        },
                        editor:{
                            type:'combobox',
                            options:{
                                valueField:'key',
                                textField:'name',
                                multiple: true,
                                data:confsList
                            }
                        }
                    },{
                        field:'Status',
                        title:'状态',
                        width:50,
                        align: 'center',
                        formatter:function(value){
                            for(var i=0; i<statuslist.length; i++){
                                if (statuslist[i].statusid == value) return statuslist[i].name;
                            }
                            return value;
                        },
                        editor:{
                            type:'combobox',
                            options:{
                                valueField:'statusid',
                                textField:'name',
                                data:statuslist,
                                required:true
                            }
                        }
                    }
                    /* , {
                                        field: 'Createtime',
                                        title: '添加时间',
                                        width: 100,
                                        align: 'center',
                                        formatter: function(value, row, index) {
                                            if (value) return phpjs.date("Y-m-d H:i:s", phpjs.strtotime(value));
                                            return value;
                                        }
                                    } */
                ]
            ],
            onAfterEdit: function(index, data, changes) {
                if (vac.isEmpty(changes)) {
                    return;
                }
                changes.Id = data.Id;
                vac.ajax(URL + '/update', changes, 'POST', function(r) {
                    if (!r.status) {
                        vac.alert(r.info);
                    } else {
                        $("#datagrid").datagrid("reload");
                    }
                })
            },
            onDblClickRow: function(index, row) {
                editrow();
            },
            onClickRow: function(index, row) {
//                 console.log(index, row)
                var _row = row
//                $("#datagrid_run").datagrid({
//                    title: '运行情况',
//                    idField: "pid",
//                    url: URL + '/ajax_get_stat?id=' + row.Id + "&_random="+new Date().getTime(),
//                    method: 'POST',
//                    pagination: false,
//                    toolbar: null,
//                    singleSelect: true,
//                    cache: false,
//                    toolbar: "#tb2",
//                    onLoadSuccess: function (data) {
//                        if ( data == undefined || data == "" ){
//                            vac.alert("正在开启应用, 请再等待几秒后读取运行情况")
//                            return
//                        }
//                        console.log(data)
//                    },
//                    rowStyler: function(index, row) {
//                        var confs = _row.Confs.split(",")
////                        console.log(row, _row)
//                        var show = false
//                        $.each(confs, function (i, e) {
////                            console.log(row.confile, row.confile.indexOf(e))
//                            if (row.confile.indexOf(e) == 0 || e == "all"){
//                                show = true
//                                return
//                            }
//                        })
//                        if (!show){
//                            return "display: none";
//                        }
//                        if (!row.is_run) {
//                            return 'color:#b74949;font-weight:bold;';
//                        }
//                    },
//                    columns: [
//                        [{
//                            field: 'server_id',
//                            title: '服务器ID',
//                            width: 200,
//                            align: 'center',
//                        }, {
//                            field: 'pid',
//                            title: 'PID',
//                            width: 200,
//                            align: 'center',
//                        }, {
//                            field: 'confile',
//                            title: '配置文件',
//                            width: 200,
//                            align: 'center',
//                        }, {
//                            field: 'stime',
//                            title: '开始时间',
//                            width: 200,
//                            align: 'center',
//                        }, {
//                            field: 'ltime',
//                            title: '最近心跳时间',
//                            width: 200,
//                            align: 'center',
//                        }, {
//                            field: 'onlineCount',
//                            title: '在线人数',
//                            width: 200,
//                            align: 'center',
//                        }, {
//                            field: 'is_run',
//                            title: 'RUN',
//                            width: 200,
//                            align: 'center'
//                        }, {
//                            field: 'warn',
//                            title: 'Warn',
//                            width: 200,
//                            align: 'center',
//                            formatter:function(value){
//                                return JSON.stringify(value);
//                            },
//                        }]
//                    ],
//                    onRowContextMenu: function(e, index, row) {
//                        e.preventDefault();
//                        $(this).datagrid("selectRow", index);
//                        $('#mm_run').menu('show', {
//                            left: e.clientX,
//                            top: e.clientY
//                        });
//                    },
//                });


                $.post(URL + '/ajax_get_stat?format=html&id=' + row.Id + "&_random="+new Date().getTime(), {}, function(res){
                    $('#datagrid_run').html(res)
                 })
            },
            onRowContextMenu: function(e, index, row) {
                console.log(e, index, row)
                e.preventDefault();
                $(this).datagrid("selectRow", index);
                $('#mm').menu('show', {
                    left: e.clientX,
                    top: e.clientY
                });
            }
        });



        $('#datagrid_online').datagrid({
            title: '在线情况',
            url: URL + '/ajax_get_online',
            method: 'POST',
            fitColumns: true,
            striped: true,
            rownumbers: true,
            singleSelect: true,
            idField: 'Id',
            pagination: true,
            pageSize: 20,
            pageList: [10, 20, 30, 50, 100],
            columns: [
                [
                    {
                        field: 'type',
                        title: '服务器',
                        width: 100,
                        sortable: true,
                        align: 'center',
                    }, {
                        field: 'count',
                        title: '在线人数',
                        width: 100,
                        align: 'center',
                    }
                ]
            ]
        });


        //创建添加服务器窗口
        $("#dialog").dialog({
            modal: true,
            resizable: true,
            top: 150,
            closed: true,
            buttons: [{
                text: '保存',
                iconCls: 'icon-save',
                handler: function() {
                    $("#form1").form('submit', {
                        url: URL + '/add',
                        onSubmit: function() {
                            return $("#form1").form('validate');
                        },
                        success: function(r) {
                            var r = $.parseJSON(r);
                            if (r.status) {
                                $("#dialog").dialog("close");
                                $("#datagrid").datagrid('reload');
                            } else {
                                vac.alert(r.info);
                            }
                        }
                    });
                }
            }, {
                text: '取消',
                iconCls: 'icon-cancel',
                handler: function() {
                    $("#dialog").dialog("close");
                }
            }]
        });

        $("#msgwindow2").dialog({
            modal: true,
            resizable: true,
            top: 150,
            closed: true,
            buttons: [{
                text: '关闭',
                iconCls: 'icon-cancel',
                handler: function() {
                    $("#msgwindow2").dialog("close");
                }
            }]
        });

    })

    function editrow() {
        var index = $("#datagrid").datagrid("getSelected")
        if (!index) {
            vac.alert("请选择要编辑的行");
            return;
        }
        var vacindex = vac.getindex("datagrid"); //
        if (editIndex != vacindex) {
            if (editIndex == undefined) {
                $('#datagrid').datagrid('selectRow', vac.getindex("datagrid"))
                    .datagrid('beginEdit', vac.getindex("datagrid"));
            } else {
                $("#datagrid").datagrid("cancelEdit", editIndex);
                // $('#dg').datagrid('selectRow', editIndex);
                $('#datagrid').datagrid('beginEdit', vac.getindex("datagrid"));
            }
            editIndex = vacindex;
        }
        // $('#datagrid').datagrid('beginEdit', vac.getindex("datagrid"));
    }

    function saverow(index) {
        if (!$("#datagrid").datagrid("getSelected")) {
            vac.alert("请选择要保存的行");
            return;
        }
        $('#datagrid').datagrid('endEdit', vac.getindex("datagrid"));
        editIndex = undefined
    }
    //取消
    function cancelrow() {
        if (!$("#datagrid").datagrid("getSelected")) {
            vac.alert("请选择要取消的行");
            return;
        }
        $("#datagrid").datagrid("cancelEdit", vac.getindex("datagrid"));
        editIndex = undefined
    }
    //刷新
    function reloadrow() {
        $("#datagrid").datagrid("reload");
        editIndex = undefined
    }

    //添加服务器弹窗
    function addrow() {
        $("#dialog").dialog('open');
        $("#form1").form('clear');
    }

    //编辑服务器密码
    // function updateuserpassword() {
    //     var dg = $("#datagrid")
    //     var selectedRow = dg.datagrid('getSelected');
    //     if (selectedRow == null) {
    //         vac.alert("请选择服务器");
    //         return;
    //     }
    //     $("#dialog2").dialog('open');
    //     $("form2").form('load', {
    //         password: ''
    //     });
    // }

    //删除
    function delrow() {
        $.messager.confirm('Confirm', '你确定要删除?', function(r) {
            if (r) {
                var row = $("#datagrid").datagrid("getSelected");
                if (!row) {
                    vac.alert("请选择要删除的行");
                    return;
                }
                vac.ajax(URL + '/del', {
                    Id: row.Id
                }, 'POST', function(r) {
                    if (r.status) {
                        $("#datagrid").datagrid('reload');
                    } else {
                        vac.alert(r.info);
                    }
                })
            }
        });
    }


    function closeApp() {
        // console.log("close app");
        var row = $("#datagrid_run").datagrid("getSelected");
        if (!row) {
            vac.alert("请选择要操作的行");
            return;
        }
        $.post(URL + '/ssh/close', {
            id: row.server_id,
            pid: row.pid,
            confile: row.confile
        }, function(res) {
            $("#datagrid_run").datagrid("reload")
        })
    }

    // function restartApp(){
    //     console.log("close app");
    //     var row = $("#datagrid_run").datagrid("getSelected");
    //     if (!row) {
    //         vac.alert("请选择要操作的行");
    //         return;
    //     }
    //     $.post(URL + '/ssh/restart', {id: row.server_id, pid: row.pid, conf: row.confile}, function(res){
    //         $("#datagrid_run").datagrid("reload")
    //     })
    // }

    function startApp() {
        // console.log("close app");
        var row = $("#datagrid_run").datagrid("getSelected");
        if (!row) {
            vac.alert("请选择要操作的行");
            return;
        }
        $.post(URL + '/ssh/start', {
            id: row.server_id,
            confile: row.confile
        }, function(res) {
            if (!res.status){
                vac.alert(res.info);
                return;
            }
            $("#datagrid_run").datagrid("reload")
//            window.location.reload()
        })
    }

    function mount() {
        var row = $("#datagrid").datagrid("getSelected");
        if (!row) {
            vac.alert("请选择要操作的行");
            return;
        }
        if (!confirm("确定要挂载")) {
            return
        }
        $.post(URL + '/ssh/mount', {
            id: row.Id
        }, function(res) {
            $("#datagrid").datagrid("reload")
        })
    }


    function reload_config_all() {
        $.messager.confirm('Confirm', '是否确定要重载配置?', function (r) {
            if (!r){
                return
            }
            $.post(URL + '/remote/reload_config_all', {}, function (res) {
                var html = ''
                $.each(res, function (k, e) {
                    html += k + ':'+ e + '\n\r\n<br/>'
                })

                $.messager.alert('读取情况', html)
            });
        });
    }



    function editconf() {
        var row = $("#datagrid_run").datagrid("getSelected");
        if (!row) {
            vac.alert("请选择要操作的行");
            return;
        }
        $('#msgwindow').dialog({
            href: URL + '/conf_content?id=' + row.server_id + '&file=' + row.confile,
            width: 600,
            height: 300,
            modal: true,
            cache: false,
            title: row.path,
            closable: true,
            onLoad: function() {
                var contentnode = $('#msgwindow').find('.panel').children('.dialog-content')
                contentnode.html('<textarea id="fcontent" style="width: 90%; height:90%">' + contentnode.html() + "</textarea>")
            },
            buttons: [{
                text: 'Ok',
                iconCls: 'icon-ok',
                handler: function() {
                    $.post(URL + '/update_conf', {
                        id: row.server_id,
                        file: row.confile,
                        content: $('#fcontent').val()
                    }, function(res) {
                        vac.alert("保存成功")
                        $('#msgwindow').dialog('close')
                    });
                    // console.log(row.server_id, $('#fcontent').html());
                }
            }]
        });
    }


    function terminal() {
        var row = $("#datagrid").datagrid("getSelected");
        if (!row) {
            vac.alert("请选择要操作的行");
            return;
        }
        $('#terminal').html('<iframe src="/rbac/servers/terminal?id=' + row.Id + '" width="99%" height="99%" frameborder="0" scrolling="no"></iframe>');
        $('#terminal').window('open');
    }

    var autoRrefreshTimer

    function autorefresh(){
        alert('开始自动刷新')
        autoRrefreshTimer= setInterval(function(){
            $('#datagrid_run').datagrid('reload');
        }, 3000)
    }

    function closeauto(){
        alert('自动刷新关闭')
        clearInterval(autoRrefreshTimer)
    }

    function downlog() {
        var row = $("#datagrid").datagrid("getSelected");
        if (!row) {
            vac.alert("请选择要操作的行");
            return;
        }
        $.post(URL+"/download_log", {server_id: row.Id}, function(res){
            if (!res.status){
                vac.alert(res.info);
                return;
            }else{
                $('#msgwindow2').html("<p>点击下载.  请用tar工具解压 tar -zxvf xxx.tar, 或者用7-zip解压两次</p>" +
                    "<a href='"+res.info+"' target='_blank'>下载</a>")
//                vac.alert("请复制一下链接进行下载: " + res.info)
                $('#msgwindow2').dialog('open');
                return
            }
        })
    }

</script>

<body>
    <table id="datagrid" toolbar="#tb"></table>
    <br/>
    <table id="datagrid_online"></table>
    <br/>
    <div class="panel" title="运行情况" style="width: auto" >
        <div class="panel-header"><div class="panel-title">运行情况</div><div class="panel-tool"></div></div>
        <ul class="panel-body" id="datagrid_run"></ul>
    </div>
    <div id="terminal" class="easyui-window" title="控制台" data-options="modal:true,closed:true" style="width:700px;height:400px;"></div>
    <!-- <div class="easyui-panel" title="Nested Panel" style="width:100%;height:auto;padding:10px;"> -->
    <!-- <div class="easyui-layout"   data-options="fit:true">
            <div data-options="region:'center'" style="; padding:10px">

            </div>
            <div data-options="region:'east'" style="width: 400px;padding:10px">
                
            </div>
        </div> -->
    <!-- </div> -->
    <div id="tb" style="padding:5px;height:auto">
        <a href="#" icon='icon-add' plain="true" onclick="addrow()" class="easyui-linkbutton">新增</a>
        <a href="#" icon='icon-edit' plain="true" onclick="editrow()" class="easyui-linkbutton">编辑</a>
        <a href="#" icon='icon-save' plain="true" onclick="saverow()" class="easyui-linkbutton">保存</a>
        <a href="#" icon='icon-cancel' plain="true" onclick="delrow()" class="easyui-linkbutton">删除</a>
        <a href="#" icon='icon-cancel' plain="true" onclick="cancelrow()" class="easyui-linkbutton">取消</a>
        <a href="#" icon='icon-reload' plain="true" onclick="reloadrow()" class="easyui-linkbutton">刷新</a>
        <a href="#" icon='icon-remove' plain="true" onclick="mount()" class="easyui-linkbutton">挂载</a>
        <a href="#" icon='icon-remove' plain="true" onclick="terminal()" class="easyui-linkbutton">控制台</a>
        <a href="#" icon='icon-reload' plain="true" onclick="reload_config_all()" class="easyui-linkbutton">重新读取配置</a>
    </div>

    <div id="tb2" style="padding:5px;height:auto;display: none;">
        <a href="#" icon='icon-add' plain="true" onclick="startApp()" class="easyui-linkbutton">启动</a>
        <a href="#" icon='icon-cancel' plain="true" onclick="closeApp()" class="easyui-linkbutton">关闭</a>
        <a href="#" icon='icon-reload' plain="true" onclick="$('#datagrid_run').datagrid('reload')" class="easyui-linkbutton">刷新</a>
        <a href="#" icon='icon-reload' plain="true" onclick="autorefresh()" class="easyui-linkbutton">自动刷新</a>
        <a href="#" icon='icon-cancel' plain="true" onclick="closeauto()" class="easyui-linkbutton">关闭自动刷新</a>
        <a href="#" icon='icon-edit' plain="true" onclick="editconf()" class="easyui-linkbutton">编辑文件</a>
    </div>

    <!--表格内的右键菜单-->
    <div id="mm" class="easyui-menu" style="width:120px;display: none">
        <div iconCls='icon-add' onclick="addrow()">新增</div>
        <div iconCls="icon-edit" onclick="editrow()">编辑</div>
        <div iconCls='icon-cancel' onclick="delrow()">删除</div>
        <div iconCls='icon-reload' onclick="reloadrow()">刷新</div>
        <div iconCls='icon-remove' onclick="terminal()">控制台</div>
        <div iconCls='icon-save' onclick="downlog()">下载日志</div>
        <div class="menu-sep"></div>
    </div>
    <div id="mm_run" class="easyui-menu" style="width:120px;display: none">
        <div iconCls='icon-add' onclick="startApp()">启动</div>
        <div iconCls="icon-cancel" onclick="closeApp()">关闭</div>
        <div iconCls="icon-edit" onclick="editconf()">编辑文件</div>
        <div class="menu-sep"></div>
    </div>
    <!--表头的右键菜单-->
    <div id="mm1" class="easyui-menu" style="width:120px;display: none">
        <div icon='icon-add' onclick="addrow()">新增</div>
    </div>
    <div id="dialog" title="添加服务器" style="width:400px;height:400px;">
        <div style="padding:20px 20px 40px 80px;">
            <form id="form1" method="post">
                <table>
                    <tr>
                        <td>服务名：</td>
                        <td><input name="ServerName" class="easyui-validatebox" required="true" /></td>
                    </tr>
                    <tr>
                        <td>服务器外网地址：</td>
                        <td><input name="OutHost" class="easyui-validatebox" required="true" /></td>
                    </tr>
                    <tr>
                        <td>服务器内网地址：</td>
                        <td><input name="Host" class="easyui-validatebox" /></td>
                    </tr>
                    <tr>
                        <td>SSH端口：</td>
                        <td><input name="Port" class="easyui-validatebox" required="true" placeholder="22" /></td>
                    </tr>
                    <tr>
                        <td>SSH账号：</td>
                        <td><input name="LoginUserName" class="easyui-validatebox" required="true" placeholder="root" /></td>
                    </tr>
                    <tr>
                        <td>SSH密码：</td>
                        <td><input name="LoginPassword" class="easyui-validatebox" required="true" /></td>
                    </tr>
                    <tr>
                        <td>Confs：</td>
                        <td><input class="easyui-combobox"
                                   name="Confs"
                                   data-options="
					data: confsList,
					valueField:'key',
					textField:'name',
					multiple:true,
					panelHeight:'auto'
			"></td>
                    </tr>
                    <tr>
                        <td>备注：</td>
                        <td><textarea name="Remark" class="easyui-validatebox" validType="length[0,200]"></textarea></td>
                    </tr>
                </table>
            </form>
        </div>
    </div>
    <div id="msgwindow"></div>
    <div id="msgwindow2" style="width:400px;height:200px;"></div>
</body>

</html>