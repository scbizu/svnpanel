//匿名登录事件
$('input[name="anycheckbox"]').on('switchChange.bootstrapSwitch', function(event, state) {
        var repo = $('#reponame').text();
    if (state){
      $.ajax({
        url:"/edit",
        dataType:"json",
        data:{reponame:repo,tag:"anon-access",action:"remark"},
        type:"PUT",
      }).done(function(data,textStatus,jqXHR){
            if (jqXHR.status==200){
              console.log("remarked anon-access");
            }
      });
    }else{
      $.ajax({
        url:"/edit",
        dataType:"json",
        data:{reponame:repo,tag:"anon-access",action:"rmremark"},
        type:"PUT",
      }).done(function(data,textStatus,jqXHR){
        if (jqXHR.status==200){
        console.log("rm remark of  anon-access");
      }
      });
    }
});
//设置权限事件
$('input[name="authzcheckbox"]').on('switchChange.bootstrapSwitch', function(event, state) {
    var repo = $('#reponame').text();
    if(state){
      $.ajax({
        url:"/edit" ,
        dataType:"json",
        data:{reponame:repo,tag:"authz-db",action:"rmremark"},
        type:"PUT",
      }).done(function(data,textStatus,jqXHR){
        if (jqXHR.status==200){
          var dataObj=jQuery.parseJSON(data)
          // console.log(data)
          $('.authz').attr("style","display:block");
          $('.authz-name').attr("style","display:inline");
          $('.adb').text(dataObj.Adb);
            console.log("remarked authz-db");
        }
      });
    }else{
      $.ajax({
        url:"/edit",
        dataType:"json",
        data:{reponame:repo,tag:"authz-db",action:"remark"},
        type:"PUT",
      }).done(function(data,textStatus,jqXHR){
        if (textStatus === "success"){
            $('.authz').attr("style","display:none");
            $('.authz-name').attr("style","display:none");
            console.log("rmed the remark of authz-db...");
        }
      });
    }
});
//用户事件
$('input[name="passwdcheckbox"]').on('switchChange.bootstrapSwitch', function(event, state) {
          var repo = $('#reponame').text();
    if(state){
      $.ajax({
        url:"/edit",
        dataType:"json",
        data:{reponame:repo,tag:"password-db",action:"rmremark"},
        type:"PUT",
      }).done(function(data,textStatus,jqXHR){
        if (textStatus === "success"){
          var dataObj=jQuery.parseJSON(data)
          $('.passwd').attr("style","display:block");
          $('.passwd-name').attr("style","display:inline");
          $('.pdb').text(dataObj.Pdb)
        }
      });
    }else{
      $.ajax({
        url:"/edit",
        dataType:"json",
        data:{reponame:repo,tag:"password-db",action:"remark"},
        type:"PUT",
      }).done(function(data,textStatus,jqXHR){
        if (textStatus === "success"){
            console.log("rmed the remark of password-db...");
            $('.passwd').attr("style","display:none");
            $('.passwd-name').attr("style","display:none");
        }
      });
    }
});

$(function() {
  $('[name="anycheckbox"]').bootstrapSwitch();
});
$(function(){
  $('[name="passwdcheckbox"]').bootstrapSwitch();
});
$(function(){
  $('[name="authzcheckbox"]').bootstrapSwitch();
});
