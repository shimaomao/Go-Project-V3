app.controller('BroadvidAdsBlacklistUrlGroupsEditCtrl', function($scope, $stateParams, $http, $resource) {
  $scope.group = {}

  $scope.uNew = "";

  $scope.removeUrl = function(idx) {
    var url_to_delete = $scope.group.Urls[idx];

    $scope.group.Urls.splice(idx, 1);
  }

  $scope.addUrl = function() {
    if ($scope.group.Urls === null) {
      $scope.group.Urls = [$scope.uNew];
    } else {
      $scope.group.Urls.push($scope.uNew);
    }
    $scope.uNew = "";
  }

  $scope.saveGroup = function() {
    $http.post("/broadvidads/blacklisturlgroups/save", $scope.group).
    success(function() {
    })
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get("/broadvidads/blacklisturlgroups/view/" + $stateParams.id).
    success(function(data) {
      $scope.group = data;
    })
  }
})
