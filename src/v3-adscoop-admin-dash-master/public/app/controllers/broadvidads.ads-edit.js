app.controller('BroadvidAdsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.ad = {};

  $http.get("/broadvidads/ads/view/" + $stateParams.id).
  success(function(data) {
    $scope.ad = data;
  });

  $scope.removeAdEmbed = function(ae) {
    $http.post('/broadvidads/ad_embeds/remove', ae).success(function() {
      $scope.updateAdEmbeds()
    })
  }

  $scope.pauseAdEmbed = function(ae) {
    $http.post('/broadvidads/ad_embeds/pause', ae).success(function() {
      $scope.updateAdEmbeds()
    })
  }

  $scope.updateAdEmbeds = function() {
      $http.get('/broadvidads/ads/view/' + $stateParams.id).
      success(function(data) {
        $scope.ad.AdDesktop = data.AdDesktop;
        $scope.ad.AdHTML5 = data.AdHTML5;
        $scope.ad.PlayerDesktop = data.PlayerDesktop;
        $scope.ad.PlayerHTML5 = data.PlayerHTML5;
        $scope.ad.DefaultTag = data.DefaultTag;
      })
  }

  $scope.copyAdEmbed = function(ae) {
    $http.post('/broadvidads/ad_embeds/copy', ae).success(function() {
      $scope.updateAdEmbeds();
    })
  }


})
