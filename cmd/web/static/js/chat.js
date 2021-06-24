let conn;
let msg = document.querySelector('input[name="message"]');
let chat = document.querySelector(".chat-scroll");
let clientUsername = document.querySelector(".badge").innerHTML
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
                appendStart(msg.text, msg.username, new Date(msg.created));
            })
        } else {
            update.messages.forEach((msg) => {
                appendEnd(msg.text, msg.username, new Date(msg.created));
            })
        }
    };
} else {
    appendEnd("<b>Your browser does not support WebSockets.</b>", new Date());
}

chat.onscroll = () => {
    if (
        !reachedHistoryEnd &&
        (chat.scrollHeight + chat.scrollTop - chat.clientHeight) <= 500
    ) {
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

function createMessage(text, username, date) {
    let messageItem = document.createElement("div");
    let textTimeWrapper = document.createElement("div")
    textTimeWrapper.setAttribute("class", "d-flex flex-wrap")

    let textItem = document.createElement("div");
    textItem.innerHTML = text;
    textItem.setAttribute("class", "text")

    let timeItem = document.createElement("div");
    timeItem.innerHTML = new Date(date).toLocaleTimeString(undefined, {hour: "2-digit", minute: "2-digit"});

    if (username !== clientUsername) {
        messageItem.setAttribute("class", "px-2 py-1 my-2 rounded-3 bg-primary message");
        timeItem.setAttribute("class", "time ps-2");

        let usernameItem = document.createElement("div");
        usernameItem.innerHTML = username;
        usernameItem.setAttribute("class", "username")
        messageItem.appendChild(usernameItem);
    } else {
        messageItem.setAttribute("class", "px-2 py-1 my-2 rounded-3 bg-secondary client-message");
        timeItem.setAttribute("class", "client-time ps-2");
    }

    textTimeWrapper.appendChild(textItem);
    textTimeWrapper.appendChild(timeItem);
    messageItem.appendChild(textTimeWrapper)

    return messageItem;
}

function appendStart(text, username, date) {
    // insert after add to the top because our chat is column-reverse flexbox container
    chat.append(createMessage(text, username, date));
}

function appendEnd(text, username, date) {
    let message = createMessage(text, username, date);
    if (chat.children.length === 0) {
        chat.append(message);
        return;
    }

    // insert before adds to the bottom because our chat is column-reverse flexbox container
    chat.insertBefore(createMessage(text, username, date), chat.firstChild);
}
