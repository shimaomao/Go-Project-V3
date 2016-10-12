app.controller('AdscoopsClientsViewAllCtrl', function($scope, $resource) {
  $resource("/adscoops/clients/viewall").query().$promise.then(function(clients) {
    $scope.clientsViewAll = clients;
  });
});
