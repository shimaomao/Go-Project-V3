app.controller('BroadvidAdsEmbedsViewallCtrl', function($scope, $stateParams, $http, $resource) {
  $resource("/broadvidads/embeds/viewall").query().$promise.then(function(embeds) {
    $scope.embeds = embeds;
  })
})
