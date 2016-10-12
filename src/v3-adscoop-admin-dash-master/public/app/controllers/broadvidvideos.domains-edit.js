app.controller('BroadvidVideosDomainsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.domain = {}
  $scope.redirects = {}
  $scope.themes = {}

  $http.get("/broadvidvideos/redirects/viewall").
  success(function(data) {
    $scope.redirects = data;
  })

  $http.get("/broadvidvideos/themes/viewall").
  success(function(data) {
    $scope.themes = data;
  })

  $scope.saveDomain = function() {
    $http.post("/broadvidvideos/domains/save", $scope.domain).
    success(function() {

    })
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidvideos/domains/view/" + $stateParams.id).
    success(function(data) {
      $scope.domain = data;
    });
  }
})
