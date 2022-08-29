// UMD 魔法代码
// if the module has no dependencies, the above pattern can be simplified to
(function (root, factory) {
    if (typeof define === 'function' && define.amd) {
        // AMD. Register as an anonymous module.
        define([], factory);
    } else if (typeof module === 'object' && module.exports) {
        // Node. Does not work with strict CommonJS, but
        // only CommonJS-like environments that support module.exports,
        // like Node.
        module.exports = factory();
    } else {
        // Browser globals (root is window)
        root.handle = factory();
    }
}(this, function () {
    /**
     * 初始化需要的leaflet
     * @returns {leaflet_map}
     * @param domName
     * @param option
     */
    let registerEChartsOptions = function registerEchartsOption(domName, option) {
        let chartDom = document.getElementById(domName);
        let chart = echarts.init(chartDom);
        option && chart.setOption(option);
    }
    let registerEchartsWaterfall = function registerEchartsWaterfall(domName, chartTitle, xAxisData, allData, addData, subData) {
        let chartDom = document.getElementById(domName);
        let chart = echarts.init(chartDom);
        let option = {
            title: {
                text: chartTitle,
                subtext: ''
            },
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'shadow'
                },
                formatter: function (params) {
                    // var tar = params[1];
                    // let tar2 = params[2];
                    // return tar.name + '<br/>' + tar.seriesName + ' : ' + tar.value + '<br/>' +
                    //     tar2.seriesName + ' : ' + tar2.value;
                    var tar;
                    if (params[1].value !== "-") {
                        tar = params[1];
                    }
                    else {
                        tar = params[2];
                    }
                    // console.log(params);
                    return `${tar.name}<br/>${tar.seriesName} : ${tar.value}`;
                }
            },
            grid: {
                left: '3%',
                right: '4%',
                bottom: '3%',
                containLabel: true
            },
            xAxis: {
                type: 'category',
                splitLine: { show: false },
                data: xAxisData
            },
            yAxis: {
                    type: 'value'
                },
            dataZoom: [
                {
                    type: 'slider',
                    realtime: true,
                    start: 30,
                    end: 100,
                    xAxisIndex: [0]
                },
                {
                    type: 'inside',
                    realtime: true,
                    start: 30,
                    end: 100,
                    xAxisIndex: [0]
                },
            ],
            series: [
                {
                    name: '资产总量',
                    type: 'bar',
                    stack: 'Total',
                    itemStyle: {
                        barBorderColor: 'rgba(0,0,0,0)',
                        color: 'rgba(0,0,0,0)'
                    },
                    emphasis: {
                        itemStyle: {
                            barBorderColor: 'rgba(0,0,0,0)',
                            color: 'rgba(0,0,0,0)'
                        }
                    },
                    data: allData
                },
                {
                    name: '增量',
                    type: 'bar',
                    stack: 'Total',
                    label: {
                        show: true,
                        position: 'top'
                    },
                    itemStyle: {
                        color: 'rgba(220,38,38,100)',
                        barBorderColor: 'rgba(0,0,0,0)',
                        borderRadius: 4,
                        borderColor: '#fff',
                        borderWidth: 2,
                        shadowColor: 'rgba(0, 0, 0, 0.5)',
                        shadowBlur: 10
                    },
                    data: addData
                },
                {
                    name: '减量',
                    type: 'bar',
                    stack: 'Total',
                    label: {
                        show: true,
                        position: 'bottom'
                    },
                    itemStyle: {
                        color: 'rgba(5,150,105,100)',
                        barBorderColor: 'rgba(0,0,0,0)',
                        borderRadius: 4,
                        borderColor: '#fff',
                        borderWidth: 2,
                        shadowColor: 'rgba(0, 0, 0, 0.5)',
                        shadowBlur: 10
                    },
                    data: subData
                }
            ]
        };
        option && chart.setOption(option);
    }

    return {
        registerEChartsOption: registerEChartsOptions,
        registerEchartsWaterfall: registerEchartsWaterfall,
    }
}));