app.controller('AdscoopsClientsEditCtrl', function($scope, $stateParams, $http) {
  $scope.client = {};

  $scope.removeEmail = function(idx) {
    $scope.client.Emails.splice(idx, 1);
  }

  $scope.removeCampaignEmail = function(idx) {
    $scope.client.CampaignEmails.splice(idx, 1);
  }

  $scope.addEmail = function() {

    console.log("ENEW :"+$scope.eNew)
    console.log($scope.eNew);
    emails = $scope.eNew.split('');
    console.log("EMAIL :"+$scope.emails)
    $scope.eNew = "";
    for (i = 0; i < emails.length; i++) {
      if (emails[i].isEmpty()) {
        continue;
      }
      if ($scope.client.Emails === null) {
        $scope.client.Emails = [emails[i]];
      } else {
        $scope.client.Emails.push(emails[i]);
      }
    }
  }

  $scope.addCampaignEmail = function() {
    emails = $scope.eClientNew.split('\n');
    $scope.eClientNew = "";
    for (i = 0; i < emails.length; i++) {
      if (emails[i].isEmpty()) {
        continue;
      }
      if ($scope.client.CampaignEmails === null) {
        $scope.client.CampaignEmails = [emails[i]];
      } else {
        $scope.client.CampaignEmails.push(emails[i]);
      }
    }
  }

  if (typeof $stateParams.id !== 'undefined') {
    $http.get('/adscoops/clients/view/' + $stateParams.id).
    success(function(data) {
      $scope.client = data;
    })
  } else {
    $scope.client.ID = 0;
    $scope.client.Name = "New Client";
  }

  $scope.saveClient = function() {
    $http.post('/adscoops/clients/save', $scope.client).success(function() {

    })
  }
})
