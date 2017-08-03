{{template "../public/header.tpl"}}

<script type="text/javascript">
    var editIndex = undefined;
    var statuslist = [{
        statusid: '1',
        name: '禁用'
    }, {
        statusid: '2',
        name: '启用'
    }];
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
                // console.log(index, row)
                $("#datagrid_run").datagrid({
                    title: '运行情况',
                    idField: "pid",
                    url: URL + '/ajax_get_stat?id=' + row.Id,
                    method: 'POST',
                    pagination: false,
                    toolbar: null,
                    singleSelect: true,
                    toolbar: "#tb2",
                    rowStyler: function(index, row) {
                        if (!row.is_run) {
                            return 'color:#b74949;font-weight:bold;';
                        }
                    },
                    columns: [
                        [{
                            field: 'server_id',
                            title: '服务器ID',
                            width: 200,
                            align: 'center',
                        }, {
                            field: 'pid',
                            title: 'PID',
                            width: 200,
                            align: 'center',
                        }, {
                            field: 'confile',
                            title: '配置文件',
                            width: 200,
                            align: 'center',
                        }, {
                            field: 'stime',
                            title: '开始时间',
                            width: 200,
                            align: 'center',
                        }, {
                            field: 'ltime',
                            title: '最近心跳时间',
                            width: 200,
                            align: 'center',
                        }, {
                            field: 'is_run',
                            title: 'RUN',
                            width: 200,
                            align: 'center'
                        }]
                    ],
                    onRowContextMenu: function(e, index, row) {
                        e.preventDefault();
                        $(this).datagrid("selectRow", index);
                        $('#mm_run').menu('show', {
                            left: e.clientX,
                            top: e.clientY
                        });
                    },
                });
                // $.post(URL+'/ajax_get_stat', {id: row.Id}, function(res){
                // console.log(res)
                // console.log(res.data)
                // $("#datagrid_run").datagrid("loadData", {"total": 2, "rows":[{
                //       "comd": "./union",
                //       "confile": "conf/test.json",
                //       "pid": "24128"
                //     }, {
                //       "comd": "./union",
                //       "confile": "conf/pvp.json",
                //       "pid": "24140"
                //     }]}
                //     )
                // })
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
        //创建修改密码窗口
        // $("#dialog2").dialog({
        //     modal: true,
        //     resizable: true,
        //     top: 150,
        //     closed: true,
        //     buttons: [{
        //         text: '保存',
        //         iconCls: 'icon-save',
        //         handler: function() {
        //             var selectedRow = $("#datagrid").datagrid('getSelected');
        //             var password = $('#password').val();
        //             vac.ajax(URL + '/UpdateUser', {
        //                 Id: selectedRow.Id,
        //                 Password: password
        //             }, 'post', function(r) {
        //                 if (r.status) {
        //                     $("#dialog2").dialog("close");
        //                 } else {
        //                     vac.alert(r.info);
        //                 }
        //             })
        //         }
        //     }, {
        //         text: '取消',
        //         iconCls: 'icon-cancel',
        //         handler: function() {
        //             $("#dialog2").dialog("close");
        //         }
        //     }]
        // });

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
            $("#datagrid_run").datagrid("reload")
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
</script>

<body>
    <table id="datagrid" toolbar="#tb"></table>
    <br/>
    <table id="datagrid_run"></table>
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
    </div>

    <div id="tb2" style="padding:5px;height:auto;display: none;">
        <a href="#" icon='icon-add' plain="true" onclick="startApp()" class="easyui-linkbutton">启动</a>
        <a href="#" icon='icon-cancel' plain="true" onclick="closeApp()" class="easyui-linkbutton">关闭</a>
        <a href="#" icon='icon-reload' plain="true" onclick="$('#datagrid_run').datagrid('reload')" class="easyui-linkbutton">刷新</a>
        <a href="#" icon='icon-edit' plain="true" onclick="editconf()" class="easyui-linkbutton">编辑文件</a>
    </div>

    <!--表格内的右键菜单-->
    <div id="mm" class="easyui-menu" style="width:120px;display: none">
        <div iconCls='icon-add' onclick="addrow()">新增</div>
        <div iconCls="icon-edit" onclick="editrow()">编辑</div>
        <div iconCls='icon-cancel' onclick="delrow()">删除</div>
        <div iconCls='icon-reload' onclick="reloadrow()">刷新</div>
        <div iconCls='icon-remove' onclick="terminal()">控制台</div>
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
                        <td>备注：</td>
                        <td><textarea name="Remark" class="easyui-validatebox" validType="length[0,200]"></textarea></td>
                    </tr>
                </table>
            </form>
        </div>
    </div>
    <div id="msgwindow"></div>
</body>

</html>