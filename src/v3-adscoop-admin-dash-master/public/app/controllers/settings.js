app.controller('AdminSettingsCtrl', function($scope, $resource, $http) {
  $scope.users = {};
  $http.get("/settings/users/viewall")
  .success(function(data) {
    $scope.users = data;
  })
});
