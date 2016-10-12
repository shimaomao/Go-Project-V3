String.prototype.isEmpty = function() {
    return (this.length === 0 || !this.trim());
};

﻿'use strict';
angular.module('app')
    .run(
        [
            '$rootScope', '$state', '$stateParams',
            function($rootScope, $state, $stateParams) {
                $rootScope.$state = $state;
                $rootScope.$stateParams = $stateParams;
            }
        ]
    )
    .config(
        [
            '$stateProvider', '$urlRouterProvider',
            function($stateProvider, $urlRouterProvider) {
              $urlRouterProvider
                  .otherwise('/admin/adscoops');
                $stateProvider
                    .state('admin', {
                        abstract: true,
                        url: '/admin',
                        templateUrl: 'views/layout.html?1',
                        resolve: {
                          deps: [
                            '$ocLazyLoad',
                            function($ocLazyLoad) {
                              return $ocLazyLoad.load({
                                serie: true,
                                files: [
                                  'app/controllers/layout.js',
                                ]
                              });
                            }
                          ]
                        }
                    })
                    .state('admin.settings', {
                      url: '/settings',
                      templateUrl: 'views/settings.index.html',
                      ncyBreadcrumb: {
                        label: 'Settings',
                        description: ''
                      },
                      resolve: {
                        deps: [
                          '$ocLazyLoad',
                          function($ocLazyLoad) {
                            return $ocLazyLoad.load({
                              serie: true,
                              files: [
                                'app/controllers/settings.js',
                              ]
                            });
                          }
                        ]
                      }
                    })
                    .state('admin.settings.users-add', {
                      url: '/user/add',
                      templateUrl: 'views/settings.users-edit.html',
                      resolve: {
                        deps: [
                          '$ocLazyLoad',
                          function($ocLazyLoad) {
                            return $ocLazyLoad.load('toaster').then(
                              function() {
                                  return $ocLazyLoad.load({
                                          serie: true,
                                          files: [
                                              'app/controllers/settings.users-edit.js'
                                          ]
                                      }
                                  );
                              }
                            );
                          }
                        ]
                      }
                    })
                    .state('admin.adscoops',  {
                        url: '/adscoops',
                        templateUrl: 'views/adscoops.index.html',
                        ncyBreadcrumb: {
                            label: 'Dashboard',
                            description: ''
                        },
                        resolve: {
                            deps: [
                                '$ocLazyLoad',
                                function($ocLazyLoad) {
                                    return $ocLazyLoad.load({
                                        serie: true,
                                        cache: true,
                                        files: [
                                            'lib/jquery/charts/sparkline/jquery.sparkline.js',
                                            'lib/jquery/charts/easypiechart/jquery.easypiechart.js',
                                            'lib/jquery/charts/flot/jquery.flot.js',
                                            'lib/jquery/charts/flot/jquery.flot.resize.js',
                                            'lib/jquery/charts/flot/jquery.flot.pie.js',
                                            'lib/jquery/charts/flot/jquery.flot.tooltip.js',
                                            'lib/jquery/charts/flot/jquery.flot.orderBars.js',
                                            'lib/jquery/charts/flot/jquery.flot.selection.js',
                                            'lib/jquery/charts/flot/jquery.flot.crosshair.js',
                                            'lib/jquery/charts/flot/jquery.flot.stack.js',
                                            'lib/jquery/charts/flot/jquery.flot.time.js',
                                            'app/controllers/adscoops.dashboard.js',
                                            // 'app/directives/realtimechart.js'
                                        ]
                                    });
                                }
                            ]
                        }
                    })
            				.state('admin.adscoops.clients-viewall', {
            					url: '/clients/viewall',
            					templateUrl: 'views/adscoops.clients-viewall.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.clients-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
            				})
            				.state('admin.adscoops.clients-new', {
            					url: '/clients/new',
            					templateUrl: 'views/adscoops.clients-new.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.clients-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
            				})
            				.state('admin.adscoops.clients-edit', {
            					url: '/clients/edit/{id:[0-9]{1,11}}',
            					templateUrl: 'views/adscoops.clients-new.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.clients-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
            				})
        	          .state('admin.adscoops.clients-edit.manual-charge', {
        	            url: '/manual-charge',
        	            templateUrl: 'views/adscoops.clients-manual-charge.html',
        	            controller: function($scope, $stateParams, $http) {
                        $scope.chargeType = "Manual";
        	              $scope.clientID = $stateParams.id;

                        $scope.chargeClient = function() {
                          $http.post('/adscoops/clients/manualCharge/' + $stateParams.id, {Charge: parseInt($scope.chargeAmount)}).success(function() {
                            alert("Client charged")
                          })
                        }
        	            }
        	          })
        	          .state('admin.adscoops.clients-edit.auto-charge', {
        	            url: '/auto-charge',
        	            templateUrl: 'views/adscoops.clients-manual-charge.html',
        	            controller: function($scope, $stateParams, $http) {
                        $scope.chargeType = "Auto";
                        $scope.clientID = $stateParams.id;

                        $scope.chargeClient = function() {
                          $http.post('/adscoops/clients/autoCharge/' + $stateParams.id, {Charge: parseInt($scope.chargeAmount)}).success(function() {
                            alert("Client charged")
                          })
                        }
        	            }
        	          })
        	          .state('admin.adscoops.clients-edit.viewcampaigns', {
        	            url: '/campaigns',
        	            templateUrl: 'views/adscoops.campaigns-viewall.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaigns-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
                    .state('admin.adscoops.campaign-groups-viewall', {
                      url:"/campaign-groups/viewall",
                      templateUrl: "/views/adscoops.campaign-groups-viewall.html",
                      resolve: {
                        deps: [
                          "$ocLazyLoad",
                          function($ocLazyLoad) {
                            return $ocLazyLoad.load({
                              serie: true,
                              files: [
                                "app/controllers/adscoops.campaign-groups-viewall.js"
                              ]
                            })
                          }
                        ]
                      }
                    })
        	          .state('admin.adscoops.campaign-groups-edit', {
        	            url: '/campaign-groups/edit/{campaigngroupid:[0-9]{1,11}}',
        	            templateUrl: 'views/adscoops.campaign-groups-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaign-groups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.campaign-groups-new', {
        	            url: '/campaign-groups/new',
        	            templateUrl: 'views/adscoops.campaign-groups-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaign-groups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.campaigns-viewall', {
        	            url: '/campaigns/viewall',
        	            templateUrl: 'views/adscoops.campaigns-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaigns-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
                    .state('admin.adscoops.clients-edit.available-campaign-groups', {
                      url:"/available-campaign-groups",
                      templateUrl: 'views/adscoops.client-available-campaign-groups.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.client-available-campaign-groups.js',
                                      ]
                                  });
                              }
                          ]
                      }
                    })
        	          .state('admin.adscoops.clients-edit.viewcampaigns.edit', {
        	            url: '/edit/{campaignid:[0-9]{1,11}}',
        	            templateUrl: 'views/adscoops.campaigns-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaigns-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.clients-edit.viewcampaigns.new', {
        	            url: '/new',
        	            templateUrl: 'views/adscoops.campaigns-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaigns-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.clients-edit.viewcampaigns.edit.viewschedules', {
        	            url: '/schedules',
        	            templateUrl: 'views/adscoops.campaigns-schedules-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaigns-viewschedules.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.clients-edit.viewcampaigns.edit.viewschedules.edit', {
        	            url: '/edit/{scheduleid:[0-9]{1,11}}',
        	            templateUrl: 'views/adscoops.campaigns-schedules-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.campaigns-editschedules.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.redirects-viewall', {
        	            url: '/redirects/viewall',
        	            templateUrl: 'views/adscoops.redirects-viewall.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.redirects-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.redirects-edit', {
        	            url: '/redirects/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/adscoops.redirects-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.redirects-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.redirects-edit.editcampaigns', {
        	            url: '/edit-campaign-groups',
        	            templateUrl: 'views/adscoops.redirects-editcampaigns.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.redirects-editcampaigns.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.redirects-edit.editcampaigngroups', {
        	            url: '/edit-campaigns',
        	            templateUrl: 'views/adscoops.redirects-editcampaigngroups.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.redirects-editcampaigngroups.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.adscoops.redirects-new', {
        	            url: '/redirects/new',
        	            templateUrl: 'views/adscoops.redirects-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.redirects-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })

        	        // Broadvid Ads


        	        .state('admin.broadvidads', {
        	            url: '/broadvidads',
        	            templateUrl: 'views/broadvidads.index.html',
                      ncyBreadcrumb: {
                        label: 'Broadvid Ads',
                        description: ''
                      },

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.dashboard.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.ads-viewall', {
        	            url: '/ads/viewall',
        	            templateUrl: 'views/broadvidads.ads-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.ads-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.ads-new', {
        	            url: '/ads/new',
        	            templateUrl: 'views/broadvidads.ads-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.ads-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.ads-edit', {
        	            url: '/ads/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidads.ads-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.ads-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
                    .state('admin.broadvidads.ads-edit.embeds-new', {
                      url: "/embed-new/{type}",
        	            templateUrl: 'views/broadvidads.ads_embed-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.ads_embed-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
                    })
                    .state('admin.broadvidads.ads-edit.embeds-edit', {
                      url: "/embed-edit/{embedId:[0-9]{1,11}}",
        	            templateUrl: 'views/broadvidads.ads_embed-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.ads_embed-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
                    })
        	          .state('admin.broadvidads.embeds-viewall', {
        	            url: '/embeds/viewall',
        	            templateUrl: 'views/broadvidads.embeds-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.embeds-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.embeds-new', {
        	            url: '/embeds/new',
        	            templateUrl: 'views/broadvidads.embeds-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.embeds-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.embeds-edit', {
        	            url: '/embeds/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidads.embeds-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.embeds-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelisturlgroups-viewall', {
        	            url: '/whitelist-url-groups/viewall',
        	            templateUrl: 'views/broadvidads.whitelist-url-groups-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelisturlgroups-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelisturlgroups-new', {
        	            url: '/whitelist-url-groups/new',
        	            templateUrl: 'views/broadvidads.whitelist-url-groups-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelisturlgroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelisturlgroups-edit', {
        	            url: '/whitelist-url-groups/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidads.whitelist-url-groups-edit.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelisturlgroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.blacklisturlgroups-viewall', {
        	            url: '/blacklist-url-groups/viewall',
        	            templateUrl: 'views/broadvidads.blacklist-url-groups-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.blacklisturlgroups-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.blacklisturlgroups-new', {
        	            url: '/blacklist-url-groups/new',
        	            templateUrl: 'views/broadvidads.blacklist-url-groups-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.blacklisturlgroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.blacklisturlgroups-edit', {
        	            url: '/blacklist-url-groups/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidads.blacklist-url-groups-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.blacklisturlgroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelistuagroups-viewall', {
        	            url: '/whitelist-ua-groups/viewall',
        	            templateUrl: 'views/broadvidads.whitelist-ua-groups-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelistuagroups-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelistuagroups-new', {
        	            url: '/whitelist-ua-groups/new',
        	            templateUrl: 'views/broadvidads.whitelist-ua-groups-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelistuagroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelistuagroups-edit', {
        	            url: '/whitelist-ua-groups/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidads.whitelist-ua-groups-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelistuagroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelistcountrygroups-viewall', {
        	            url: '/whitelist-country-groups/viewall',
        	            templateUrl: 'views/broadvidads.whitelist-country-groups-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelistcountrygroups-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelistcountrygroups-new', {
        	            url: '/whitelist-country-groups/new',
        	            templateUrl: 'views/broadvidads.whitelist-country-groups-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelistcountrygroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.whitelistcountrygroups-edit', {
        	            url: '/whitelist-country-groups/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidads.whitelist-country-groups-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidads.whitelistcountrygroups-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidads.flushadcache', {
        	            url: '/flush-ad-cache',
        	            templateUrl: 'views/broadvidads.flushadcache.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/adscoops.clients-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })

        	        // Broadvid Videos

        	        .state('admin.broadvidvideos', {
        	            url: '/broadvidvideos',
        	            templateUrl: 'views/broadvidvideos.index.html',
                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.dashboard.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.rss-viewall', {
        	            url: '/rss/viewall',
        	            templateUrl: 'views/broadvidvideos.rss-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.rss-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.rss-new', {
        	            url: '/rss/new',
        	            templateUrl: 'views/broadvidvideos.rss-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.rss-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.rss-edit', {
        	            url: '/rss/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidvideos.rss-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.rss-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.embeds-viewall', {
        	            url: '/embeds/viewall',
        	            templateUrl: 'views/broadvidvideos.embeds-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.embeds-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.embeds-new', {
        	            url: '/embeds/new',
        	            templateUrl: 'views/broadvidvideos.embeds-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.embeds-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.embeds-edit', {
        	            url: '/embeds/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidvideos.embeds-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.embeds-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.redirects-viewall', {
        	            url: '/redirects/viewall',
        	            templateUrl: 'views/broadvidvideos.redirects-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.redirects-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.redirects-new', {
        	            url: '/redirects/new',
        	            templateUrl: 'views/broadvidvideos.redirects-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.redirects-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.redirects-edit', {
        	            url: '/redirects/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidvideos.redirects-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.redirects-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.injectjs-viewall', {
        	            url: '/injectjs/viewall',
        	            templateUrl: 'views/broadvidvideos.injectjs-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.injectjs-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.injectjs-new', {
        	            url: '/injectjs/new',
        	            templateUrl: 'views/broadvidvideos.injectjs-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.injectjs-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.injectjs-edit', {
        	            url: '/injectjs/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidvideos.injectjs-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.injectjs-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.themes-viewall', {
        	            url: '/themes/viewall',
        	            templateUrl: 'views/broadvidvideos.themes-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.themes-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.themes-new', {
        	            url: '/themes/new',
        	            templateUrl: 'views/broadvidvideos.themes-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.thems-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.themes-edit', {
        	            url: '/themes/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidvideos.themes-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.themes-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.domains-viewall', {
        	            url: '/domains/viewall',
        	            templateUrl: 'views/broadvidvideos.domains-viewall.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.domains-viewall.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.domains-new', {
        	            url: '/domains/new',
        	            templateUrl: 'views/broadvidvideos.domains-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.domains-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          })
        	          .state('admin.broadvidvideos.domains-edit', {
        	            url: '/domains/edit/{id:[0-9]{1,11}}',
        	            templateUrl: 'views/broadvidvideos.domains-edit.html',

                      resolve: {
                          deps: [
                              '$ocLazyLoad',
                              function($ocLazyLoad) {
                                  return $ocLazyLoad.load({
                                      serie: true,
                                      files: [
                                          'app/controllers/broadvidvideos.domains-edit.js',
                                      ]
                                  });
                              }
                          ]
                      }
        	          });
            }
        ]
    );
