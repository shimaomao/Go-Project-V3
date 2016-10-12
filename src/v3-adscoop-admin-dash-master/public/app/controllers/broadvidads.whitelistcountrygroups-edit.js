app.controller('BroadvidAdsWhitelistCountryGroupsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.group = {}

  $scope.uNew = "";

  $scope.removeUrl = function(idx) {
    var url_to_delete = $scope.group.Countries[idx];

    $scope.group.Countries.splice(idx, 1);
  }

  $scope.addUrl = function() {
    if ($scope.group.Countries === null) {
      $scope.group.Countries = [$scope.uNew];
    } else {
      $scope.group.Countries.push($scope.uNew);
    }
    $scope.uNew = "";
  }

  $scope.saveGroup = function() {
    $http.post("/broadvidads/whitelistcountrygroups/save", $scope.group).
    success(function() {

    })
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidads/whitelistcountrygroups/view/" + $stateParams.id).
    success(function(data) {
      $scope.group = data;
    })
  }

})
