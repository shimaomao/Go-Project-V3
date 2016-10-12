'use strict';

app
// Dashboard Box controller
  .controller('AdscoopsDashboardCtrl', [
  '$rootScope', '$scope', '$http', '$resource', '$interval', '$filter', 'adscoopsSocket', '$timeout',
  function($rootScope, $scope, $http, $resource, $interval, $filter, adscoopsSocket, $timeout) {

  var plotClients, plotRevenue, plotRedir, updatedLength;

    $scope.topStats = {};
    $scope.topStats.Impressions = 0;

      adscoopsSocket.on('dailyImpressionsStats', function(payload) {
        var data = JSON.parse(payload.message);
          $scope.impressionsBilledPercent = parseInt((data.Impressions / data.Limit) * 100);
          $("#impressionsPieChart").easyPieChart({
            percent: $scope.impressionsBilledPercent,
            lineWidth: 3,
            barColor: '#fff',
            trackColor: 'rgba(255,255,255,0.1)',
            scaleColor: false,
            size: 47,
            lineCap: 'butt',
            animate: 500
          })

          $scope.impressionsRequiredPercent = parseInt(((data.Limit - data.Impressions) / data.Limit) * 100);
          $("#impressionsRequiredPieChart").easyPieChart({
            percent: $scope.impressionsRequiredPercent,
            lineWidth: 3,
            barColor: '#fff',
            trackColor: 'rgba(255,255,255,0.1)',
            scaleColor: false,
            size: 47,
            lineCap: 'butt',
            animate: 500
          })
          data.eCPM = (1000 * parseFloat(data.Revenue) / parseInt(data.Impressions)).toFixed(2);
          $scope.topStats = data;
      })

      adscoopsSocket.on('getVerticalStats', function(payload) {
          var data = JSON.parse(payload.message);
          $scope.statsVerticals = data;

          $scope.tcd = [];
          angular.forEach(data.Clicks.Breakdown, function(d) {
            $scope.tcd.push(parseInt(d.Count));
          })
          $scope.statsVerticals = data;

          $scope.tld = [];
          angular.forEach(data.Load.Breakdown, function(d) {
            $scope.tld.push(parseInt(d.Count));
          })
          $scope.statsVerticals = data;

          $scope.tvd = [];
          angular.forEach(data.Verification.Breakdown, function(d) {
            $scope.tvd.push(parseInt(d.Count));
          })

          $("#clicksDataGraph").sparkline($scope.tcd, {
            type: 'bar',
            height: 82,
            width: '100%',
            barColor: '#b0dc81',
            barWidth: 10,
            barSpacing: 1
          });
          $("#loadsDataGraph").sparkline($scope.tld, {
            type: 'bar',
            height: 82,
            width: '100%',
            barColor: '#66CBE4',
            barWidth: 10,
            barSpacing: 1
          });
          $("#verificationsDataGraph").sparkline($scope.tvd, {
            type: 'bar',
            height: 82,
            width: '100%',
            barColor: '#FD9F8D',
            barWidth: 10,
            barSpacing: 1
          });
      })

  $scope.series = ['Yesterday', 'Today'];


  $(window).focus(function() {
        console.log("has focus")
        $scope.$apply();
  });

  $scope.chartoptions = {animationSteps : 5000, scaleBeginAtZero: true, scaleShowLine : false};

      adscoopsSocket.on('getRealtimeRedirStats', function(payload) {
        var retDataYesterday = JSON.parse(payload.message);
        var retDataToday = JSON.parse(payload.title);

        var tempArr = [];
        angular.forEach(retDataYesterday.Data, function(x) {
          tempArr.push(x.AllCount);
        });


        $scope.realTimedata = tempArr;

        updatedLength = $scope.realTimedata.length;

        var tempArr = [];
        var minutesAgo = 30;
        $scope.realtimeRedirLabels = [];
        angular.forEach(retDataYesterday.Data, function(x) {
          $scope.realtimeRedirLabels.push(minutesAgo--);
          tempArr.push(x.Revenue);
        })

        $scope.realTimeRevenue = tempArr;

        var retData = JSON.parse(payload.message);

        var tempArr = [];
        angular.forEach(retDataToday.Data, function(x) {
          tempArr.push(x.AllCount);
        });

        $scope.realTimedata2 = tempArr;

        $scope.realTimeRedir = [
          $scope.realTimedata,
          $scope.realTimedata2
        ];

        var tempArr = [];
        angular.forEach(retDataToday.Data, function(x) {
          tempArr.push(x.Revenue);
        })

        $scope.realTimeRevenue2 = tempArr;


        $scope.realTimeRev = [
          $scope.realTimeRevenue,
          $scope.realTimeRevenue2
        ];

      });

      $scope.options = {
        animation: true,
        showScale: true,
        datasetStrokeWidth: 0.1,
        pointHitDetectionRadios: 1,
        scaleBeginAtZero : true,
        colors : ['#EFEFEF', '#3EC7E9']
      };

      adscoopsSocket.on('getRealTimeClientStats', function(payload) {
        var data = JSON.parse(payload.message);

        if (typeof plotClients === 'undefined') {
          plotClients = $.plot($("#clientRevenuePlot"), {}, clientRealtimeOpts)
        }

        plotClients.setData(data.Clients);
        plotClients.draw();

      });


    $scope.messages = $rootScope.messages;

    $scope.$watch(function() {
      return $rootScope.messages;
    }, function() {
      $scope.messages = $rootScope.messages;
    }, true);

    var today = new Date();
    $scope.dayOfMonth = today.getDate();
    $scope.max = new Date(today.getYear(), today.getMonth(), 0).getDate();
    $scope.daysLeft = new Date(today.getYear(), today.getMonth(), 0).getDate() - today.getDate();

    $scope.statsVerticals = {}

    $scope.showReport = function(c) {

      var startDate = c.ReportStartDate.getFullYear() + "-" + ("0" + (c.ReportStartDate.getMonth() + 1)).slice(-2) + "-" + ("0" + c.ReportStartDate.getDate()).slice(-2);
      var endDate = c.ReportEndDate.getFullYear() + "-" + ("0" + (c.ReportEndDate.getMonth() + 1)).slice(-2) + "-" + ("0" + c.ReportEndDate.getDate()).slice(-2);

      $http.post('/adscoops/clients/showReport/' + c.ID, {ClientID: c.ID, StartDate: startDate, EndDate: endDate}).
      success(function(data) {
        c.reportOutput = data;
        c.showReport = true;
      })
    }

    $scope.setReportToToday = function(c) {
      var today = new Date();
      c.ReportStartDate =  new Date(today.getFullYear(), today.getMonth(), today.getDate());
      c.ReportEndDate =  new Date(today.getFullYear(), today.getMonth(), today.getDate());

      $scope.showReport(c);
    }

    $scope.setReportToYesterday = function(c) {
      var today = new Date();
      c.ReportStartDate =  new Date(new Date(today.getFullYear(), today.getMonth(), today.getDate()).setDate(today.getDate()-1));
      c.ReportEndDate =  new Date(new Date(today.getFullYear(), today.getMonth(), today.getDate()).setDate(today.getDate()-1));

      $scope.showReport(c);
    }

    $scope.setReportToMTD = function(c) {
      var today = new Date();
      c.ReportStartDate =  new Date(today.getFullYear(), today.getMonth(), 1);
      c.ReportEndDate =  new Date(today.getFullYear(), today.getMonth(), today.getDate());

      $scope.showReport(c);
    }

    $scope.downloadReport = function(c) {

            var startDate = c.ReportStartDate.getFullYear() + "-" + ("0" + (c.ReportStartDate.getMonth() + 1)).slice(-2) + "-" + ("0" + c.ReportStartDate.getDate()).slice(-2);
            var endDate = c.ReportEndDate.getFullYear() + "-" + ("0" + (c.ReportEndDate.getMonth() + 1)).slice(-2) + "-" + ("0" + c.ReportEndDate.getDate()).slice(-2);

            var form = document.createElement("form");
            form.setAttribute("method", "post");
            form.setAttribute("action", "/adscoops/clients/showReport/" + c.ID + "?download=true");

            form.setAttribute("target", "_blank");

            var hiddenField = document.createElement("input");
            hiddenField.setAttribute("type", "hidden");
            hiddenField.setAttribute("name", "client_id");
            hiddenField.setAttribute("value", c.ID);
            form.appendChild(hiddenField);

            var hiddenField1 = document.createElement("input");
            hiddenField1.setAttribute("type", "hidden");
            hiddenField1.setAttribute("name", "start_date");
            hiddenField1.setAttribute("value", startDate);
            form.appendChild(hiddenField1);

            var hiddenField2 = document.createElement("input");
            hiddenField2.setAttribute("type", "hidden");
            hiddenField2.setAttribute("name", "end_date");
            hiddenField2.setAttribute("value", endDate);
            form.appendChild(hiddenField2);

            document.body.appendChild(form);

            form.submit();
    }

    var updateDailyImpsCounts = function() {
      $resource("/adscoops/clients/viewVisible").query().$promise.then(function(clients) {
        $scope.clients = clients;


        angular.forEach($scope.clients, function(d, c) {
          var iCount = 0;
          var eCount = 0;
          var lCount = 0;

          var icd = [];
          angular.forEach(d.ImpressionStats, function(stats) {
            var nCount = parseInt(stats.Count)
            iCount = iCount + nCount;
            icd.push(nCount);
          })
          var ecd = [];
          angular.forEach(d.EngagementStats, function(stats) {
            var nCount = parseInt(stats.Count)
            eCount = eCount + nCount;
            ecd.push(nCount);
          })
          var lcd = [];
          angular.forEach(d.LoadStats, function(stats) {
            var nCount = parseInt(stats.Count)
            lCount = lCount + nCount;
            lcd.push(nCount);
          })


          setTimeout(function() {
            if (iCount !== 0) {
              $("#impressionsDataByClient" + d.ID).sparkline(icd, {
                type: 'bar',
                height: 24,
                width: 95,
                barColor: '#b0dc81',
                barWidth: 3,
                barSpacing: 1
              });
            } else {
              $('#impressionsDataByClientCont' + d.ID).remove()
            }

            if (eCount !== 0) {
              $("#engagementsDataByClient" + d.ID).sparkline(ecd, {
                type: 'bar',
                height: 24,
                width: 95,
                barColor: '#FD9F8D',
                barWidth: 3,
                barSpacing: 1
              });
            } else {
              $('#engagementsDataByClientCont' + d.ID).remove()
            }

            $scope.loadRedirSideview = function(points, evt) {
              console.log("points", points);
              console.log("evt", evt);
            }

            if (lCount !== 0) {
              $("#loadsDataByClient" + d.ID).sparkline(lcd, {
                type: 'bar',
                height: 24,
                width: 95,
                barColor: '#66CBE4',
                barWidth: 3,
                barSpacing: 1
              });
            } else {
              $('#loadsDataByClientCont' + d.ID).remove()
            }
          }, 250);

          if (d.UserSettings.ShowInfo) {

            d.Campaigns = [];

            $http.get("/adscoops/campaigns/client-viewextradetails/" + d.ID).
            success(function(data) {
              d.Campaigns = data
              d.editCampaigns = [];
              angular.forEach(jQuery.extend(true,{}, data), function(cmp) {
                d.editCampaigns[cmp.ID] = cmp
              })
            })

            $http.get("/adscoops/clients/viewRedirStats/" + d.ID).
            success(function(data) {
              $scope.clientRedirData["client_redir_data_" + d.ID] = [];
              $scope.clientRedirData["client_redir_labels_" + d.ID] = [];
              angular.forEach(data, function(x, y) {
                $scope.clientRedirData["client_redir_data_" + d.ID].push(x)
                $scope.clientRedirData["client_redir_labels_" + d.ID].push(y)
              })
            });

            $http.get("/adscoops/clients/viewAssociatedRedirects/" + d.ID).
            success(function(data) {
              $scope.clientRedirData["client_redir_list_" + d.ID] = data;
            })

            $scope["client_redirs_" + d.ID] = $interval(function() {
              $http.get("/adscoops/clients/viewRedirStats/" + d.ID).
              success(function(data) {
                $scope.clientRedirData["client_redir_data_" + d.ID] = [];
                $scope.clientRedirData["client_redir_labels_" + d.ID] = [];
                angular.forEach(data, function(x, y) {
                  $scope.clientRedirData["client_redir_data_" + d.ID].push(x)
                  $scope.clientRedirData["client_redir_labels_" + d.ID].push(y)
                })
              })
            }, 5000);
          }
        })
      });
    }

    $scope.clientRedirData = [];

    updateDailyImpsCounts();


    $scope.clientInfo = function(client) {
      client.UserSettings.ShowInfo = !client.UserSettings.ShowInfo;

      $scope.saveUserOptions(client);

      if (!client.UserSettings.ShowInfo) {
        $interval.cancel($scope["client_redirs_" + client.ID]);
        return;
      }
      client.Campaigns = [];

      $http.get("/adscoops/campaigns/client-viewextradetails/" + client.ID).
      success(function(data) {
        client.Campaigns = data
        client.editCampaigns = [];
        angular.forEach(jQuery.extend(true,{}, data), function(cmp) {
          client.editCampaigns[cmp.ID] = cmp
        })
      })


      $http.get("/adscoops/clients/viewRedirStats/" + client.ID).
      success(function(data) {
        $scope.clientRedirData["client_redir_data_" + client.ID] = [];
        $scope.clientRedirData["client_redir_labels_" + client.ID] = [];
        angular.forEach(data, function(x, y) {
          $scope.clientRedirData["client_redir_data_" + client.ID].push(x)
          $scope.clientRedirData["client_redir_labels_" + client.ID].push(y)
        })
      });


      $http.get("/adscoops/clients/viewAssociatedRedirects/" + client.ID).
      success(function(data) {
        $scope.clientRedirData["client_redir_list_" + client.ID] = data;
      })

      $scope["client_redirs_" + client.ID] = $interval(function() {
        $http.get("/adscoops/clients/viewRedirStats/" + client.ID).
        success(function(data) {
          $scope.clientRedirData["client_redir_data_" + client.ID] = [];
          $scope.clientRedirData["client_redir_labels_" + client.ID] = [];
          angular.forEach(data, function(x, y) {
            $scope.clientRedirData["client_redir_data_" + client.ID].push(x)
            $scope.clientRedirData["client_redir_labels_" + client.ID].push(y)
          })
        })
      }, 5000);

    }

    $scope.saveUserOptions = function(client) {
      $http.post('/adscoops/clients/updateCampaignSort', {
        ClientID: client.ID,
        CampaignSort: client.UserSettings.CampaignSort,
        ShowInfo: client.UserSettings.ShowInfo,
        ClientOrder: client.UserSettings.ClientOrder
      }).
      success(function(data) {
      })
    }

    $scope.showInfo = function(client, campaign) {
      angular.forEach(client.Campaigns, function(cmp, key) {
        if (cmp.ID == campaign.ID) {
          console.log("updating campaign")
            if (cmp.Paused != campaign.Paused || cmp.DailyImpsLimit != campaign.DailyImpsLimit || cmp.Cpc != campaign.Cpc) {
              campaign.DailyImpsLimit = parseInt(campaign.DailyImpsLimit);
              $http.post('/adscoops/campaigns/basicSave', campaign).success(function() {
                client.Campaigns[key] = jQuery.extend(true,{}, campaign)
              });
            }
        }
      })
    }

    $scope.boxWidth = $('.box-tabbs').width() - 20;

    var realtimeRevenueOptions = {
        yaxis: {
            color: '#f3f3f3',
            min: 0,
            tickDecimals: 0,
        },
        xaxis: {
            color: '#f3f3f3',
            min: 0,
            tickFormatter: function(val, axis) {
                return "";
            }
        },
        grid: {
            hoverable: true,
            clickable: false,
            borderWidth: 0,
            aboveData: false
        },
        tooltip: true,
        tooltipOpts: {
          defaultTheme: false,
          content: "<span>$%y</span>",
        },
        colors: ['#eee', $scope.settings.color.themeprimary],
    };

    var clientRealtimeOpts = {
      series: {
        lines: {
          show: true
        },
        points: {
          show: true
        }
      },
      legend: {
        noColumns: 4
      },
      xaxis: {
        tickFormatter: function(val, axis) {
            return "";
        },
        color: '#eee'
      },
      yaxis: {
        min: 0,
        color: '#eee'
      },
      selection: {
        mode: "x"
      },
      grid: {
        hoverable: true,
        clickable: false,
        borderWidth: 0,
        aboveData: false
      },
      tooltip: true,
      tooltipOpts: {
        defaultTheme: false,
        content: "<strong>%s</strong> <span>$%y</span>",
      },
      crosshair: {
        mode: "x"
      }
    };

    var getSeriesObj;
    $scope.realTimedata = [];
    $scope.realTimeRevenue = [];

     function  getSeriesObj() {
        return [
            {
                data: getRandomData(),
                lines: {
                    show: true,
                    lineWidth: 1,
                    fill: true,
                    fillColor: {
                        colors: [
                            {
                                opacity: 0
                            }, {
                                opacity: 1
                            }
                        ]
                    },
                    steps: false
                },
                shadowSize: 0
            }, {
                data: getRandomData2(),
                lines: {
                    lineWidth: 0,
                    fill: true,
                    fillColor: {
                        colors: [
                            {
                                opacity: .5
                            }, {
                                opacity: 1
                            }
                        ]
                    },
                    steps: false
                },
                shadowSize: 0
            }
        ];
    }
    var getSeriesRevenueObj = function() {
        return [
            {
                data: getRevenueData(),
                lines: {
                    show: true,
                    lineWidth: 1,
                    fill: true,
                    fillColor: {
                        colors: [
                            {
                                opacity: 0
                            }, {
                                opacity: 1
                            }
                        ]
                    },
                    steps: false
                },
                shadowSize: 0
            }, {
                data: getRevenueData2(),
                lines: {
                    lineWidth: 0,
                    fill: true,
                    fillColor: {
                        colors: [
                            {
                                opacity: .5
                            }, {
                                opacity: 1
                            }
                        ]
                    },
                    steps: false
                },
                shadowSize: 0
            }
        ];
    }
    function getRandomData() {
        var res = [];

        if (typeof $scope.realTimedata !== 'undefined') {
          for (var i = 0; i < $scope.realTimedata.length; ++i) {
              res.push([i, $scope.realTimedata[i]]);
          }
        }

        return res;
    }
     function getRandomData2() {
      var res = [];

      if (typeof $scope.realTimedata2 !== 'undefined') {
        for (var i = 0; i < $scope.realTimedata2.length; ++i) {
            res.push([i, $scope.realTimedata2[i]]);
        }
      }

      return res;
    }
    var getRevenueData = function() {
        var res = [];

        for (var i = 0; i < $scope.realTimedata.length; ++i) {
            res.push([i, $scope.realTimeRevenue[i]]);
        }

        return res;
    }
    var getRevenueData2 = function() {
      var res = [];
      for (var i = 0; i < $scope.realTimedata2.length; ++i) {
          res.push([i, $scope.realTimeRevenue2[i]]);
      }

      return res;
    }

    var realtimeImpressionsOptions = {
        yaxis: {
            color: '#f3f3f3',
            min: 0
        },
        xaxis: {
            color: '#f3f3f3',
            min: 0,
            tickFormatter: function(val, axis) {
                return "";
            }
        },
        grid: {
            hoverable: true,
            clickable: true,
            borderWidth: 0,
            aboveData: false
        },
        tooltip: true,
        tooltipOpts: {
          defaultTheme: false,
          content: "<span>%y imps</span>",
        },
        colors: ['#eee', $scope.settings.color.themeprimary],
    }

    $scope.redirectQuickEdit = function(redirID, clientID) {
      $http.get("/adscoops/redirects/view/" + redirID).
      success(function(data) {
        $scope.clientRedirData["client_redir_editdata_" + clientID] = data;
      });
    }

    $scope.saveRedirectQuickview = function(qvRedirect) {
      qvRedirect.Min = parseInt(qvRedirect.Min);
      qvRedirect.Max = parseInt(qvRedirect.Max);
      qvRedirect.Iframe = parseInt(qvRedirect.Iframe);
      qvRedirect.AdvertisingDailySpend = parseInt(qvRedirect.AdvertisingDailySpend);
      qvRedirect.RedirType = parseInt(qvRedirect.RedirType);
      qvRedirect.BapiScoring = parseInt(qvRedirect.BapiScoring);

      qvRedirect.RedirType = parseInt(qvRedirect.RedirType);
      $http.post("/adscoops/redirects/save", qvRedirect).
      success(function() {

      })
    }

    $scope.clearRedirEditScreen = function(c) {
      $timeout(function() {
        $scope.clientRedirData['client_redir_editdata_' + c.ID] = {};
        $scope.clientRedirData['client_redir_chosen' + c.ID] = {};
      }, 100);
    }


    $scope.$on('$destroy', function() {
      $interval.cancel(realtimeInterval);
      $interval.cancel(fiveminUpdates);
    })

  }
]);
