export default class SocketClient {
    initSocketConnection = (connKey, callback) => {
        if (window["WebSocket"]) {
            let conn = new WebSocket("ws://localhost:3456/api/v1/ws");
            conn.onclose = function (evt) {
                console.log('Socket Connection Closed!');
            };
            conn.onopen = function (evt) {
                console.log('Socket Connection Opened!');
                conn.send(connKey);
            }
            conn.onmessage = function (evt) {
                const data = JSON.parse(JSON.parse(evt.data));
                console.log('Data received from socket', data);
                callback(data);
            };
            return conn;
        }
    };
}