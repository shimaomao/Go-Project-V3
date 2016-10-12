app.controller('AdscoopCampaignGroupsEdit', function($scope, $stateParams, $http, $resource) {
  $scope.embed = {};


  if (typeof $stateParams.campaigngroupid != 'undefined') {
    $http.get("/adscoops/campaign-groups/view/" + $stateParams.campaigngroupid).
    success(function(campaignGroup) {
      $scope.campaignGroup = campaignGroup;
    })
  }

  $scope.saveCampaignGroup = function() {
    $http.post('/adscoops/campaign-groups/save', $scope.campaignGroup);
  }
})
