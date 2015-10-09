$(function () {
  var ws = new WebSocket("ws://" + location.host + "/ws");
  /*
  ws.onopen = function () {
    var status = $("conn-status");
    status.removeClass("disconnected");
    status.addClass("connected");
  }
  
  ws.onclose = function() {
    var status = $("conn-status");
    status.removeClass("connected");
    status.addClass("disconnected");
  };*/

  ws.onmessage = function (event) {
    var msg = $("<p />").text(event.data);
    $(".msg-box").append(msg);
  };
  /*
  ws.onerror = function() {
    var status = $("conn-status");
    status.removeClass("connected");
    status.addClass("disconnected");
  };*/

  $("#msg-form").submit(function () {
    var msgIn = $(".msg-input");

    if (!msgIn.val()) {
      return false;
    }

    ws.send(msgIn.val());
    msgIn.val("");
    return false;
  })

})
