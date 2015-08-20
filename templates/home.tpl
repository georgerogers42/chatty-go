{{ define "head" }}
<title>Chatty</title>
{{ end }}
{{ define "body" }}
<div id="chat-scroll">
  <pre id="chat-msgs"></pre>
</div>
<form action="/send" method="POST" id="chat-send">
<input name="name">
<input name="message">
<input type="submit" value="Send">
</form>
<script data-main="/static/js/app.js" src="/static/js/jam/require.js"></script>
{{ end }}
