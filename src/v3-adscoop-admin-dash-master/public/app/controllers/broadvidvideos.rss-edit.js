app.controller('BroadvidVideosRssEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.rss = {}
  $scope.embeds = {};

  $scope.saveRSS = function() {
    $http.post("/broadvidvideos/rss/save", $scope.rss).
    success(function() {

    })
  }

  $http.get("/broadvidvideos/embeds/viewall").
  success(function(data) {
    $scope.embeds = data;
  })

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidvideos/rss/view/" + $stateParams.id).
    success(function(data) {
      $scope.rss = data;
    })
  }
})
