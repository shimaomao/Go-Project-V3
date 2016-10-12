app.controller('BroadvidVideosEmbedsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidvideos/embeds/viewall").query().$promise.then(function(groups) {
    $scope.embeds = groups;
  })
})
