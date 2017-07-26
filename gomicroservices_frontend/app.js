
var bodyParser = require('body-parser');
var express = require('express');
var request = require('request');

var app = express();
const port = 80;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended : true }));

app.use('/', express.static('.'));

app.post('/create', function (req, rez) {

	service = 'http://gomicroservices-api-service/create'
	//service = "http://foxley.co:22222/create"
    console.log( JSON.stringify(req.body))

    var options = {
        url: service,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=UTF-8'
        },
        json: req.body
    };

    request(options, function(err, res, body) {
        console.log("running")
        if (res && (res.statusCode === 200 || res.statusCode === 201)) {
            rez.json(body)
        }
    });
});

app.listen(port, function () {
	console.log('Example app listening on port ' + process.env.EXPOSE_PORT + '!');
});
