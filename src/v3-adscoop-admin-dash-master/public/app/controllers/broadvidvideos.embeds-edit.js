app.controller('BroadvidVideosEmbedsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.embed = {}

  $scope.saveEmbed = function() {
    $http.post("/broadvidvideos/embeds/save", $scope.embed).
    success(function() {

    })
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidvideos/embeds/view/" + $stateParams.id).
    success(function(data) {
      $scope.embed = data;
    })
  }

})
