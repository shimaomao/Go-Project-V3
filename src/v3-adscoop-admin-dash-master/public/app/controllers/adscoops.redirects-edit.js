String.prototype.isEmpty = function() {
    return (this.length === 0 || !this.trim());
};
app.controller('AdscoopRedirectsEditCtrl', function($scope, $stateParams, $http, $resource, $q) {
  $scope.redirect = {};
  $scope.qNew = "";

  var dep1 = $http.get('/adscoops/hosts/viewall').
  success(function(data) {
    $scope.hosts = data;
  })
  var dep2 = $http.get('/adscoops/whitelisturlgroups/viewall').
  success(function(data) {
    $scope.urlgroups = data;
  })
  var dep3 = $http.get('/adscoops/whitelistuagroups/viewall').
  success(function(data) {
    $scope.uagroups = data;
  })
  var dep4 = $http.get('/adscoops/advertiser/campaigns/viewall').
  success(function(data) {
    $scope.advertisingCampaigns = data;
  })

  $q.all([dep1, dep2, dep3, dep4]).then(function() {
    if (typeof $stateParams.id !== 'undefined') {
      $http.get("/adscoops/redirects/view/" + $stateParams.id).
      success(function(data) {
        $scope.redirect = data;
      })
    }
  })

  $scope.qNew = "";

  $scope.removeQS = function(idx) {
    $scope.redirect.AllowedQueryStrings.splice(idx, 1);
  }

  $scope.addQS = function() {
    console.log('added')
    qs = $scope.qNew.split('\n');
    $scope.qNew = "";
    for (i = 0; i < qs.length; i++) {
      if (qs[i].isEmpty()) {
        continue;
      }
      console.log($scope.redirect.AllowedQueryStrings);
      console.log(typeof $scope.redirect.AllowedQueryStrings);
      if ($scope.redirect.AllowedQueryStrings === null) {
        console.log('aqs new');
        $scope.redirect.AllowedQueryStrings = [qs[i]];
      } else {
        console.log('aqs existing');
        $scope.redirect.AllowedQueryStrings.push(qs[i]);
      }
    }

    console.log($scope.redirect.AllowedQueryStrings);
  }

  $scope.saveRedirect = function() {
    $scope.redirect.Min = parseInt($scope.redirect.Min);
    $scope.redirect.Max = parseInt($scope.redirect.Max);
    $scope.redirect.Iframe = parseInt($scope.redirect.Iframe);
    $scope.redirect.AdvertisingDailySpend = parseInt($scope.redirect.AdvertisingDailySpend);
    $scope.redirect.RedirType = parseInt($scope.redirect.RedirType);
    $scope.redirect.BapiScoring = parseInt($scope.redirect.BapiScoring);

    $scope.redirect.RedirType = parseInt($scope.redirect.RedirType);
    $http.post("/adscoops/redirects/save", $scope.redirect).
    success(function() {

    })
  }
});
