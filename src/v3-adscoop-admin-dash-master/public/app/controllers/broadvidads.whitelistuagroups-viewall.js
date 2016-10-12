app.controller('BroadvidAdsWhitelistUaGroupsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidads/whitelistuagroups/viewall").query().$promise.then(function(groups) {
    $scope.groups = groups;
  })
})
