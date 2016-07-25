/* globals WebSocket */
"use strict";

var jsonrpc = function(url, onopen) {
    var jrpc = {
        ws: new WebSocket(url)
    };

    jrpc.ws.onopen = onopen.bind(jrpc);

    jrpc.onerror = function(result, error) {
        console.log("Notification", result, error);
    };

    jrpc._register = {};
    jrpc.register = function(method, cb) {
        jrpc._register[method] = cb;
    };

    jrpc._requests = {};
    jrpc._id = 0;

    jrpc.call = function(method, params, cb) {
        var request = {
            id: jrpc._id++,
            method: method,
            params: params
        };
        jrpc.ws.send(JSON.stringify(request));
        jrpc._requests[request.id] = cb;
    };

    jrpc.notification = function(method, params) {
        var request = {
            method: method,
            params: params
        };
        jrpc.ws.send(JSON.stringify(request));
    };

    jrpc.ws.onerror = function(evt) {
        console.log("ws error", evt);
    };

    jrpc.ws.onmessage = function(evt) {
        var msg = JSON.parse(evt.data);
        if (msg.result || msg.error) {
            if (msg.id) {
                jrpc._requests[msg.id].call(jrpc, msg.resp , msg.error);
            } else {
                jrpc.onerror.call(jrpc, msg.resp , msg.error);
            }
        } else {
            if (msg.method) {
                if (msg.id) {
                    jrpc._register[msg.method].call(this, msg.params,
                            function(result, error) {
                                var response = {
                                    id: msg.id,
                                    result: result,
                                    error: error
                                };
                                jrpc.ws.send(JSON.stringify(response));
                            });
                } else {
                    jrpc._register[msg.method].call(this, msg.params);
                }
            }
        }
    };

    return jrpc;
};
