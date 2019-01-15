/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./src/main.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./src/chunck.js":
/*!***********************!*\
  !*** ./src/chunck.js ***!
  \***********************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


class Chunck {
    constructor(id, state, points_list){
        this.id = id;
        this.state = state;
        this.draw_poins = points_list;
        
        this.normal_color = '#000000';
        this.negative_color = '#FF0000';
        this.positive_color = '#00FF00';
        this.color = this.normal_color;
    }
}

exports.Chunck = Chunck;

/***/ }),

/***/ "./src/main.js":
/*!*********************!*\
  !*** ./src/main.js ***!
  \*********************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


var chunck = __webpack_require__(/*! ./chunck */ "./src/chunck.js");
var websocket = __webpack_require__(/*! ./websockets */ "./src/websockets.js");
var mouse_manager = __webpack_require__(/*! ./mouse_manager */ "./src/mouse_manager.js");

var canvas;
var ctx;
var MousManager;
var selectChunck;

var MessageIntervalID = null;

// Заглушки.
var grid_coordinats = {
    // [0]: new chunck.Chunck(0, [[0, 0],     [100, 0],   [100, 100], [0, 100]]),
    // [1]:new chunck.Chunck(0, [[100, 0],   [200, 0],   [200, 100], [100, 100]]),
    // [2]:new chunck.Chunck(0, [[200, 0],   [300, 0],   [300, 100], [200, 100]]),

    // [3]:new chunck.Chunck(0, [[0,100],   [100, 100], [100, 200], [0, 200]]),
    // [4]:new chunck.Chunck(0, [[100,100], [200, 100], [200, 200], [200, 100]]),
    // [5]:new chunck.Chunck(0, [[200,100], [300, 100], [300, 200], [200, 200]]),

    // [6]:new chunck.Chunck(0, [[0,200],   [100, 200], [100, 300], [0, 300]]),
    // [7]:new chunck.Chunck(0, [[100,200], [200, 200], [200, 300], [100, 300]]),
    // [8]:new chunck.Chunck(0, [[200,200], [300, 200], [300, 300], [200, 300]]),
}

// Запуск приложения после загрузки html страницы.
if (document.readyState == 'complete'){
    app_start();
} else {
    document.addEventListener('DOMContentLoaded', app_start);
}
//-------------------------------------------------------------//

function app_start(){
    canvas = document.getElementById('myCanvas');

    if (!canvas){
        // Вообще надо что бы система работала через коды ошибок.
        // Коды ошибок с сообщениями должен выдавать сервер.
        console.log('Не смог получить myCanvas!');
        return 0;
    }

    ctx = canvas.getContext('2d');
    websocket.connect();
    MousManager = new mouse_manager.MouseManager(canvas);
    var roomCatalog = document.getElementById("roomCatalog");
    roomCatalog.addEventListener("change", changeRoom);

    game_loop();
}

function game_loop(){
    set_chuncks_color();
    draw_grid();

    requestAnimationFrame(game_loop);
}

function draw_grid(){
    // Очищаем канвас.
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    for (var key in grid_coordinats){
        ctx.strokeStyle = grid_coordinats[key].color;
        let chunck_points = grid_coordinats[key].draw_poins
        ctx.beginPath();
        ctx.moveTo(chunck_points[0][0], chunck_points[0][1]);
        for (let pos_index = 1; pos_index < chunck_points.length; pos_index++){
            ctx.lineTo(chunck_points[pos_index][0], chunck_points[pos_index][1])
        }
        ctx.closePath();
        ctx.stroke();

        if (grid_coordinats[key].state == 1) {
            //Надо нарисовать крестик
            let draw_poins = grid_coordinats[key].draw_poins
            let centerX = draw_poins[0][0] + 50
            let centerY = draw_poins[0][1] + 50

            ctx.beginPath();
            ctx.moveTo(centerX, centerY);
            ctx.lineTo(centerX - 25, centerY - 25)
            ctx.lineTo(centerX, centerY)

            ctx.lineTo(centerX + 25, centerY - 25)
            ctx.lineTo(centerX, centerY)

            ctx.lineTo(centerX + 25, centerY + 25)
            ctx.lineTo(centerX, centerY)

            ctx.lineTo(centerX - 25, centerY + 25)
            ctx.lineTo(centerX, centerY)
            
            ctx.stroke();

        } else if (grid_coordinats[key].state == 2){
            let draw_poins = grid_coordinats[key].draw_poins
            let centerX = draw_poins[0][0] + 50
            let centerY = draw_poins[0][1] + 50

            ctx.beginPath();
            ctx.arc(centerX,centerY, 25, 0, Math.PI*2, true);
            ctx.stroke();
        }
    }
}

function set_grid(new_map) {
    grid_coordinats = [];
    for (var key in new_map){
        let newChunck = new chunck.Chunck(
            new_map[key].id,
            new_map[key].state,
            new_map[key].coordinates
        );

        grid_coordinats[key] = newChunck;
    } 
}

