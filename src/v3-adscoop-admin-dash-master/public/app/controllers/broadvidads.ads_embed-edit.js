app.controller('BroadvidAdsEmbedsEditCtrl', function($scope, $stateParams, $http, $resource, $q) {
  $scope.adembed = {};

  $scope.embedId = $stateParams.id;


  $scope.embeds = {};

  var dep1 = $http.get("/broadvidads/embeds/viewall").
  success(function(data) {
    $scope.embeds = data;
  });

  $scope.UrlsWhitelist = {};

  var dep2 = $http.get("/broadvidads/whitelisturlgroups/viewall").
  success(function(data) {
    $scope.UrlsWhitelist = data;
  });

  $scope.UrlsBlacklist = {};

  var dep3 = $http.get("/broadvidads/blacklisturlgroups/viewall").
  success(function(data) {
    $scope.UrlsBlacklist = data;
  });

  $scope.CountriesWhitelist = {};

  var dep4 = $http.get("/broadvidads/whitelistcountrygroups/viewall").
  success(function(data) {
    $scope.CountriesWhitelist = data;
  });

  $q.all([dep1, dep2, dep3, dep4]).then(function() {

      if (typeof $stateParams.embedId === 'undefined') {
        $scope.adembed.AdID = $stateParams.id;
      } else {
        $http.get("/broadvidads/ad_embeds/view/" + $stateParams.embedId).
        success(function(data) {
          $scope.adembed = data;
        });
      }
  })

  $scope.saveAdEmbed = function() {
    $http.post('/broadvidads/ad_embeds/save', $scope.adembed).success(function() {

    });
  }
})
