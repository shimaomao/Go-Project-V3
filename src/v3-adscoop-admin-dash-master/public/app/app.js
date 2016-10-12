'use strict';
angular.module('app', [
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngSanitize',
    'ngTouch',
    'ngStorage',
    'ui.router',
    'ncy-angular-breadcrumb',
    'ui.bootstrap',
    'ui.utils',
    'oc.lazyLoad',
    'ngWebsocket',
    'toaster',
    'ui.select2',
    'chart.js',
    'angular-animation-counter',
]).run(function($websocket, toaster, $rootScope, $http) {
  $rootScope.messages = [];

  // $http.get('/messages/getForUser').success(function() {
  //   $rootScope.message = data;
  // })

  var secureConnection = 'ws'
  if (location.protocol === 'https:') {
    secureConnection = 'wss'
  }
  var ws = $websocket.$new({'url': secureConnection + '://' + location.host + '/rtupdates', 'reconnect': true, 'protocols': [], 'subprotocols': ['base46'] });

  ws.$on('$open', function() {
    ws.$emit('ping', 'hi, listening websocket server');
  })

  function pushToRoot(data) {
    $rootScope.messages.unshift(data);
    $rootScope.messages = $rootScope.messages.slice(0,5);
  }


  ws.$on('message', function(data) {
    pushToRoot(data);
    toaster.pop(data.type, data.title, data.message);
  })
}).factory('adscoopsSocket', function($rootScope, $websocket) {
  var secureConnection = 'ws'
  if (location.protocol === 'https:') {
    secureConnection = 'wss'
  }
  var asws = $websocket.$new({'url': secureConnection + '://' + location.host + '/adscoopsupdates', 'lazy': true, 'reconnect': true, 'protocols': [], 'subprotocols': ['base46']});

  setTimeout(function() {
    asws.$open();
  }, 1000);

  return {
    on: function(eventName, callback) {
      asws.$on(eventName, function() {
        var args = arguments;
        $rootScope.$apply(function () {
          callback.apply(asws, args);
        })
      })
    }
  }
}).config(function(ChartJsProvider) {
  ChartJsProvider.setOptions({
    pointHitDetectionRadius:5,
  })
});
