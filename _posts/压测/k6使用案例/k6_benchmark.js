import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const requestDuration = new Trend('request_duration');

// 配置参数 - 可以通过环境变量修改
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8092/sdk/spi/test';
const APP_KEY = __ENV.APP_KEY || '3409409348479354011';
const TIMESTAMP = __ENV.TIMESTAMP || Date.now();
const SIGN = __ENV.SIGN || '8abb21bcfc4cc7ba4a501e2dc73a5e0c';

// 为旧版 K6 兼容的 UUID 生成器
function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// 从文件加载测试数据
const testData = new SharedArray('testData', function () {
    try {
        const fileContent = open('./new2new.jsonl');
        const lines = fileContent.split('\n').filter(line => line.trim() !== '');
        
        return lines.map((line, index) => {
            try {
                return JSON.parse(line);
            } catch (e) {
                console.warn(`解析第 ${index + 1} 行失败: ${e.message}`);
                return null;
            }
        }).filter(item => item !== null);
    } catch (e) {
        console.error(`无法读取文件 new2new.jsonl: ${e.message}`);
        console.error('请确保 new2new.jsonl 文件在脚本同目录下');
        return [];
    }
});

// 生成请求 URL
function buildURL(baseUrl, appKey, timestamp, sign) {
    try {
        const url = new URL(baseUrl);
        url.searchParams.set('app_key', appKey);
        url.searchParams.set('timestamp', timestamp);
        url.searchParams.set('sign', sign);
        return url.toString();
    } catch (e) {
        const separator = baseUrl.includes('?') ? '&' : '?';
        return `${baseUrl}${separator}app_key=${encodeURIComponent(appKey)}&timestamp=${encodeURIComponent(timestamp)}&sign=${encodeURIComponent(sign)}`;
    }
}

// 构建请求体
function buildRequestBody(dataItem) {
    const requestBody = JSON.parse(JSON.stringify(dataItem));
    requestBody.timestamp = Date.now(); // 使用毫秒级时间戳

    if (requestBody.evaluation_data_list && Array.isArray(requestBody.evaluation_data_list)) {
        requestBody.evaluation_data_list.forEach(item => {
            if (item) {
                item.data_id = uuidv4(); // 使用自定义的UUID生成器
            }
        });
    }

    return JSON.stringify(requestBody);
}

// 压测配置: 固定QPS测试
export const options = {
    scenarios: {
        constant_qps_scenario: {
            executor: 'constant-arrival-rate',
            rate: 5,
            timeUnit: '1s',  
            duration: '1m',  // 总持续时间
            preAllocatedVUs: 20,  // 预分配的虚拟用户数
            maxVUs: 50, // 最大虚拟用户数
        },
    },
    thresholds: {
        'http_req_duration': ['p(95)<2000'],
        'http_req_failed': ['rate<0.01'],
    },
};

// 主压测函数
export default function () {
    if (testData.length === 0) {
        console.error('没有可用的测试数据，跳过请求');
        return;
    }
    
    const dataIndex = (__VU * 1000 + __ITER) % testData.length;
    const dataItem = testData[dataIndex];
    
    const url = buildURL(BASE_URL, APP_KEY, TIMESTAMP, SIGN);
    const payload = buildRequestBody(dataItem);
    
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
        tags: {
            name: 'api_test',
            data_id: dataItem.evaluation_data_list?.[0]?.data_id || 'unknown',
        },
    };
    
    const response = http.post(url, payload, params);
    requestDuration.add(response.timings.duration);
    
    const success = check(response, {
        '状态码是 200': (r) => r.status === 200,
        '响应时间 < 5000ms': (r) => r.timings.duration < 5000,
        '响应体不为空': (r) => r.body && r.body.length > 0,
    });
    
    if (!success || response.status !== 200) {
        errorRate.add(1);
        if (response.status !== 200) {
            console.error(`请求失败: status=${response.status}, data_id=${dataItem.evaluation_data_list?.[0]?.data_id || 'unknown'}`);
            if (response.body) {
                console.error(`响应体: ${response.body.substring(0, 200)}`);
            }
        }
    } else {
        errorRate.add(0);
    }
    
    sleep(0.1);
}

// 测试设置函数（可选）
export function setup() {
    console.log('========== K6 压测开始 ==========');
    console.log(`目标URL: ${BASE_URL}`);
    console.log(`APP_KEY: ${APP_KEY}`);
    console.log(`测试数据条数: ${testData.length}`);
    console.log('================================');
    
    return {
        dataCount: testData.length,
        startTime: new Date().toISOString(),
    };
}

// 测试结束函数（可选）
export function teardown(data) {
    console.log('\n========== K6 压测结束 ==========');
    console.log(`测试开始时间: ${data.startTime}`);
    console.log(`测试数据条数: ${data.dataCount}`);
    console.log('================================\n');
}
