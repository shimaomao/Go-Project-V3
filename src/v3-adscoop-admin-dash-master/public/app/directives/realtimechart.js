angular.module('app').directive("adscoopsRealtimeChart", [
    "$http", "$interval", function($http, $interval) {
      var realTimedata,
          realTimedata2;

      updateValues();
      $interval(function() {
        updateValues();
      }, 5000);

        return {
            restrict: "AE",
            link: function(scope, ele) {
                var totalPoints,
                    getSeriesObj,
                    getRandomData,
                    getRandomData2,
                    updateInterval,
                    plot,
                    update,
                    updatedLength = 300;

                return realTimedata = [],
                    realTimedata2 = [],
                    totalPoints = updatedLength,
                    plotRevenue = $.plot($("#realtimeRevenue"), getSeriesRevenueObj(), realtimeRevenueOptions),
                    update = function() {

                        plotRedir.setData(getSeriesObj());
                        plotRevenue.setData(getSeriesRevenueObj());

                        plotRedir.draw();
                        plotRevenue.draw();
                        setTimeout(update, updateInterval);
                    },
                    update();
            }
        };
    }
]);

angular.module('app').directive("databoxFlotChartRealtime", [
    function() {
        return {
            restrict: "AE",
            link: function(scope, ele) {
                var data = [],
                    totalPoints = 300,
                    updateInterval = 100,
                    plot,
                    update,
                    getRandomData;
                return getRandomData = function() {

                        if (data.length > 0)
                            data = data.slice(1);

                        // Do a random walk

                        while (data.length < totalPoints) {

                            var prev = data.length > 0 ? data[data.length - 1] : 50,
                                y = prev + Math.random() * 10 - 5;

                            if (y < 0) {
                                y = 0;
                            } else if (y > 100) {
                                y = 100;
                            }

                            data.push(y);
                        }

                        // Zip the generated y values with the x values

                        var res = [];
                        for (var i = 0; i < data.length; ++i) {
                            res.push([i, data[i]]);
                        }

                        return res;
                    },

                    // Set up the control widget
                    plot = $.plot(ele[0], [getRandomData()], {
                        yaxis: {
                            color: '#f3f3f3',
                            min: 0,
                            max: 100,
                            tickFormatter: function(val, axis) {
                                return "";
                            }
                        },
                        xaxis: {
                            color: '#f3f3f3',
                            min: 0,
                            max: 100,
                            tickFormatter: function(val, axis) {
                                return "";
                            }
                        },
                        colors: ['#fff'],
                        series: {
                            lines: {
                                lineWidth: 2,
                                fill: false,
                                fillColor: {
                                    colors: [
                                        {
                                            opacity: 0.5
                                        }, {
                                            opacity: 0
                                        }
                                    ]
                                },
                                steps: false
                            },
                            shadowSize: 0
                        },
                        grid: {
                            show: false,
                            hoverable: true,
                            clickable: false,
                            borderWidth: 0,
                            aboveData: false
                        }
                    }), update = function() {

                        plot.setData([getRandomData()]);
                        plot.draw();
                        setTimeout(update, updateInterval);
                    },
                    update();
            }
        };
    }
]);
