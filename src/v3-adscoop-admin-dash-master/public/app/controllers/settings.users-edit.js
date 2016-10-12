app.controller('SettingsUsersEditCtrl', function ($scope, toaster, $http, $stateParams) {

  $scope.user = {};
  if ($scope.stateParams !== 'undefined') {

  } else {
    $http.get('/settings/users/view/' + $scope.stateParams)
    .success(function(data) {
      $scope.user = data;
    })
  }

  $scope.saveUser = function() {
    $http.post('/settings/users/save', $scope.user).
    success(function() {
      toaster.success('User saved');
    })
  }
})
