require(["jquery", "underscore"], function($, _) {
  "use strict";
  var messages = [];
  var rescroll = function() {
    var $msgs = $("#chat-scroll");
    $msgs.scrollTop($msgs.children().height());
  };
  var update = function(d) {
    messages.push.apply(messages, d);
    $(function() {
      $("#chat-msgs").text(_.map(messages, function(l) {
        return l["Time"] + ": " + l["Name"] + ": " + l["Message"];
      }).join("\n"));
    });
  };
  (function upoll(u) {
    $.ajax(u).success(function(d) {
      upoll(u);
      update(d);
      rescroll();
    });
  }("/await"));
  $.ajax("/recv").success(function(d) {
    update(d);
    rescroll();
  });
  $(function() {
    $("form#chat-send").on("submit", function(evt) {
      evt.preventDefault();
      var self = this;
      var $self = $(self);
      var $msg = $self.find("[name=message]");
      var $user = $self.find("[name=name]");
      $.ajax({
        url: "/send",
        method: "POST",
        data: {message: $msg.val(), name: $user.val()}
      }).success(function() {
        $msg.val("");
        rescroll();
      });
    });
  });
});
