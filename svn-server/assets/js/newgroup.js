$("#new_group_btn").on('click',function(){
  $("#new_group").css("display","none")
  $("#new_group").before(`<tr id="new_group_tr">
        <td><input type="text" class="form-control" id="new_group_gname" value=""></td>
        <td><input type="text" class="form-control" id="new_group_gmember" value=""></td>
        <td>
        <button type="button" id="group_new_ok" class="btn btn-success " style="display:inline"><span class="glyphicon glyphicon-ok"></span></button>
        &nbsp;&nbsp;&nbsp;
        <button id="group_new_del" type="button" class="btn btn-danger "><span class="glyphicon glyphicon-trash"></span></button>
        </td>
  </tr>`)

  $("#group_new_del").on('click',function(){
    console.log("quit edit")
    $("#new_group_tr").remove();
    $("#new_group").css("display","block");
  });

  $("#group_new_ok").on('click',function(){
    console.log("adding user ....")
    var repo = $('#reponame').text();
    var gname = $("#new_group_gname").val();
    var gmember = $("#new_group_gmember").val();
    $.ajax({
      url:"/addgroup",
      data:{reponame:repo,new_groupname:gname,new_users:gmember},
      type:"POST",
      dataType:"json",
    }).done(function(data,statusText,jqXHR){
        if (jqXHR.status ===200){
          $("#new_group_tr").remove();
          $("#new_group").before(`<tr>
              <td id="gname_`+gname+`">`+gname+`</td>

              <td id="gmember_`+gname+`">`+gmember+`</td>
              <td id="group_op_`+gname+`">
              <button id="group_edit_`+gname+`" type="button" class="btn btn-default btnuser" style="display:inline"><span class="glyphicon glyphicon-pencil"></span></button>
              <button type="button" id="group_ok_`+gname+`" class="btn btn-success btnuser" style="display:none"><span class="glyphicon glyphicon-ok"></span></button>
              &nbsp;&nbsp;&nbsp;<button  id="group_del_`+gname+`"  type="button" class="btn btn-danger btnuser"><span class="glyphicon glyphicon-trash"></span></button>
              </td>
            </tr>`);

            $('#group_edit_'+gname).on('click',function(e){
              console.log('edited')
              var groupname = $('#gname_'+gname).text()
              var groupmember = $('#gmember_'+gname).text()
              //remove td
                $('#gname_'+gname).remove();
                $('#gmember_'+gname).remove();
              //add new td
                $('#group_op_'+gname).before("<td id='gname_"+gname+"'><input type='text' class='form-control' id= 'edited_gname_"+gname+"' value="+groupname+" ></td>")
                $('#group_op_'+gname).before("<td id='gmember_"+gname+"'><input type='text' class='form-control' id='edited_gmember_"+gname+"' value="+groupmember+" ></td>")
                //change button
                $('#group_edit_'+gname).css("display","none");
                $('#group_ok_'+gname).css("display","inline");
            });
            // confirm edit ..
            $('#group_ok_'+gname).on('click',function(e){
              //handle ajax  edit post
              var o_gname= gname
              var o_gmember= gmember
              var new_gname = $('#edited_gname_'+gname).val()
              var new_gmember = $('#edited_gmember_'+gname).val()
              var repo = $('#reponame').text();
              $.ajax({
                url:"/groups" ,
                data:{reponame:repo,old_groupname:o_gname,old_users:o_gmember,new_groupname:new_gname,new_users:new_gmember},
                dataType:"json",
                type:"POST",
              }).done(function(data,statusText,jqXHR){

                if (jqXHR.status===200){

                  $("#gname_"+gname).remove();
                  $("#gmember_"+gname).remove();
                  $("#group_op_"+gname).before("<td id='gname_"+gname+"'>"+new_gname+"</td>");
                  $("#group_op_"+gname).before("<td id='gmember_"+gname+"'>"+new_gmember+"</td>");
                  $('#group_edit_'+gname).css("display","inline");
                  $('#group_ok_'+gname).css("display","none");
                }
              });
            });

            $('#group_del_'+gname).on('click',function(e){
              //handle ajax delete
              var repo = $('#reponame').text();
              $.ajax({
                url:"/delgroup",
                data:{reponame:repo,old_groupname:gname,old_users:gmember},
                type:"POST",
                dataType:"json",
              }).done(function(data,statusText,jqXHR){
                 if (jqXHR.status ===200){
                   console.log("delete a passwd key-pair..")
                   $("#gname_"+gname).remove();
                   $("#gmember_"+gname).remove();
                   $("#group_op_"+gname).remove();
                 }
              });
            });
        }
    });
    $("#new_group").css("display","block");
  });
});
