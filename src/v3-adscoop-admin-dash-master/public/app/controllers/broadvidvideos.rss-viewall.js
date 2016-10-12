app.controller('BroadvidVideosRssViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidvideos/rss/viewall").query().$promise.then(function(groups) {
    $scope.rss = groups;
  })
})
