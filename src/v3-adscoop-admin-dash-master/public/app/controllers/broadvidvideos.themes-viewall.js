app.controller('BroadvidVideosThemesViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("broadvidvideos/themes/viewall").query().$promise.then(function(groups) {
    $scope.themes = groups;
  })
})
