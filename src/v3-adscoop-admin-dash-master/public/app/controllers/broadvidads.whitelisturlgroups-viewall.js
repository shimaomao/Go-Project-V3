app.controller('BroadvidAdsWhitelistUrlGroupsViewallCtrl', function($scope, $stateParams, $http, $resource) {

    $resource("broadvidads/whitelisturlgroups/viewall").query().$promise.then(function(groups) {
      $scope.groups = groups;
    })
})
