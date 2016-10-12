app.controller('AdscoopRedirectsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("/adscoops/redirects/viewall").query().$promise.then(function(redirects) {
    $scope.redirects = redirects;
  })

  $scope.search = {
  	HideFromDash: false,
  }
});
