app.controller('AdscoopCampaignScheduleEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.campaign = {}

  $http.get("/adscoops/campaign-schedules/view/" + $stateParams.scheduleid).
  success(function(data) {
    data.StartDatetime = new Date(data.StartDatetime);
    data.EndDatetime = new Date(data.StartDatetime);
    $scope.campaign = data;
  })
});
