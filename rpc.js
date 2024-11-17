var ws = new WebSocket("ws://localhost:8080/ws");
ws.onopen = function () {
    console.log("Connected to WebSocket server");
    ws.send("Hello, server!");
};
ws.onmessage = function (event) {
    console.log("Message from server:", event.data);
};
ws.onclose = function () {
    console.log("Disconnected from WebSocket server");
};
function rpc_call(functionName) {
    var args = [];
    for (var _i = 1; _i < arguments.length; _i++) {
        args[_i - 1] = arguments[_i];
    }
    return new Promise(function (resolve, reject) {
        ws.send(JSON.stringify({ "type": "rpc-call", function: functionName, args: args }));
        ws.onmessage = function (event) {
            resolve(JSON.parse(event.data));
        };
    });
}
