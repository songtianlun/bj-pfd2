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
    let registerEchartsLinesChart = function registerEchartsLinesChart(domName, showDataZoom ,chartTitle,
                                                                       xAxisData,...dates) {
        let option = {
            title: {
                text: chartTitle
            },
            xAxis: {
                type: 'category',
                data: xAxisData
            },
            yAxis: {
                type: 'value'
            },
            tooltip: {
                trigger: 'axis'
            },
            series: []
        }
        if (showDataZoom) {
            option.dataZoom = [
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
            ]
        }
        for (let i = 0; i < dates.length; i++) {
            option.series.push({
                data: dates[i],
                type: 'line',
                smooth: true
            })
        }

        registerEChartsOptions(domName, option);
    }
    let registerEchartsOverView = function registerEchartsOverView(domName, chartTitle,
                                                                   firstTitle, firstData, secondTitle,secondData) {
        let option = {
            title: {
                text: chartTitle,
                subtext: '',
            },
            tooltip: {
                trigger: 'item',
                formatter: '{a} <br/>{b}: {c} ({d}%)',
                textStyle: {
                    fontWeight: 300,
                    fontSize: 10,
                    width: 4,
                },
            },
            series: [
                {
                    name: firstTitle,
                    type: 'pie',
                    selectedMode: 'single',
                    radius: [0, '30%'],
                    label: {
                        position: 'inner',
                        fontSize: 14
                    },
                    labelLine: {
                        show: false
                    },
                    itemStyle: {
                        borderRadius: 4,
                        // borderColor: '#fff',
                        borderWidth: 1,
                        shadowColor: 'rgba(0, 0, 0, 0.5)',
                        shadowBlur: 10
                    },
                    data: firstData
                },
                {
                    name: secondTitle,
                    type: 'pie',
                    radius: ['45%', '60%'],
                    labelLine: {
                        length: 30
                    },
                    label: {
                        formatter: '{a|{a}}{abg|}\n{hr|}\n  {b|{b}：}{c}  {per|{d}%}  ',
                        backgroundColor: '#F6F8FC',
                        borderColor: '#8C8D8E',
                        borderWidth: 1,
                        borderRadius: 4,
                        rich: {
                            a: {
                                color: '#6E7079',
                                lineHeight: 22,
                                align: 'center'
                            },
                            hr: {
                                borderColor: '#8C8D8E',
                                width: '100%',
                                borderWidth: 1,
                                height: 0
                            },
                            b: {
                                color: '#4C5058',
                                fontSize: 14,
                                fontWeight: 'bold',
                                lineHeight: 33
                            },
                            per: {
                                color: '#fff',
                                backgroundColor: '#4C5058',
                                padding: [3, 4],
                                borderRadius: 4
                            }
                        }
                    },
                    itemStyle: {
                        borderRadius: 4,
                        // borderColor: '#fff',
                        borderWidth: 1,
                        shadowColor: 'rgba(0, 0, 0, 0.5)',
                        shadowBlur: 10
                    },
                    data: secondData
                }
            ]
        }
        registerEChartsOptions(domName, option);
    }
    let registerEchartsWaterfall = function registerEchartsWaterfall(domName, showDataZoom, chartTitle,
                                                                     xAxisData, allData, addData, subData) {
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
        if (showDataZoom) {
            option.dataZoom =  [
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
            ]
        }
        registerEChartsOptions(domName, option);
    }

    let checkHttps = function checkHttps(access_local=true) {
        if (location.protocol !== "https:") {
            if ((location.hostname === "localhost" ||
                    location.hostname === "127.0.0.1") &&
                access_local === true) {
                return true
            }
            let r = confirm("使用 HTTP 协议传输将导致您的隐私信息被泄漏，帮您重定向到 HTTPS ？");
            if (r === true) {
                window.location.href = "https:" + window.location.href.substring(window.location.protocol.length);
            } else {
                alert("禁止使用 HTTPS 协议传输私密数据！");
            }
            return false
        } else {
            return true
        }
    }


    return {
        registerEChartsOption: registerEChartsOptions,
        registerEchartsWaterfall: registerEchartsWaterfall,
        registerEchartsOverView: registerEchartsOverView,
        registerEchartsLinesChart: registerEchartsLinesChart,
        checkHttps: checkHttps,
    }
}));