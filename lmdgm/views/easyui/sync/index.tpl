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
    var URL = "/rbac/sync/process?group=";
    $(function() {
        //服务器列表
        $("#datagrid").datagrid({
            title: '应用同步列表',
            url: URL + 'app',
            method: 'POST',
            pagination: true,
            fitColumns: true,
            striped: true,
            rownumbers: false,
            singleSelect: true,
            idField: 'Id',
            pagination: false,
            columns: [
                [
                    /*  {
                         field: 'Id',
                         title: 'ID',
                         width: 50,
                         sortable: true
                     },  */
                    {
                        field: 'ServerCode',
                        title: '标识',
                        width: 100,
                        sortable: true,
                        align: 'center',
                    }, {
                        field: 'Host',
                        title: '服务器地址',
                        width: 50,
                        align: 'center',
                    }, {
                        field: 'DestFile',
                        title: '文件路径',
                        width: 100,
                        align: 'center',
                    }, {
                        field: 'StartTimeFormat',
                        title: '开始时间',
                        width: 80,
                        align: 'center',
                    }, {
                        field: 'FinishTimeFormat',
                        title: '完成时间',
                        width: 80,
                        align: 'center',
                    }, {
                        field: 'md5',
                        title: 'md5',
                        width: 80,
                        align: 'center',
                    }, {
                        field: 'Progress',
                        title: '进度',
                        width: 100,
                        align: 'center',
                        formatter: function(value, row, index) {
                            return value + "%" + ' ('+row.total+' byte)'
                        }
                    } 
                ]
            ],
        }); 

        $("#datagrid_conf").datagrid({
            title: '配置列表',
            url: URL + 'conf',
            method: 'POST',
            pagination: true,
            fitColumns: true,
            striped: true,
            rownumbers: false,
            singleSelect: true,
            idField: 'Id',
            pagination: false,
            columns: [
                [
                    /*  {
                         field: 'Id',
                         title: 'ID',
                         width: 50,
                         sortable: true
                     },  */
                    {
                        field: 'ServerCode',
                        title: '标识',
                        width: 100,
                        sortable: true,
                        align: 'center',
                    }, {
                        field: 'Host',
                        title: '服务器地址',
                        width: 50,
                        align: 'center',
                    }, {
                        field: 'DestFile',
                        title: '文件路径',
                        width: 100,
                        align: 'center',
                    }, {
                        field: 'StartTimeFormat',
                        title: '开始时间',
                        width: 80,
                        align: 'center',
                    }, {
                        field: 'FinishTimeFormat',
                        title: '完成时间',
                        width: 80,
                        align: 'center',
                    }, {
                        field: 'Progress',
                        title: '进度',
                        width: 100,
                        align: 'center',
                        formatter: function(value, row, index) {
                            return (row.total > 0 ? value : 100.000) + "%" + ' ('+row.total+' byte)'
                        } 
                    }, {
                        field: 'md5',
                        title: 'md5',
                        width: 80,
                        align: 'center',
                    }
                ]
            ],
        }); 

    });
    //刷新
    function reloadrow() {
        $("#datagrid").datagrid("reload");
        $("#datagrid_conf").datagrid("reload");
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
 
</script>

<body>
<div id="tb" style="padding:5px;height:auto">
    <a href="#" icon='icon-reload' plain="true" onclick="startReload()" class="easyui-linkbutton" >自动刷新</a>
    <a href="#" icon='icon-cancel' plain="true" onclick="closeReload()" class="easyui-linkbutton" >停止自动刷新</a>
    <a href="#" icon='icon-reload' plain="true" onclick="reloadrow()" class="easyui-linkbutton" >刷新一次</a>
    <a href="/rbac/dir/index" class="easyui-linkbutton">返回文件夹管理</a>
</div>
    <table id="datagrid"></table>
    <br/></br/>
    <table id="datagrid_conf"></table>
    <!--表格内的右键菜单-->
     
    <!--表头的右键菜单-->
</body>
<script type="text/javascript">
    // var timer = setInterval("reloadrow()", 1000)

    function startReload(){
        timer = setInterval("reloadrow()", 1000)
    }

    function closeReload(){
        clearInterval(timer)
    }

</script>
</html>