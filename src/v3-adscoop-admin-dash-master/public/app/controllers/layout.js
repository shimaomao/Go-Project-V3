app.controller('NavbarCtrl', function($scope, $http, $rootScope) {
  $http.get("/user/info").
  success(function(data) {
    $scope.user = data;
  });

  $scope.messages = $rootScope.messages;

  $scope.$watch(function() {
    return $rootScope.messages;
  }, function() {
    $scope.messages = $rootScope.messages;
  }, true);

  })
