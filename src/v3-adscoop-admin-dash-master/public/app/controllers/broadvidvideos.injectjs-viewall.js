app.controller('BroadvidVideosInjectJSViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidvideos/injectjs/viewall").query().$promise.then(function(groups) {
    $scope.injectjs = groups;
  })
})
