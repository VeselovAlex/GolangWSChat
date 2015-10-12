/* global $ */
$(function () {
  var ws = new WebSocket("ws://" + location.host + "/ws");
  
  ws.onclose = function () {
    document.cookie = "login=; max-age=-1"
  }
  
  ws.onmessage = function (event) {
    var wrapper = $("<div />").addClass("message");
    console.log(event.data);
    var msg = JSON.parse(event.data);
    $("<div />").addClass("author").text(msg.Author).appendTo(wrapper);
	  var tzH = (new Date()).getTimezoneOffset() / 60;
	  var tzM = (new Date()).getTimezoneOffset() % 60;
    var serverTime = msg.Timestamp.match(/(\d{2}):(\d{2}):(\d{2})/);
    var serverTZ = msg.Timestamp.match(/([+-]\d{2}):(\d{2})/);
    if (serverTZ) {
        serverTime[1] -= (tzH + parseInt(serverTZ[1]));
        serverTime[2] -= (tzM + parseInt(serverTZ[2]));
    }
    var time = serverTime.slice(1).join(":");
    $("<div />").addClass("timestamp").text(time).appendTo(wrapper);
    $("<div />").addClass("content").text(msg.Content).appendTo(wrapper);
    $(".msg-box").append(wrapper);
    wrapper.get(0).scrollIntoView();
  };
  
  ws.onerror = function(err) {
    console.log(err);
  };
  
  $("#msg-form").submit(function () {
    var msgIn = $(".msg-input");

    if (!msgIn.val().trim()){
      return false;
    }

    ws.send(msgIn.val().trim());
    msgIn.val("");
    return false;
  });
  
  $(window).unload(function(){
    ws.close()
  })
})
