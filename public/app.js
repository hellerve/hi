document.getElementById("usr").onkeyup = function(e) {
  if (e.keyCode == 13) {
    join();
  }
};
function join() {
  var username = document.getElementById("usr").value;
  var proto = location.protocol == "http:" ? "ws" : "wss";
  var ws = new WebSocket(proto + "://" + window.location.host + "/ws?username=" + username);

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
    msgs.appendChild(node);
  });

  var roomInsert = document.getElementById("room");

  document.getElementById("send").onkeyup = function(e) {
    if (e.keyCode == 13) {
      ws.send(
        JSON.stringify({
            From: username,
            Message: e.target.value,
            Room: roomInsert.value
        })
      );
      e.target.value = "";
    }
  };
}
