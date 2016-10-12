app.controller('BroadvidVideosInjectJSEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.injectjs = {}
  $scope.redirects = {}

  $http.get("/broadvidvideos/redirects/viewall").
  success(function(data) {
    $scope.redirects = data;
  })

  $scope.saveInjectJS = function() {
    $http.post("/broadvidvideos/injectjs/save", $scope.injectjs).
    success(function() {

    })
  }
  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidvideos/injectjs/view/" + $stateParams.id).
    success(function(data) {
      $scope.injectjs = data;
    })
  }
})
