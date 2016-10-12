String.prototype.isEmpty = function() {
    return (this.length === 0 || !this.trim());
};
app.controller('AdscoopRedirectsEditCampaignGroupsCtrl', function($scope, $stateParams, $http, $resource, $q) {
  $scope.campaignsGroups = {}
  $scope.asscCampaignGroups = [];

  $scope.clientID = $stateParams.id;

  $http.get("/adscoops/campaign-groups/byRedirect/" + $stateParams.id).
  success(function(data) {
    $scope.asscCampaignGroups = data;
  })

  $http.get("/adscoops/campaign-groups/viewall").
  success(function(data) {
    $scope.campaignGroups = data;
  })

  $scope.addCampaignGroup = function() {
    var newCampaign = {
      CampaignGroupID: $scope.campaignGroups[$scope.addCampaignGroupIndex].ID,
      RedirectID: parseInt($stateParams.id),
      Name: $scope.campaignGroups[$scope.addCampaignGroupIndex].Name
    }
    $scope.asscCampaignGroups.push(newCampaign);
  }

  $scope.removeCampaignGroup = function(index) {
    $scope.asscCampaignGroups.splice(index, 1);
  }


  $scope.saveCampaignGroup = function() {
    $http.post("/adscoops/campaign-groups/saveByRedirect/" + $stateParams.id, $scope.asscCampaignGroups);
  }


});
