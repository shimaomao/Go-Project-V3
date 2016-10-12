app.controller('BroadvidEmbedsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.embed = {};


  if (typeof $stateParams.id != 'undefined') {
    $http.get("/broadvidads/embeds/view/" + $stateParams.id).
    success(function(data) {
      $scope.embed = data;
    })  
  }

  $scope.saveEmbed = function() {
    $http.post('/broadvidads/embeds/save', $scope.embed);
  }
})
