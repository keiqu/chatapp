window.onload = function () {
    let conn;
    let msg = document.querySelector('input[name="message"]');
    let chat = document.querySelector(".chat-scroll");

    function appendChat(message) {
        // create new message item
        let item = document.createElement("div")
        item.setAttribute("class", "px-3 py-1 my-1 bg-primary rounded-3 trim")
        item.innerHTML = message

        chat.insertBefore(item, chat.firstChild)

        // scroll to bottom
        chat.scrollTop = chat.scrollHeight - chat.clientHeight;
    }

    document.querySelector("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        conn.send(msg.value);
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (ev) {
            appendChat("<b>Connection closed.</b>");
        };
        conn.onmessage = function (ev) {
            appendChat(ev.data)
        };
    } else {
        appendChat("<b>Your browser does not support WebSockets.</b>");
    }
};