function set_chuncks_color(){
    // TODO: интересно, что правильнее, собрать информацию о том, с каким чанком мы пересечены и после
    // пробегатся - задавая цвет, или все делать в одном цикле???
    for (var key in grid_coordinats){
        let is_collision = check_collision(grid_coordinats[key])
        
        if (is_collision){
            selectChunck = grid_coordinats[key]
            if (grid_coordinats[key].state == 0){
                grid_coordinats[key].color = grid_coordinats[key].positive_color
            } else{
                grid_coordinats[key].color = grid_coordinats[key].negative_color
            }
        } else {
            grid_coordinats[key].color = grid_coordinats[key].normal_color
        }
    }
} 

// Данный метод - некоторого рода защита на клиенте -
// мы не должны отправлять запрос, если клетка занята.
function check_collision(cunck_for_check){
    var result = false;
    var resultList = [];

    for (let i = 0; i < cunck_for_check.draw_poins.length - 1; i++){ 
        let pos1X = cunck_for_check.draw_poins[i][0];
        let pos1Y = cunck_for_check.draw_poins[i][1];

        let pos2X = cunck_for_check.draw_poins[i + 1][0];
        let pos2Y = cunck_for_check.draw_poins[i + 1][1];

        resultList.push(
            (pos2X - pos1X) * (MousManager.Y - canvas.offsetTop - pos1Y) - (pos2Y - pos1Y) * (MousManager.X - canvas.offsetLeft - pos1X)
        )
    }

    if (resultList.length > 0){
        result = resultList.every(function (element){
            return element >= 0;
        })
    }    

    return result
}

function set_chunck_state(){
    if(selectChunck != undefined){
        if (selectChunck.state == 0){
            websocket.set_chunck_state(selectChunck.id)
        } else {
            // Защита на клиенте.
            send_error("Нельзя изменить значение!");
        }
        
    }
}

// Необходимо вывести ошибку.
function send_error(message=""){
    var messageLabel = document.getElementById('messageLabel');

    if (MessageIntervalID != null && message == ""){
        clearInterval(MessageIntervalID);
        MessageIntervalID = null;
        messageLabel.innerText = "";
        return
    }

    if (message != "" && MessageIntervalID != null){
        clearInterval(MessageIntervalID);
        MessageIntervalID = null;
        messageLabel.innerText = "";
    }

    messageLabel.innerText = message;

    MessageIntervalID = setInterval(send_error, 5000);
}

function setRoomCatalog(roomsList){
    var roomCatalog = document.getElementById("roomCatalog");
    // сначала очищу список.
    while (roomCatalog.length > 0){
        roomCatalog.remove(roomCatalog.length-1);
    }

    // А теперь заполняем.
    for ( let i = 0; i < roomsList.length; i++) {
        let option = document.createElement("option");
        option.value = roomsList[i];
        option.text = roomsList[i];
        roomCatalog.add(option);
    }
}

function changeRoom(event){
    websocket.sendChangeRoomID(event.target.value)
}

exports.set_grid = set_grid;
exports.set_chunck_state = set_chunck_state;
exports.send_error = send_error;
exports.setRoomCatalog = setRoomCatalog;

/***/ }),

/***/ "./src/mouse_manager.js":
/*!******************************!*\
  !*** ./src/mouse_manager.js ***!
  \******************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


var main = __webpack_require__(/*! ./main */ "./src/main.js")

class MouseManager {
    // не понимаю почему импорт канваса не дает результата...
    constructor(canvas){
        var X = 0;
        var Y = 0;

        canvas.addEventListener('mousemove', this.mouseMove.bind(this));
        canvas.addEventListener('click', this.mouseClick.bind(this));
    }

    mouseMove(event){
        this.X = event.pageX;
        this.Y = event.pageY;
    }

    mouseClick(event){
        main.set_chunck_state();
    }
}



exports.MouseManager = MouseManager;

/***/ }),

/***/ "./src/websockets.js":
/*!***************************!*\
  !*** ./src/websockets.js ***!
  \***************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


var main = __webpack_require__(/*! ./main */ "./src/main.js");

var ws

// Набор функций получаемых от сервера
var handlers = {
    'set_grid': set_grid,
    'send_error': send_error,
    'set_rooms_catalog': set_rooms_catalog,
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

// Пришла ошибка от сервера.
function send_error(message){
    main.send_error(message);
}

// Отправляем запрос на постановку символа в чанк
function set_chunck_state(chunck_id){
    var data = {
        'chunck_id': chunck_id
    }

    var message = {
        'handler_name': 'setChunckState',
        'data': JSON.stringify(data)
    }

    ws.send(JSON.stringify(message));
}

function sendChangeRoomID(roomID){
    var data = {
        'room_id': parseInt(roomID)
    }

    var message = {
        'handler_name': 'chengeRoomID',
        'data': JSON.stringify(data)
    }

    console.log("Отправляем id = ", message)

    ws.send(JSON.stringify(message));
}

function set_rooms_catalog(roomsIDs){
    console.log("Католог комнат:", roomsIDs)
    main.setRoomCatalog(roomsIDs);
}

exports.connect = connect;
exports.set_chunck_state = set_chunck_state;
exports.sendChangeRoomID = sendChangeRoomID;

/***/ })

/******/ });
//# sourceMappingURL=main_build.js.map