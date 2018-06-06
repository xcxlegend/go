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
    var URL = "/rbac/log";
    $(function() {

        //服务器列表
        $("#datagrid").datagrid({
            title: '日志',
            url: URL + '/index',
            method: 'POST',
            fitColumns: true,
            striped: true,
            // rownumbers: true,
            singleSelect: true,
            idField: 'Id',
            pagination: true,
            pageSize: 20,
            pageList: [10, 20, 30, 50, 100],
            columns: [
                [

                    {
                        field: 'Id',
                        title: 'ID',
                        width: 50,
                        // sortable: true,
                        align: 'center',
                    }, {
                        field: 'NodeTitle',
                        title: '类型',
                        width: 50,
                        // sortable: true,
                        align: 'center',
                    }, {
                        field: 'Uid',
                        title: 'UID',
                        width: 100,
                        align: 'center',
                    }, {
                        field: 'Nickname',
                        title: '名称',
                        width: 50,
                        align: 'center',
                    }, {
                        field: 'Createtime',
                        title: '时间',
                        width: 30,
                        align: 'center',
                        formatter: function(value, row, index) {
                            if (value) return phpjs.date("Y-m-d H:i:s", phpjs.strtotime(value));
                            return value;
                        }
                    }

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

    //刷新
    function reloadrow() {
        $("#datagrid").datagrid("reload");
        editIndex = undefined
    }

    function doSearch() {
        var type = $('#searchType').combobox('getValue')
        $('#datagrid').datagrid('load', {
            type: type,
        })
    }

    function show() {
        var row = $("#datagrid").datagrid("getSelected")
        if (!row) {
            vac.alert("请选择要取消的行");
            return;
        }
        $('#msgwindow').dialog({
                content: '<pre>' + row.Remark + '</pre>',
                width: 600,
                height: 300,
                modal: true,
                cache: false,
                title: "LOG详情",
                closable: true,
            })
            // $('#msgwindow').dialog('open')
    }
</script>

<body>
    <table id="datagrid" toolbar="#tb"></table>
    <br/>

    <div id="tb" style="padding:5px;height:auto">
        <a href="#" icon='icon-reload' plain="true" onclick="reloadrow()" class="easyui-linkbutton">刷新</a> 类型:
        <select id="searchType" class="easyui-combobox" panelHeight="auto" style="width:150px">
                <option value="">all</option>
                {{range $k, $v:=.types}}
                <option value="{{$k}}">{{$v}}</option>
                {{end}}
			</select>
        <a href="javascript:void(0)" class="easyui-linkbutton" iconCls="icon-search" onclick="doSearch()">Search</a>
        <a href="#" icon='icon-search' plain="true" onclick="show()" class="easyui-linkbutton">查看详细</a>
    </div>



    <!--表格内的右键菜单-->
    <div id="mm" class="easyui-menu" style="width:120px;display: none">

        <div iconCls='icon-reload' onclick="reloadrow()">刷新</div>
        <div iconCls='icon-search' onclick="show()">查看详细</div>
        <div class="menu-sep"></div>
    </div>



    <div id="msgwindow" style="padding: 20px;"></div>
</body>

</html>