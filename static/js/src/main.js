'use strict';

var chunck = require('./chunck');
var websocket = require('./websockets');
var mouse_manager = require('./mouse_manager');

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
        option.value = roomsList[i].ID;
        option.text = roomsList[i].Name;
        roomCatalog.add(option);
    }
}

function changeRoom(event){
    websocket.sendChangeRoomID(event.target.value)
}

function setSelectRoom(roomID){
    document.getElementById("roomCatalog").selectedIndex = selectRoomID;
}

exports.set_grid = set_grid;
exports.set_chunck_state = set_chunck_state;
exports.send_error = send_error;
exports.setRoomCatalog = setRoomCatalog;
exports.setSelectRoom = setSelectRoom;