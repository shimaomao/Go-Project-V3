app.controller('BroadvidAdsWhitelistCountryGroupsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidads/whitelistcountrygroups/viewall").query().$promise.then(function(groups) {
    $scope.groups = groups;
  })
})
