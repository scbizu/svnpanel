$("#new_passwd_btn").on('click',function(){
  $("#new_passwd").css("display","none")
  $("#new_passwd").before(`<tr id="new_passwd_tr">
        <td><input type="text" class="form-control" id="new_passwd_username" value=""></td>
        <td><input type="text" class="form-control" id="new_passwd_pwd" value=""></td>
        <td>
        <button type="button" id="passwd_new_ok" class="btn btn-success " style="display:inline"><span class="glyphicon glyphicon-ok"></span></button>
        &nbsp;&nbsp;&nbsp;
        <button id="passwd_new_del" type="button" class="btn btn-danger "><span class="glyphicon glyphicon-trash"></span></button>
        </td>
  </tr>`)

  $("#passwd_new_del").on('click',function(){
    console.log("quit edit")
    $("#new_passwd_tr").remove();
    $("#new_passwd").css("display","block");
  });

  $("#passwd_new_ok").on('click',function(){
    console.log("adding user ....")
    var repo = $('#reponame').text();
    var user = $("#new_passwd_username").val();
    var passwd = $("#new_passwd_pwd").val();
    $.ajax({
      url:"/newpasswd",
      data:{reponame:repo,new_username:user,new_pwd:passwd},
      type:"POST",
      dataType:"json",
    }).done(function(data,statusText,jqXHR){
        if (jqXHR.status ===200){
          $("#new_passwd_tr").remove();
          $("#new_passwd").before(`<tr>
              <td id="username_`+user+`">`+user+`</td>
              <td id="password_`+user+`">`+passwd+`</td>
              <td id="op_`+user+`">
              <button id="edit_`+user+`" type="button" class="btn btn-default btnuser" style="display:inline"><span class="glyphicon glyphicon-pencil"></span></button>
              <button type="button" id="ok_`+user+`" class="btn btn-success btnuser" style="display:none"><span class="glyphicon glyphicon-ok"></span></button>
              &nbsp;&nbsp;&nbsp;<button  id="del_`+user+`"  type="button" class="btn btn-danger btnuser"><span class="glyphicon glyphicon-trash"></span></button>
              </td>
            </tr>`);

            $('#edit_'+user).on('click',function(e){
              var username = $('#username_'+user).text()
              var password = $('#password_'+user).text()
              //remove td
                $('#username_'+user).remove();
                $('#password_'+user).remove();
              //add new td
                $('#op_'+user).before("<td id='username_"+user+"'><input type='text' class='form-control' id= 'edited_username_"+user+"' value="+username+" ></td>")
                $('#op_'+user).before("<td id='password_"+user+"'><input type='text' class='form-control' id='edited_password_"+user+"' value="+password+" ></td>")
                //change button
                $('#edit_'+user).css("display","none");
                $('#ok_'+user).css("display","inline");
            });
            // confirm edit ..
            $('#ok_'+user).on('click',function(e){
              //handle ajax  edit post
              var o_username= user
              var o_pwd= passwd
              var username = $('#edited_username_'+user).val()
              var password = $('#edited_password_'+user).val()
              var repo = $('#reponame').text();
              $.ajax({
                url:"/passwd" ,
                data:{reponame:repo,old_username:o_username,old_pwd:o_pwd,new_username:username,new_pwd:password},
                dataType:"json",
                type:"POST",
              }).done(function(data,statusText,jqXHR){

                if (jqXHR.status===200){

                  $("#username_"+user).remove();
                  $("#password_"+user).remove();
                  $("#op_"+user).before("<td id='username_"+user+"'>"+username+"</td>");
                  $("#op_"+user).before("<td id='password_"+user+"'>"+password+"</td>");
                  $('#edit_'+user).css("display","inline");
                  $('#ok_'+user).css("display","none");
                }
              });
              console.log(username)
              console.log(password)
            });

            $('#del_'+user).on('click',function(e){
              //handle ajax delete
              var repo = $('#reponame').text();
              $.ajax({
                url:"/delpasswd",
                data:{reponame:repo,old_username:user,old_pwd:passwd},
                type:"POST",
                dataType:"json",
              }).done(function(data,statusText,jqXHR){
                 if (jqXHR.status ===200){
                   console.log("delete a passwd key-pair..")
                   $("#username_"+user).remove();
                   $("#password_"+user).remove();
                   $("#op_"+user).remove();
                 }
              });
            });
        }
    });
    $("#new_passwd").css("display","block");
  });
});
