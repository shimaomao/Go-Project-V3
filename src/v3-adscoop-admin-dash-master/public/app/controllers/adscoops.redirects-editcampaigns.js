String.prototype.isEmpty = function() {
    return (this.length === 0 || !this.trim());
};
app.controller('AdscoopRedirectsEditCampaignsCtrl', function($scope, $stateParams, $http, $resource, $q) {
  $scope.campaigns = {}
  $scope.asscCampaigns = {}

  $scope.redirID = $stateParams.id;

  $http.get("/adscoops/campaigns/byRedirect/" + $stateParams.id).
  success(function(data) {
    $scope.asscCampaigns = data;
  })

  $http.get("/adscoops/campaigns/viewall").
  success(function(data) {
    $scope.campaigns = data;
  })

  $scope.addCampaign = function() {
    var newCampaign = {
      CampaignID: $scope.campaigns[$scope.addCampaignIndex].ID.toString(),
      RedirectID: $stateParams.id,
      Name: $scope.campaigns[$scope.addCampaignIndex].Name,
      Weight: "1"
    }
    $scope.asscCampaigns.push(newCampaign);
  }

  $scope.removeCampaign = function(index) {
    $scope.asscCampaigns.splice(index, 1);
  }


  $scope.saveCampaigns = function() {
    $http.post("/adscoops/campaigns/saveByRedirect/" + $stateParams.id, $scope.asscCampaigns);
  }


});
