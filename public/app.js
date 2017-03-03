function maybejoin(e) {
  if (e.value.endsWith("\n")) {
    e.value = e.value.replace("\n", "");
  }
  join();
}
function join() {
  var username = document.getElementById("usr").value;
  var ws = new WebSocket("ws://" + window.location.host + "/ws?username=" + username);

  var joinr = document.getElementById("join");
  joinr.style.display = "none";

  var chat = document.getElementById("chat");
  chat.style.display = "block";

  var msgs = document.getElementById("messages");

  ws.addEventListener('message', function(e) {
    var msg = JSON.parse(e.data);

    var node = document.createElement("p");
    var k = document.createElement("span");
    k.setAttribute("class", "user");
    k.appendChild(document.createTextNode(msg["From"]));
    node.appendChild(k);
    var v = document.createElement("span");
    v.setAttribute("class", "message");
    v.appendChild(document.createTextNode(msg["Message"]));
    node.appendChild(v);
    msgs.append(node);
  });

  document.getElementById("send").onkeyup = function(e) {
    if (e.keyCode == 13) {
      ws.send(
        JSON.stringify({
            From: username,
            Message: e.target.value
        })
      );
      e.target.value = "";
    }
  };
}
