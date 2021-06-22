window.onload = function () {
    let conn;
    let msg = document.querySelector('input[name="message"]');
    let chat = document.querySelector(".chat-scroll");
    let reachedHistoryEnd = false;

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");

        conn.onopen = loadMore;

        conn.onclose = function (ev) {
            appendEnd("<b>Connection closed.</b>", new Date());
        };

        conn.onmessage = function (ev) {
            let update = JSON.parse(ev.data);

            if (update.request !== undefined && update.request.action === "loadMore") {
                if (update.messages.length === 0) {
                    reachedHistoryEnd = true;
                    return;
                }
                update.messages.forEach((msg) => {
                    appendStart(msg.text, new Date(msg.created));
                })
            } else {
                update.messages.forEach((msg) => {
                    appendEnd(msg.text, new Date(msg.created));
                })
            }
        };
    } else {
        appendEnd("<b>Your browser does not support WebSockets.</b>", new Date());
    }

    chat.onscroll = () => {
        if (!reachedHistoryEnd && chat.scrollHeight + chat.scrollTop - chat.clientHeight === 0) {
            loadMore();
        }
    }

    document.querySelector("#msg-form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        conn.send(JSON.stringify({
            "action": "broadcast",
            "message": msg.value
        }));

        msg.value = "";
        return false;
    };

    function loadMore() {
        if (!conn) {
            return;
        }

        // try to load more messages
        conn.send(JSON.stringify({
            "action": "loadMore",
            "offset": chat.children.length
        }));
    }


    function createMessage(text, date) {
        let messageItem = document.createElement("div");
        messageItem.setAttribute("class", "px-2 py-1 my-2 bg-primary rounded-3 trim");

        let textItem = document.createElement("div");
        textItem.innerHTML = text;
        messageItem.appendChild(textItem);

        let timeItem = document.createElement("div");
        timeItem.setAttribute("class", "time");
        timeItem.innerHTML = new Date(date).toLocaleTimeString(undefined, {hour: "2-digit", minute: "2-digit"});
        messageItem.appendChild(timeItem);

        return messageItem;
    }

    function appendStart(text, date) {
        // insert after add to the top because our chat is column-reverse flexbox container
        chat.append(createMessage(text, date));
    }

    function appendEnd(text, date) {
        let message = createMessage(text, date);
        if (chat.children.length === 0) {
            chat.append(message);
            return;
        }

        // insert before adds to the bottom because our chat is column-reverse flexbox container
        chat.insertBefore(createMessage(text, date), chat.firstChild);
    }
}
