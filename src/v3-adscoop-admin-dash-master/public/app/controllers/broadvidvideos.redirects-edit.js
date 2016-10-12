app.controller('BroadvidVideosRedirectsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.redirect = {}
  $scope.rss = {}

  $http.get("/broadvidvideos/rss/viewall").
  success(function(data) {
    $scope.rss = data;
  })

  $scope.saveRedirect = function() {
    $http.post("/broadvidvideos/redirects/save", $scope.redirect).
    success(function() {

    })
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidvideos/redirects/view/" + $stateParams.id).
    success(function(data) {
      $scope.redirect = data;
    })
  }
})
