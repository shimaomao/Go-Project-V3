app.controller('BroadvidAdsBlacklistUrlGroupsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidads/blacklisturlgroups/viewall").query().$promise.then(function(groups) {
    $scope.groups = groups;
  })
})
