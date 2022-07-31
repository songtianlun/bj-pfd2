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

    return {
        registerEChartsOption: registerEChartsOptions,
    }
}));