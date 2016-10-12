app.controller('BroadvidAdsWhitelistUaGroupsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.group = {}

  $scope.uNew = "";

  $scope.removeUrl = function(idx) {
    var url_to_delete = $scope.group.Uas[idx];

    $scope.group.Uas.splice(idx, 1);
  }

  $scope.addUrl = function() {
    if ($scope.group.Uas === null) {
      $scope.group.Uas = [$scope.uNew];
    } else {
      $scope.group.Uas.push($scope.uNew);
    }
    $scope.uNew = "";
  }

  $scope.saveGroup = function() {
    $http.post("/broadvidads/whitelistuagroups/save", $scope.group).
    success(function() {

    })
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidads/whitelistuagroups/view/" + $stateParams.id).
    success(function(data) {
      $scope.group = data;
    })
  }

})
