app.controller('BroadvidVideosDomainsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidvideos/domains/viewall").query().$promise.then(function(groups) {
    $scope.domains = groups;
  })
})
