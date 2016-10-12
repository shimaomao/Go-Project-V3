app.controller('AdscoopCampaignGroupsViewall', function($scope, $resource) {
  $resource("/adscoops/campaign-groups/viewall").query().$promise.then(function(campaignGroups) {
    $scope.campaignGroups = campaignGroups;
  });
});
