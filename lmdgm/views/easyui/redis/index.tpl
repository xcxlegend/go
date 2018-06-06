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
    var URL = "/rbac/redis";
    $(function() {

        //服务器列表
        $("#datagrid").datagrid({
            title: 'redis数据库列表',
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
                    {
                         field: 'Index',
                         title: '数据库序号',
                         width: 50,
                     },
                    /* {
                        field: 'Name',
                        title: '服务名',
                        width: 100,
                        sortable: true,
                        align: 'center',
                        editor: 'text'
                    },  */
                    {
                        field: 'Host',
                        title: '服务地址',
                        width: 100,
                        align: 'center',
                        editor: 'text'
                    }, {
                        field: 'Port',
                        title: '端口',
                        width: 30,
                        align: 'center',
                        editor: 'numberbox'
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

    function addrow() {
        $("#dialog").dialog('open');
        $("#form1").form('clear');
        $('#form1').form('load', {
            Port: 6379
        })
    }


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
</script>

<body>
    <table id="datagrid" toolbar="#tb"></table>
    <br/>
    <table id="datagrid_run"></table>

    <div id="tb" style="padding:5px;height:auto">
        <a href="#" icon='icon-add' plain="true" onclick="addrow()" class="easyui-linkbutton">新增</a>
        <a href="#" icon='icon-edit' plain="true" onclick="editrow()" class="easyui-linkbutton">编辑</a>
        <a href="#" icon='icon-save' plain="true" onclick="saverow()" class="easyui-linkbutton">保存</a>
        <!--<a href="#" icon='icon-cancel' plain="true" onclick="delrow()" class="easyui-linkbutton">删除</a>-->
        <a href="#" icon='icon-cancel' plain="true" onclick="cancelrow()" class="easyui-linkbutton">取消</a>
        <a href="#" icon='icon-reload' plain="true" onclick="reloadrow()" class="easyui-linkbutton">刷新</a>
    </div>



    <!--表格内的右键菜单-->
    <div id="mm" class="easyui-menu" style="width:120px;display: none">
        <div iconCls='icon-add' onclick="addrow()">新增</div>
        <div iconCls="icon-edit" onclick="editrow()">编辑</div>
        <!--<div iconCls='icon-cancel' onclick="delrow()">删除</div>-->
        <div iconCls='icon-reload' onclick="reloadrow()">刷新</div>
        <div class="menu-sep"></div>
    </div>


    <div id="dialog" title="添加Redis数据库" style="width:400px;height:400px;">
        <div style="padding:20px 20px 40px 80px;">
            <form id="form1" method="post">
                <table>
                    <!-- <tr>
                        <td>服务名：</td>
                        <td><input name="Name" class="easyui-validatebox" required="true" /></td>
                    </tr> -->
                    <tr>
                        <td>服务地址：</td>
                        <td><input name="Host" class="easyui-validatebox" required="true" /></td>
                    </tr>
                    <tr>
                        <td>端口：</td>
                        <td><input name="Port" class="easyui-validatebox" placeholder="6379" value="6379" /></td>
                    </tr>
                    <tr>
                        <td>序号：</td>
                        <td><input name="Index" class="easyui-validatebox" placeholder="0" value="0" /></td>
                    </tr>
                    <tr>
                        <td></td>
                        <td>（主数据库必须设置为0. 而且不能重复, 不能修改）</td>
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