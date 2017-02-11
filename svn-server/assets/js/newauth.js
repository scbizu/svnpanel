$("#new_auth_btn").on('click',function(){
  var repo = $('#reponame').text();
  $("#new_auth").css("display","none")
  $("#new_auth").before(`<tr id="new_auth_tr">
          <td style="width:25%">
            <div class="input-group" >
              <span class="input-group-addon">`+repo+`:</span><input type="text" class="form-control" id="new_cpath" value="">
            </div>
          </td>
          <td><!--input type="text" class="form-control" id="new_cuser" value=""--><select multiple="multiple" size="2" class="form-control" id="new_cuser" for=""><option value='*'>*</option></select></td>

          <td><input type="text" class="form-control" id="new_cauth" value=""></td>
          <td>
          <button type="button" id="auth_new_ok" class="btn btn-success " style="display:inline"><span class="glyphicon glyphicon-ok"></span></button>
          &nbsp;
          <button id="auth_new_del" type="button" class="btn btn-danger "><span class="glyphicon glyphicon-trash"></span></button>
          </td>
        </tr>
  `)
    // load data from users column
    $.ajax({
      url:"/users",
      data:{repo:repo},
    }).done(function(data,textStatus,jqXHR){
      var obj = jQuery.parseJSON(data)
      for (key in obj) {
        // console.log(key)
        $("#new_cuser").append("<option value='"+key+"'>"+key+"</option>")
      }

    })
    // load data from groups column
    $.ajax({
      url:"/groups",
      data:{repo:repo},
    }).done(function(data,textStatus,jqXHR){
      var obj = jQuery.parseJSON(data)
      for (key in obj){
        $("#new_cuser").append(`<option value="@`+key+`">@`+key+`</option>`)
      }
    })

  $("#auth_new_del").on('click',function(){
    $("#new_auth_tr").remove();
    $("#new_auth").css("display","block");
  });

  $("#auth_new_ok").on('click',function(){
     var index=parseInt($("#addindex").attr("index"),10)+1
     var newpath = repo+":"+$("#new_cpath").val();
     var userarray = $("#new_cuser").val();
     var newuser="";
     userarray.forEach(function(user,index){
       if(index == userarray.length-1){
         newuser+=user
       }else{
         newuser+=user+","
       }
     });
     var newauth = $("#new_cauth").val();
     $.ajax({
       url:"/addauth" ,
       data:{tag:newpath,new_users:newuser,new_auth:newauth,reponame:repo},
       dataType:"json",
       type:"post",
     }).done(function(data,statusText,jqXHR){
       if(jqXHR.status===200){
          $("#new_auth_tr").remove();
          $("#new_auth").before(`<tr>
              <td id="cpath_`+index+`">`+newpath+`</td>
              <td id='cuser_`+index+`'>`+newuser+`</td>
              <td id="cauth_`+index+`">`+newauth+`</td>
              <td id="config_op_`+index+`">
                  <button id="config_edit_`+index+`" type="button" class="btn btn-default " style="display:inline"><span class="glyphicon glyphicon-pencil"></span></button>
                  <button type="button" id="config_ok_`+index+`" class="btn btn-success " style="display:none"><span class="glyphicon glyphicon-ok"></span></button>
                   &nbsp;<button id="config_del_`+index+`" type="button" class="btn btn-danger "><span class="glyphicon glyphicon-trash"></span></button>
              </td>
            </tr>`);

          $("#config_edit_"+index).on('click',function(){
             var users = $("#cuser_"+index).text()
             var auth = $("#cauth_"+index).text()
              $("#cuser_"+index).remove();
              $("#cauth_"+index).remove();
            //  $("#config_op_"+index).before(`<td id="cuser_`+index+`"><input type="text" class="form-control" id="edited_cuser_`+index+`" value="`+users+`"></td>`)
              $("#config_op_"+index).before(`<td id="cuser_`+index+`"><select multiple="multiple" size="2" class="form-control" id="edited_cuser_`+index+`" for="`+users+`"><option value='*'>*</option></select></td>`)
              $("#config_op_"+index).before(`<td id="cauth_`+index+`"><input type="text" class="form-control" id="edited_cauth_`+index+`" value="`+auth+`"></td>`)

              // load data from users column
              $.ajax({
                url:"/users",
                data:{repo:repo},
              }).done(function(data,textStatus,jqXHR){
                var obj = jQuery.parseJSON(data)
                for (key in obj) {
                  // console.log(key)
                  $("#edited_cuser_"+index).append("<option value='"+key+"'>"+key+"</option>")
                }

              })
              // load data from groups column
              $.ajax({
                url:"/groups",
                data:{repo:repo},
              }).done(function(data,textStatus,jqXHR){
                var obj = jQuery.parseJSON(data)
                for (key in obj){
                  $("#edited_cuser_"+index).append(`<option value="@`+key+`">@`+key+`</option>`)
                }
              })


              $("#config_edit_"+index).css('display',"none")
              $("#config_ok_"+index).css('display',"inline")
          });

          $('#config_ok_'+index).on('click',function(e){
            //handle ajax  edit post
            var o_user= newuser
            var o_auth= newauth
            var userarray = $('#edited_cuser_'+index).val()
            var new_user="";
            userarray.forEach(function(user,index){
              if(index == userarray.length-1){
                new_user+=user
              }else{
                new_user+=user+","
              }
            });
            var new_auth = $('#edited_cauth_'+index).val()
            var repo = $('#reponame').text();
            var tag =  newpath
            $.ajax({
              url:"/editauth" ,
              data:{reponame:repo,old_users:o_user,old_auth:o_auth,new_users:new_user,new_auth:new_auth,tag:tag},
              dataType:"json",
              type:"POST",
            }).done(function(data,statusText,jqXHR){

              if (jqXHR.status===200){
                $("#cuser_"+index).remove();
                $("#cauth_"+index).remove();
                $("#config_op_"+index).before("<td id='cuser_"+index+"'>"+new_user+"</td>");
                $("#config_op_"+index).before("<td id='cauth_"+index+"'>"+new_auth+"</td>");
                $('#config_edit_'+index).css("display","inline");
                $('#config_ok_'+index).css("display","none");
              }
            });
          });

          $('#config_del_'+index).on('click',function(e){
            //handle ajax delete
            var repo = $('#reponame').text();

            $.ajax({
              url:"/delauth",
              data:{reponame:repo,old_users:newuser,old_auth:newauth,tag:newpath},
              type:"POST",
              dataType:"json",
            }).done(function(data,statusText,jqXHR){
               if (jqXHR.status ===200){
                 $("#cuser_"+index).remove();
                 $("#cpath_"+index).remove();
                 $("#cauth_"+index).remove();
                 $("#config_op_"+index).remove();
               }
            });
          });

        $("#new_auth").css("display","block");
       }
     });
  });
});
