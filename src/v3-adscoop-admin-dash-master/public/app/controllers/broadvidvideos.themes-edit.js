app.controller('BroadvidVideosThemesEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.theme = {}

  $scope.saveTheme = function() {
    $http.post("/broadvidvideos/themes/save", $scope.theme).
    success(function() {

    })
  }


  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidvideos/themes/view/" + $stateParams.id).
    success(function(data) {
      $scope.theme = data;
    });
  }
})
