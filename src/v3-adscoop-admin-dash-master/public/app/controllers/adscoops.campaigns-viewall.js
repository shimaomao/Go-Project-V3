app.controller('AdscoopCampaignsViewall', function($scope, $stateParams, $http, $resource) {
  $scope.clientID = 0;

  if (typeof $stateParams.id !== 'undefined') {
    $resource("/adscoops/campaigns/clientviewall/" + $stateParams.id).query().$promise.then(function(campaigns) {
      $scope.clientID = $stateParams.id;
      $scope.campaigns = campaigns;
    });
  } else {
    $resource("/adscoops/campaigns/viewall").query().$promise.then(function(campaigns) {
      $scope.campaigns = campaigns;
    });
  }
})
