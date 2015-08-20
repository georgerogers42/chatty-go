require(["jquery", "underscore"], function($, _) {
  "use strict";
  var messages = [];
  var rescroll = function() {
    var $msgs = $("#msgs-scroll");
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
    $("form#send").on("submit", function(evt) {
      evt.preventDefault();
      var self = this;
      var $self = $(self);
      var $msg = $self.find("[name=msg]");
      var $user = $self.find("[name=user]");
      $.ajax({
        url: "/send",
        method: "POST",
        data: {msg: $msg.val(), user: $user.val()}
      }).success(function() {
        $msg.val("");
        rescroll();
      });
    });
  });
});
