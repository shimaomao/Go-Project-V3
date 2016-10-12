app.controller('BroadvidAdsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("/broadvidads/ads/viewall").query().$promise.then(function(ads) {
    $scope.ads = ads;
  })
})
