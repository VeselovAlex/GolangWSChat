$(function () {
  var ws = new WebSocket("ws://" + location.host + "/ws");
  
  ws.onclose = function () {
    document.cookie = "login=; expires=-1"
  }
  
  ws.onmessage = function (event) {
    var wrapper = $("<div />").addClass("message");
    console.log(event.data);
    var msg = JSON.parse(event.data);
    console.log(msg);
    $("<span />").addClass("author").text(msg.Author).appendTo(wrapper);
    $("<span />").addClass("content").text(msg.Content).appendTo(wrapper);
    $(".msg-box").append(wrapper);
  };
  
  ws.onerror = function(err) {
    console.log(err);
  };
  
  $("#msg-form").submit(function () {
    var msgIn = $(".msg-input");

    if (!msgIn.val()) {
      return false;
    }

    ws.send(msgIn.val());
    msgIn.val("");
    return false;
  });
})
