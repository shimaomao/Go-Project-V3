app.controller('AdscoopCampaignSchedulesViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("/adscoops/campaign-schedules/viewall/" + $stateParams.campaignid).query().$promise.then(function(campaignschedules) {
    $scope.campaignschedules = campaignschedules;
  });
});
