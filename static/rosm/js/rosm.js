$(document).ready(function(){

    $('#container').highcharts({
	chart : {
	    type : 'line',
	    events : {
		load : function() {

		    var connection = new WebSocket('ws://127.0.0.1:8888/ws/lastsec/');
		    var self = this;

		    connection.onmessage = function (event) {
		    	var data = JSON.parse(event.data);         
		    	var series = self.series[0];
			console.log(event.data);
			var point = [(new Date()).getTime(),data];
		    	series.addPoint(point);
		    };

		}
	    }
	},
	title : {
	    text : false
	},
	xAxis : {
	    type : 'datetime',
	    minRange : 60 * 1000
	},
	yAxis : {
	    title : {
		text : false
	    }
	},
	legend : {
	    enabled : true
	},
	plotOptions : {
	    series : {
		threshold : 0,
		marker : {
		    enabled : false
		}
	    }
	},
	series : [{name:"tuiles/sec"}]
    });

});
