app.controller('AdscoopCampaignsEditCtrl', function($scope, $http, $stateParams) {
  $scope.campaign = {};

  $scope.uNew = "";

  $scope.removeUrl = function(idx) {
    var url_to_delete = $scope.campaign.urls[idx];

    $scope.campaign.urls.splice(idx, 1);
  }

  $scope.reactivateUrl = function(idx) {

    if ($scope.campaign.urls === null) {
      $scope.campaign.urls = []
    }
    $scope.campaign.urls.unshift($scope.campaign.all_urls[idx]);
    $scope.campaign.all_urls.splice(idx, 1);
  }

  $scope.addUrl = function() {
    urls = $scope.uNew.split('\n');
    $scope.uNew = "";
    for (i = 0; i < urls.length; i++) {
      if (urls[i].isEmpty()) {
        continue;
      }
      if ($scope.campaign.urls === null) {
        $scope.campaign.urls = [{
          Url: urls[i],
          Weight: 1
        }];
      } else {
        $scope.campaign.urls.unshift({
          Url: urls[i],
          Weight: 1
        });
      }
    }
  }

  $scope.campaignGroups = [];

  $http.get("/adscoops/campaign-groups/byClient/" + $stateParams.id).
  success(function(data) {
    $scope.campaignGroups = data;
  })

  if (typeof $stateParams.campaignid !== 'undefined') {
    $http.get('/adscoops/campaigns/view/' + $stateParams.campaignid).
    success(function(data) {
      data.StartDatetime = new Date(data.StartDatetime);
      data.EndDatetime = new Date(data.EndDatetime);
      $scope.campaign = data;
    })
  } else {
    $scope.campaign.ClientID = parseInt($stateParams.id);
    $scope.campaign.StartDatetime = new Date();
    $scope.campaign.EndDatetime = new Date();
    $scope.campaign.urls = null;
  }

  $scope.saveCampaign = function() {
    $scope.campaign.DailyImpsLimit = parseInt($scope.campaign.DailyImpsLimit);
    $scope.campaign.WeightVariance = parseInt($scope.campaign.WeightVariance);
    $scope.campaign.Type = parseInt($scope.campaign.Type);

    $http.post("/adscoops/campaigns/save", $scope.campaign).
    success(function() {

    });
  }
})
