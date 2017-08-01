{{template "../public/header.tpl"}}
<script type="text/javascript" src="/static/easyui/jquery-easyui/ajaxfileupload.js?v=1.2"></script>

<body>
    <div style="padding:5px 0;">
        <a href="conf" class="easyui-linkbutton">配置</a>
        <a href="app" class="easyui-linkbutton">应用</a>
        <!-- <a href="#" class="easyui-linkbutton">Save</a> -->
    </div>
    <div style="margin:20px 0;"></div>
    <table id="tg" title="[配置] 文件历史" class="easyui-treegrid ContextMenu" style="width:700px;height:250px" data-options="
				url: '/rbac/dir/conf',
				method: 'get',
				lines: true,
				rownumbers: true,
				idField: 'full_path_64',
				treeField: 'path',
                fitColumns: true,
                toolbar: '#tb',
                onContextMenu: onContextMenu,
                onLoadSuccess: function(){
                    $('#tg').treegrid('collapseAll')
                }
			">
        <thead>
            <tr>
                <th data-options="field:'path'" width="220">path</th>
                <th data-options="field:'is_dir'" width="220">Dir</th>
                <th data-options="field:'size'" width="100" align="right">Size</th>
                <th data-options="field:'mod_time',formatter:getLocalTime" width="250">Modified Time</th>
                <!-- <th data-options="field:'base',formatter:getTrAction" width="250">Action</th> -->
            </tr>
        </thead>
    </table>
    <div id="file_menu" class="easyui-menu" style="width:120px;">
        <div onclick="show_file()" data-options="iconCls:'icon-add'">Show</div>
    </div>
    <div id="base_menu" class="easyui-menu" style="width:120px;">
        <div onclick="sync()" data-options="iconCls:'icon-add'">Sync</div>
    </div>
    <div id="msgwindow"></div>
    <div id="tb" style="height:auto">
        <a href="javascript:void(0)" id="fileButton" title="请选择.zip压缩文件" class="easyui-linkbutton easyui-tooltip" data-options="iconCls:'icon-add',plain:true" onclick="uploadWindow()">Upload</a>
        <a href="/rbac/sync/index" id="goToSchedule" title="查看同步列表" class="easyui-linkbutton" data-options="iconCls:'icon-search',plain:true">同步列表</a>
        <a href="javascript:void(0)" id="fileButton" title="Sync" class="easyui-linkbutton easyui-tooltip" data-options="iconCls:'icon-add',plain:true" onclick="sync()">Sync</a>
        <input type="file" style="display: none" id="file" name="file" />
    </div>
    <div id="file_alert" title="文件上传" data-options="iconCls:'icon-save'" style="width:300px;height:100px;padding:10px"></div>

    <script type="text/javascript">
        // $(function () {
        $('#fileButton').tooltip({
                position: "right"
            })
            // });

        function uploadWindow() {
            $('#file').click()
        }

        $('#file').change(function() {
            $('#file_alert').dialog({
                content: '文件正在上传, 请勿关闭页面等待上传完成...<br/>如果需要重新上传, 请刷新页面后点击Upload',
            })
            $('#file_alert').dialog('open');
            $.ajaxFileUpload({
                url: '/rbac/dir/upload_conf', //处理图片的脚本路径
                type: 'post', //提交的方式
                secureuri: false, //是否启用安全提交
                fileElementId: 'file', //file控件ID
                dataType: 'json', //服务器返回的数据类型      
                timeout: 0,
                success: function(data, status) { //提交成功后自动执行的处理函数
                    $('#file_alert').dialog({
                        content: data.info,
                    })
                    $("#tg").treegrid('reload')
                    $('#file').val("")
                    return
                },
                error: function(data, status, e) { //提交失败自动执行的处理函数
                    $('#file_alert').dialog({
                        content: "服务器错误",
                    })
                    return;
                }
            })
        });

        function sync() {
            var selectedRow = $("#tg").datagrid('getSelected');
            // console.log(selectedRow)
            $.post('/rbac/sync/sync_conf', {
                dir: selectedRow.path
            }, function(res) {
                if (!res.status) {
                    $.messager.alert('提示', res.info, 'error')
                } else {
                    // $.messager.alert('提示', "开始同步, 请点击同步列表查看",  'info')
                    window.location.href = "/rbac/sync/index"
                }

            })
        }

        function show_file(row) {
            var row = row || $("#tg").datagrid("getSelected");
            $('#msgwindow').dialog({
                href: '/rbac/dir/file_content?file=' + row.full_path,
                width: 600,
                height: 300,
                modal: true,
                cache: false,
                title: row.path,
                closable: true,
                onLoad: function() {
                    var contentnode = $('#msgwindow').find('.panel').children('.dialog-content')
                    contentnode.html("<pre>" + contentnode.html() + "</pre>")
                }
            });
        }


        function onContextMenu(e, row) {
            e.preventDefault();
            $(this).treegrid('select', row.full_path_64);
            // console.log(row, row.is_dir)
            if (row.is_dir && row.base) {
                $('#base_menu').menu('show', {
                    left: e.pageX,
                    top: e.pageY
                });
            } else if (!row.is_dir) {
                $('#file_menu').menu('show', {
                    left: e.pageX,
                    top: e.pageY
                });
            }
        }


        // function getTrAction(val, row) {
        //     console.log(val, row)
        //     if (row.is_dir && row.base){
        //         return '<a onclick="setTimeout(sync(),500)" data-options="iconCls:\'icon-add\'">Sync</a>'
        //     }else if (!row.is_dir){
        //         return '<a onclick="setTimeout(show_file(),500)" data-options="iconCls:\'icon-add\'">Show</a>'
        //     }else{
        //         return '-'
        //     }
        // }

        function getLocalTime(val, row) {
            // console.log(val)
            val = val * 1000
            return formatDate(val)
        }

        function formatDate(timestamp) {
            var now = new Date(timestamp)
            return now;
            // var year=now.getYear(); 
            // var month=now.getMonth()+1; 
            // var date=now.getDate(); 
            // var hour=now.getHours(); 
            // var minute=now.getMinutes(); 
            // var second=now.getSeconds(); 
            // return year+"-"+month+"-"+date+" "+hour+":"+minute+":"+second; 
        }
        Date.prototype.toString = function() {
            return this.getFullYear() +
                "-" + (this.getMonth() > 8 ? (this.getMonth() + 1) : "0" + (this.getMonth() + 1)) +
                "-" + (this.getDate() > 9 ? this.getDate() : "0" + this.getDate()) +
                " " + (this.getHours() > 9 ? this.getHours() : "0" + this.getHours()) +
                ":" + (this.getMinutes() > 9 ? this.getMinutes() : "0" + this.getMinutes()) +
                ":" + (this.getSeconds() > 9 ? this.getSeconds() : "0" + this.getSeconds());
        }
    </script>
</body>