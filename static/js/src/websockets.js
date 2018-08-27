'use strict';

var main = require('./main');

var ws

// Набор функций получаемых от сервера
var handlers = {
    'set_grid': set_grid
};

// Пошла работа с websockets
function connect(){
    // Здесь происходит падение если не можем подключится, надо красиво обработать...
    ws = new WebSocket('ws://127.0.0.1:8081/appgame');
    ws.onopen = open;
    ws.onclose = close;
    ws.onmessage = message;
}

// websocket стартанул.
function open(event){
    console.log('websocket is open!');
}

// websocket закрылся.
function close(event){
    console.log('websocket is close!');
}

// пришло сообщение по websocket.
function message(event){
    var data = JSON.parse(event.data);
    handlers[data['handler_name']](JSON.parse(data['data']));
}

// ------------- incoming ------------------
// Пришла сетка.
function set_grid(new_map){
    main.set_grid(new_map);
}

// Отправляем запрос на постановку символа в чанк
function set_chunck_state(chunck_id){
    message = {
        'handler_name': 'setChunckState',
        'data': chunck_id.toString()
    }

    ws.send(JSON.stringify(message));
}

exports.connect = connect;
exports.set_chunck_state = set_chunck_state;