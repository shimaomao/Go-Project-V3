app.controller('BroadvidVideosRedirectsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidvideos/redirects/viewall").query().$promise.then(function(groups) {
    $scope.redirects = groups;
  })
})
