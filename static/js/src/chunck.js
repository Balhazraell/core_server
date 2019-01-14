'use strict';

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