{{template "../public/header.tpl"}}
<script type="text/javascript">
$(function(){


var URL="/rbac"

var BASE_DIR = "{{.base_dir}}"

var loadTree = function(node, sub){
    var path = node.data("path")
    var level = parseInt(node.data('level'))
    $.post(URL + "/upload/dir", {path: path}, function(res){
        // console.log(res)
        if (res.status == -1){
            alert(res.info);
            return
        }
        var html = '';
        if (sub){
            html += "<ul>"
        }
        var indexhtml = ""
        if (level > 0){
            for (var i = 0; i < level; i++) {
                indexhtml += '<span class="tree-indent"></span>'
            }
        }
        $.each(res.data, function(i, e){
            if (e.is_dir){
                var classname = "tree-icon tree-folder"
                var is_dir = 1
            }else{
                var classname = "tree-icon tree-file"
                var is_dir = 0
            }
            
            var shtml = '<li class="tree-no-li"><div data-loaded="false" data-level="'+(level+1)+'" data-size="'+e.size+'" data-open="false" data-isdir="'+is_dir+'" data-path="'+path+'/'+e.path+'" class="tree-node" style="cursor: pointer;">'+indexhtml+'<span class="'+classname+'"></span><span class="tree-title">'+e.path+'</span></div></li>';
            // console.log(html);
            html += shtml
        }) 
        if (sub){
            html += "</ul>"
            $(node).parent('li').append(html);
        }else{
            $(node).append(html);
        }
        $('.tree-node').unbind('click');
        $('.tree-node').on('click', function(){
            $('.tree-node').removeClass("tree-node-selected")
            $(this).addClass("tree-node-selected")
            if ($(this).data('isdir') == 1){
                console.log($(this).data('path'))
                $('#form-path').val($(this).data('path'))
            }else{
                $('#form-file').val($(this).data('path'))
            }
            if ($(this).data('isdir') == 1 && !$(this).data('loaded')){
                var _this = $(this)
                loadTree(_this, true);
                $(this).data('loaded', true);
                $(this).data('open', true);
                $(this).find('.tree-folder').addClass("tree-folder-open")
                return
            }
            if ($(this).data('open')){
                $(this).siblings('ul').hide()
                $(this).find('.tree-folder').removeClass("tree-folder-open")
                $(this).data('open', false);
            }else{
                $(this).find('.tree-folder').addClass("tree-folder-open")
                $(this).siblings('ul').show()
                $(this).data('open', true);
            }
        })
    })

} 



loadTree($('#tree'), false);

})

function downfile(){
    var URL="/rbac"
    var path = $('#form-file').val();
    if (path == ""){
        alert('未选择文件')
        return
    }
    window.location.href = URL + '/upload/down?path='+path;
}

</script>

<style>
.ht_nav {
    float: left;
    overflow: hidden;
    padding: 0 0 0 10px;
    margin: 0;
}
.ht_nav li{
    font:700 16px/2.5 'microsoft yahei';
    float: left;
    list-style-type: none;
    margin-right: 10px;

}
.ht_nav li a{
    text-decoration: none;
    color:#333;
}
.ht_nav li a.current, .ht_nav li a:hover{
    color:#F20;

}
</style>
<body class="easyui-layout" style="text-align:left">
<div region="west" border="false" split="true" title="目录"  tools="#toolbar" style="width:200px;padding:5px;">
    <ul id="tree" class="tree" data-level="0" data-path="">
        <li class="tree-no-li"><div data-isdir="true" data-path="" class="tree-node" style="cursor: pointer;"><span class="tree-icon tree-folder"></span><span class="tree-title">./</span></div></li>

    </ul>
</div>
<div region="center" border="false" >
    <div id="tabs" >
<div class="easyui-panel" title="单文件上传" style="width:600px">
        <div style="padding:10px 60px 20px 60px">
        <form id="upload" method="post" action="/rbac/upload/post" enctype="multipart/form-data">
            <table cellpadding="5">
                <tr>
                    <td>选择目录:</td>
                    <td><input readonly="" id="form-path" class="easyui-textbox" type="text" name="path" data-options="required:true"></input></td>
                </tr>
                
                <tr>
                    <td>上传文件:</td>
                    <td><input type="file" class="easyui-textbox" name="file" data-options="required:true"></input></td>
                </tr>

                <tr>
                    <td>选择文件:</td>
                   <td><input readonly="" id="form-file" class="easyui-textbox" type="text" name="file" data-options="required:true"></input></td>
                </tr>
                 
            </table>
        </form>
        <div style="text-align:center;padding:5px">
            <a href="javascript:void(0)" class="easyui-linkbutton" onclick="$('#upload').form('submit')">上传</a>
            <a href="javascript:void(0)" class="easyui-linkbutton" onclick="downfile()">下载</a>
        </div>
         
        </div>
    </div>
    </div>
</div>
<div id="toolbar">
</div>
</body>
</html>