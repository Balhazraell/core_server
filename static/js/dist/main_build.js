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
/******/ 	return __webpack_require__(__webpack_require__.s = "./js/src/main.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./js/src/chunck.js":
/*!**************************!*\
  !*** ./js/src/chunck.js ***!
  \**************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

class Chunck {
    constructor(points_list){
        this.id = -1;
        this.draw_poins = points_list;
        
        this.normal_color = '#000000';
        this.negative_color = '#FF0000';
        this.positive_color = '#00FF00';
        this.color = this.normal_color;
    }
}

// пока старвый вариант мне нравится больше надо посмотреть ES6
// const _Chunck = Chunck;
// export { _Chunck as Chunck };

exports.Chunck = Chunck;

/***/ }),

/***/ "./js/src/main.js":
/*!************************!*\
  !*** ./js/src/main.js ***!
  \************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var chunck = __webpack_require__(/*! ./chunck */ "./js/src/chunck.js");
var websocket = __webpack_require__(/*! ./websockets */ "./js/src/websockets.js");

var canvas;
var ctx;

// Заглушки.
var grid_coordinats = [
    new chunck.Chunck([[0,0],     [0, 100],   [100, 100], [100, 0]]),
    new chunck.Chunck([[0,100],   [0, 200],   [100, 200], [100, 100]]),
    new chunck.Chunck([[0,200],   [0, 300],   [100, 300], [100, 200]]),

    new chunck.Chunck([[100,0],   [100, 100], [200, 100], [200, 0]]),
    new chunck.Chunck([[100,100], [100, 200], [200, 200], [200, 100]]),
    new chunck.Chunck([[100,200], [100, 300], [200, 300], [200, 200]]),

    new chunck.Chunck([[200,0],   [200, 100], [300, 100], [300, 0]]),
    new chunck.Chunck([[200,100], [200, 200], [300, 200], [300, 100]]),
    new chunck.Chunck([[200,200], [200, 300], [300, 300], [300, 200]]),
]

// Запуск приложения после загрузки html страницы.
if (document.readyState != 'loading'){
    app_start();
} else {
    document.addEventListener('DOMContentLoaded', app_start);
}
//-------------------------------------------------------------//

function app_start(){
    canvas = document.getElementById('myCanvas');
    ctx = canvas.getContext('2d');
    websocket.connect();
    draw_grid();
}

function draw_grid(){
    for (let i = 0; i < grid_coordinats.length; i++){
        ctx.strokeStyle = grid_coordinats.color;
        let chunck_points = grid_coordinats[i].draw_poins
        ctx.beginPath();
        ctx.moveTo(chunck_points[0][0], chunck_points[0][1]);
        for (let pos_index = 1; pos_index < chunck_points.length; pos_index++){
            ctx.lineTo(chunck_points[pos_index][0], chunck_points[pos_index][1])
        }
        ctx.closePath();
        ctx.stroke();
    }

    requestAnimationFrame(draw_grid);
}


/***/ }),

/***/ "./js/src/websockets.js":
/*!******************************!*\
  !*** ./js/src/websockets.js ***!
  \******************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


var chunck = __webpack_require__(/*! ./chunck */ "./js/src/chunck.js");

// Набор функций получаемых от сервера
var handlers = {
    'set_grid': set_grid
};

// Пошла работа с websockets
function connect(){
    var ws = new WebSocket('ws://127.0.0.1:8081/connect');
    ws.onopen(open);
    ws.onclose(close);
    ws.onmessage(message);
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
    handlers[data['handler_name']](data['data']);
}

// ------------- incoming ------------------
// Пришла сетка.
function set_grid(data){
    new_grid = data['grid'];
    // По сути очищаем список.
    grid_coordinats = [];
    for (let i = 0; i < new_grid.length; i++){
        let chunck = new chunck.Chunck(
            // первым аргументом потом будет id.
            new_grid[i]
        );
        grid_coordinats.append(chunck);
    } 
}

exports.connect = connect;

/***/ })

/******/ });
//# sourceMappingURL=main_build.js.map