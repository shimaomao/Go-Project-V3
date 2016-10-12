	'use strict';

	String.prototype.isEmpty = function() {
	  return (this.length === 0 || !this.trim());
	};

	angular.module('app', [
	  'ngAnimate',
	  'ngCookies',
	  'ngResource',
	  'ngSanitize',
	  'ngTouch',
	  'ngStorage',
	  'ui.router',
	  'ncy-angular-breadcrumb',
	  'ui.bootstrap',
	  'ui.utils',
	  'oc.lazyLoad',
	  'datatables',
	  'ui.bootstrap.datetimepicker'
	]);

	angular.module('app')
	  .directive('inlineEdit', function() {
	    return {
	      restrict: 'A',
	      transclude: true,
	      template: '<label class="editing" data-ng-transclude></label>',
	      controller: ['$scope', '$element', '$transclude', function($scope, $element, $transclude) {
	        $transclude(function(clone) {
	          $scope.transcluded_content = clone[0].textContent;
	        });
	        $element.bind('click', function() {
	          $element.hide().after('<input type="text" value="' + $scope.transcluded_content + '" />');

	          $element.next().focus().blur(function() {
	            $scope.transcluded_content = $element.next().val();
	            $element.html($scope.transcluded_content);
	            $element.next().hide();
	            $element.show();
	          });
	        });

	      }]
	    };
	  })
	  .config(
	    ['$stateProvider', '$urlRouterProvider',
	      function($stateProvider, $urlRouterProvider) {
	        //
	        // For any unmatched url, redirect to /state1
	        $urlRouterProvider.otherwise("/adscoops");
	        //
	        // Now set up the states
	        $stateProvider

	        // Adscoops

	          .state('adscoops', {
	            url: "/adscoops",
	            templateUrl: "partials/adscoops.index.html",
	            resolve: {
	              deps: [
	                '$ocapp',
	                function($ocLazyLoad) {
	                  return $ocLazyLoad.load({
	                    serie: true,
	                    files: [
	                      'lib/jquery/charts/sparkline/jquery.sparkline.js',
	                      'lib/jquery/charts/easypiechart/jquery.easypiechart.js',
	                      'lib/jquery/charts/flot/jquery.flot.js',
	                      'lib/jquery/charts/flot/jquery.flot.resize.js',
	                      'lib/jquery/charts/flot/jquery.flot.pie.js',
	                      'lib/jquery/charts/flot/jquery.flot.tooltip.js',
	                      'lib/jquery/charts/flot/jquery.flot.orderBars.js',
	                      'app/controllers/dashboard.js',
	                      'app/directives/realtimechart.js'
	                    ]
	                  });
	                }
	              ]
	            },
	            controller: function($scope, $http, $resource) {
	              $scope.topStats = {}
	              $scope.statsVerticals = {}
	              $http.get("/adscoops/stats/dailyImpressionCount").
	              success(function(data) {
	                console.log(parseInt(data.Revenue.replace(/,/g, "")))
	                console.log(parseInt(data.Impressions.replace(/,/g, "")))
	                data.eCPM = (1000 * parseInt(data.Revenue.replace(/,/g, "")) / parseInt(data.Impressions.replace(/,/g, ""))).toFixed(2);
	                $scope.topStats = data;
	              })
	              $http.get("/adscoops/stats/getVerticalStats").
	              success(function(data) {
	                $scope.statsVerticals = data;
	              })
	              $resource("/adscoops/clients/viewall").query().$promise.then(function(clients) {
	                $scope.clients = clients;
	              });


	              $scope.clientInfo = function(client) {
	                client.ShowInfo = !client.ShowInfo;

	                if (!client.ShowInfo) {
	                  return;
	                }
	                client.Campaigns = [];

	                $http.get("/adscoops/campaigns/client-viewextradetails/" + client.ID).
	                success(function(data) {
	                  client.Campaigns = data;
	                })
	              }

	              $scope.showInfo = function(campaign) {
	                console.log(campaign.Paused);
	              }
	            }
	          })
	          .state('adscoops.clients-viewall', {
	            url: "/clients/viewall",
	            templateUrl: "partials/adscoops.clients-viewall.html",
	            controller: function($scope, $resource) {
	              $resource("/adscoops/clients/viewall").query().$promise.then(function(clients) {
	                $scope.clients = clients;
	              });
	            }
	          })
	          .state('adscoops.clients-new', {
	            url: "/clients/new",
	            templateUrl: "partials/adscoops.clients-new.html",
	            controller: adscoopClientsEdit,
	          })
	          .state('adscoops.clients-edit', {
	            url: "/clients/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/adscoops.clients-new.html",
	            controller: adscoopClientsEdit,
	          })
	          .state('adscoops.clients-edit.manual-charge', {
	            url: "/manual-charge",
	            templateUrl: "partials/adscoops.clients-manual-charge.html",
	            controller: function($scope, $stateParams, $http) {
	              $scope.clientID = $stateParams.id;
	            }
	          })
	          .state('adscoops.clients-edit.viewcampaigns', {
	            url: '/campaigns',
	            templateUrl: 'partials/adscoops.campaigns-viewall.html',
	            controller: adscoopViewClientCampaigns
	          })
	          .state('adscoops.campaigns-viewall', {
	            url: "/campaigns/viewall",
	            templateUrl: 'partials/adscoops.campaigns-viewall.html',
	            controller: adscoopViewClientCampaigns
	          })
	          .state('adscoops.clients-edit.viewcampaigns.edit', {
	            url: "/edit/{campaignid:[0-9]{1,11}}",
	            templateUrl: "partials/adscoops.campaigns-edit.html",
	            controller: adscoopCampaignsEdit
	          })
	          .state('adscoops.clients-edit.viewcampaigns.edit.viewschedules', {
	            url: "/schedules",
	            templateUrl: "partials/adscoops.campaigns-schedules-viewall.html",
	            controller: adscoopViewCampaignSchedules
	          })
	          .state('adscoops.clients-edit.viewcampaigns.edit.viewschedules.edit', {
	            url: "/edit/{scheduleid:[0-9]{1,11}}",
	            templateUrl: "partials/adscoops.campaigns-schedules-edit.html",
	            controller: adscoopCampaignScheduleEdit
	          })
	          .state('adscoops.redirects-viewall', {
	            url: "/redirects/viewall",
	            templateUrl: 'partials/adscoops.redirects-viewall.html',
	            controller: adscoopViewRedirects
	          })
	          .state('adscoops.redirects-edit', {
	            url: "/redirects/edit/{id:[0-9]{1,11}}",
	            templateUrl: 'partials/adscoops.redirects-edit.html',
	            controller: adscoopRedirectsEdit
	          })
	          .state('adscoops.redirects-new', {
	            url: "/redirects/new",
	            controller: adscoopRedirectsEdit
	          })

	        // Broadvid Ads


	        .state('broadvidads', {
	            url: "/broadvidads",
	            templateUrl: "partials/broadvidads.index.html",
	            controller: function($scope) {}
	          })
	          .state('broadvidads.ads-viewall', {
	            url: "/ads/viewall",
	            templateUrl: "partials/broadvidads.ads-viewall.html",
	            controller: broadvidadsAdsViewall
	          })
	          .state('broadvidads.ads-new', {
	            url: "/ads/new",
	            templateUrl: "partials/broadvidads.ads-edit.html",
	            controller: broadvidadsAdsEdit
	          })
	          .state('broadvidads.ads-edit', {
	            url: "/ads/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidads.ads-edit.html",
	            controller: broadvidadsAdsEdit
	          })
	          .state("broadvidads.embeds-viewall", {
	            url: "/embeds/viewall",
	            templateUrl: "partials/broadvidads.embeds-viewall.html",
	            controller: broadvidadsEmbedsViewall
	          })
	          .state("broadvidads.embeds-new", {
	            url: "/embeds/new",
	            templateUrl: "partials/broadvidads.embeds-edit.html",
	            controller: broadvidadsEmbedsEdit
	          })
	          .state("broadvidads.embeds-edit", {
	            url: "/embeds/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidads.embeds-edit.html",
	            controller: broadvidadsEmbedsEdit
	          })
	          .state('broadvidads.whitelisturlgroups-viewall', {
	            url: "/whitelist-url-groups/viewall",
	            templateUrl: "partials/broadvidads.whitelist-url-groups-viewall.html",
	            controller: broadvidadsWhitelistUrlGroupsViewall
	          })
	          .state('broadvidads.whitelisturlgroups-new', {
	            url: "/whitelist-url-groups/new",
	            templateUrl: "partials/broadvidads.whitelist-url-groups-edit.html",
	            controller: broadvidadsWhitelistUrlGroupsEdit
	          })
	          .state('broadvidads.whitelisturlgroups-edit', {
	            url: "/whitelist-url-groups/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidads.whitelist-url-groups-edit.html",
	            controller: broadvidadsWhitelistUrlGroupsEdit
	          })
	          .state('broadvidads.blacklisturlgroups-viewall', {
	            url: "/blacklist-url-groups/viewall",
	            templateUrl: "partials/broadvidads.blacklist-url-groups-viewall.html",
	            controller: broadvidadsBlacklistUrlGroupsViewall
	          })
	          .state('broadvidads.blacklisturlgroups-new', {
	            url: "/blacklist-url-groups/new",
	            templateUrl: "partials/broadvidads.blacklist-url-groups-edit.html",
	            controller: broadvidadsBlacklistUrlGroupsEdit
	          })
	          .state('broadvidads.blacklisturlgroups-edit', {
	            url: "/blacklist-url-groups/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidads.blacklist-url-groups-edit.html",
	            controller: broadvidadsBlacklistUrlGroupsEdit
	          })
	          .state('broadvidads.whitelistuagroups-viewall', {
	            url: "/whitelist-ua-groups/viewall",
	            templateUrl: "partials/broadvidads.whitelist-ua-groups-viewall.html",
	            controller: broadvidadsWhitelistUaGroupsViewall
	          })
	          .state('broadvidads.whitelistuagroups-new', {
	            url: "/whitelist-ua-groups/new",
	            templateUrl: "partials/broadvidads.whitelist-ua-groups-edit.html",
	            controller: broadvidadsWhitelistUaGroupsEdit
	          })
	          .state('broadvidads.whitelistuagroups-edit', {
	            url: "/whitelist-ua-groups/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidads.whitelist-ua-groups-edit.html",
	            controller: broadvidadsWhitelistUaGroupsEdit
	          })
	          .state('broadvidads.whitelistcountrygroups-viewall', {
	            url: "/whitelist-country-groups/viewall",
	            templateUrl: "partials/broadvidads.whitelist-country-groups-viewall.html",
	            controller: broadvidadsWhitelistCountryGroupsViewall
	          })
	          .state('broadvidads.whitelistcountrygroups-new', {
	            url: "/whitelist-country-groups/new",
	            templateUrl: "partials/broadvidads.whitelist-country-groups-viewall.html",
	            controller: broadvidadsWhitelistCountryGroupsEdit
	          })
	          .state('broadvidads.whitelistcountrygroups-edit', {
	            url: "/whitelist-country-groups/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidads.whitelist-country-groups-edit.html",
	            controller: broadvidadsWhitelistCountryGroupsEdit
	          })
	          .state('broadvidads.flushadcache', {
	            url: "/flush-ad-cache",
	            templateUrl: "partials/broadvidads.flushadcache.html",
	            controller: broadvidadsWhitelistUaGroupsViewall
	          })

	        // Broadvid Videos

	        .state('broadvidvideos', {
	            url: "/broadvidvideos",
	            templateUrl: "partials/broadvidvideos.index.html",
	            controller: broadvidvideosIndex
	          })
	          .state('broadvidvideos.rss-viewall', {
	            url: "/rss/viewall",
	            templateUrl: "partials/broadvidvideos.rss-viewall.html",
	            controller: broadvidvideosRssViewall
	          })
	          .state('broadvidvideos.rss-new', {
	            url: "/rss/new",
	            templateUrl: "partials/broadvidvideos.rss-edit.html",
	            controller: broadvidvideosRssEdit
	          })
	          .state('broadvidvideos.rss-edit', {
	            url: "/rss/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidvideos.rss-edit.html",
	            controller: broadvidvideosRssEdit
	          })
	          .state('broadvidvideos.embeds-viewall', {
	            url: "/embeds/viewall",
	            templateUrl: "partials/broadvidvideos.embeds-viewall.html",
	            controller: broadvidvideosEmbedsViewall
	          })
	          .state('broadvidvideos.embeds-new', {
	            url: "/embeds/new",
	            templateUrl: "partials/broadvidvideos.embeds-edit.html",
	            controller: broadvidvideosEmbedsEdit
	          })
	          .state('broadvidvideos.embeds-edit', {
	            url: "/embeds/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidvideos.embeds-edit.html",
	            controller: broadvidvideosEmbedsEdit
	          })
	          .state('broadvidvideos.redirects-viewall', {
	            url: "/redirects/viewall",
	            templateUrl: "partials/broadvidvideos.redirects-viewall.html",
	            controller: broadvidvideosRedirectsViewall
	          })
	          .state('broadvidvideos.redirects-new', {
	            url: "/redirects/new",
	            templateUrl: "partials/broadvidvideos.redirects-edit.html",
	            controller: broadvidvideosRedirectsEdit
	          })
	          .state('broadvidvideos.redirects-edit', {
	            url: "/redirects/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidvideos.redirects-edit.html",
	            controller: broadvidvideosRedirectsEdit
	          })
	          .state('broadvidvideos.injectjs-viewall', {
	            url: "/injectjs/viewall",
	            templateUrl: "partials/broadvidvideos.injectjs-viewall.html",
	            controller: broadvidvideosInjectjsViewall
	          })
	          .state('broadvidvideos.injectjs-new', {
	            url: "/injectjs/new",
	            templateUrl: "partials/broadvidvideos.injectjs-edit.html",
	            controller: broadvidvideosInjectjsEdit
	          })
	          .state('broadvidvideos.injectjs-edit', {
	            url: "/injectjs/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidvideos.injectjs-edit.html",
	            controller: broadvidvideosInjectjsEdit
	          })
	          .state('broadvidvideos.themes-viewall', {
	            url: "/themes/viewall",
	            templateUrl: "partials/broadvidvideos.themes-viewall.html",
	            controller: broadvidvideosThemesViewall
	          })
	          .state('broadvidvideos.themes-new', {
	            url: "/themes/new",
	            templateUrl: "partials/broadvidvideos.themes-edit.html",
	            controller: broadvidvideosThemesEdit
	          })
	          .state('broadvidvideos.themes-edit', {
	            url: "/themes/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidvideos.themes-edit.html",
	            controller: broadvidvideosThemesEdit
	          })
	          .state('broadvidvideos.domains-viewall', {
	            url: "/domains/viewall",
	            templateUrl: "partials/broadvidvideos.domains-viewall.html",
	            controller: broadvidvideosDomainsViewall
	          })
	          .state('broadvidvideos.domains-new', {
	            url: "/domains/new",
	            templateUrl: "partials/broadvidvideos.domains-edit.html",
	            controller: broadvidvideosDomainsEdit
	          })
	          .state('broadvidvideos.domains-edit', {
	            url: "/domains/edit/{id:[0-9]{1,11}}",
	            templateUrl: "partials/broadvidvideos.domains-edit.html",
	            controller: broadvidvideosDomainsEdit
	          })


	        // Remove when dev done
	        .state('state2', {
	            url: "/state2",
	            templateUrl: "partials/state2.html"
	          })
	          .state('state2.list', {
	            url: "/list",
	            templateUrl: "partials/state2.list.html",
	            controller: function($scope) {
	              $scope.things = ["A", "Set", "Of", "Things"];
	            }
	          })
	      }
	    ]
	  );

	var adscoopViewClientCampaigns = function($scope, $resource, $stateParams) {
	}

	var adscoopViewCampaignSchedules = function($scope, $stateParams, $resource) {

	}

	var adscoopCampaignScheduleEdit = function($scope, $stateParams, $http) {

	}

	var adscoopCampaignsEdit = function($scope, $stateParams, $http) {

	}

	var adscoopViewRedirects = function($scope, $resource) {


	}

	var adscoopRedirectsEdit = function($scope, $stateParams, $http) {

	}

	var broadvidadsAdsViewall = function($resource, $scope) {

	}

	var broadvidadsAdsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidadsEmbedsViewall = function($resource, $scope) {

	}


	var broadvidadsEmbedsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidadsWhitelistUrlGroupsViewall = function($resource, $scope) {

	}

	var broadvidadsWhitelistUrlGroupsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidadsBlacklistUrlGroupsViewall = function($resource, $scope) {

	}

	var broadvidadsBlacklistUrlGroupsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidadsWhitelistUaGroupsViewall = function($resource, $scope) {

	}

	var broadvidadsWhitelistUaGroupsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidadsWhitelistCountryGroupsViewall = function($resource, $scope) {

	}

	var broadvidadsWhitelistCountryGroupsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidvideosIndex = function($http, $scope) {

	}

	var broadvidvideosRssViewall = function($resource, $scope) {

	}

	var broadvidvideosRssEdit = function($http, $scope, $stateParams) {

	}

	var broadvidvideosEmbedsViewall = function($resource, $scope) {

	}

	var broadvidvideosEmbedsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidvideosRedirectsViewall = function($resource, $scope) {

	}

	var broadvidvideosRedirectsEdit = function($http, $scope, $stateParams) {


	}

	var broadvidvideosInjectjsViewall = function($resource, $scope) {

	}

	var broadvidvideosInjectjsEdit = function($http, $scope, $stateParams) {

	}

	var broadvidvideosThemesViewall = function($resource, $scope) {

	}

	var broadvidvideosThemesEdit = function($http, $scope, $stateParams) {


	}

	var broadvidvideosDomainsViewall = function($resource, $scope) {

	}

	var broadvidvideosDomainsEdit = function($http, $scope, $stateParams) {
	  	}
