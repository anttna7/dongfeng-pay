/**
 * dongfeng-pay 工具函数库
 * 纯 JavaScript 实现，替代 jQuery 和其他框架
 */

/* ===== HTTP 请求工具 (替代 jQuery.ajax) ===== */
const http = {
    /**
     * GET 请求
     * @param {string} url - 请求 URL
     * @param {object} params - 查询参数
     * @returns {Promise<any>} 响应数据
     */
    get: async function(url, params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const fullUrl = queryString ? `${url}?${queryString}` : url;

            const response = await fetch(fullUrl, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // 携带 Cookie
            });

            return await this.handleResponse(response);
        } catch (error) {
            console.error('GET 请求失败:', error);
            throw error;
        }
    },

    /**
     * POST 请求
     * @param {string} url - 请求 URL
     * @param {object} data - 请求数据
     * @returns {Promise<any>} 响应数据
     */
    post: async function(url, data = {}) {
        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify(data),
            });

            return await this.handleResponse(response);
        } catch (error) {
            console.error('POST 请求失败:', error);
            throw error;
        }
    },

    /**
     * POST 表单数据 (application/x-www-form-urlencoded)
     * @param {string} url - 请求 URL
     * @param {object} data - 表单数据
     * @returns {Promise<any>} 响应数据
     */
    postForm: async function(url, data = {}) {
        try {
            const formData = new URLSearchParams();
            for (let key in data) {
                formData.append(key, data[key]);
            }

            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                credentials: 'include',
                body: formData.toString(),
            });

            return await this.handleResponse(response);
        } catch (error) {
            console.error('POST Form 请求失败:', error);
            throw error;
        }
    },

    /**
     * 处理响应
     * @param {Response} response - Fetch Response 对象
     * @returns {Promise<any>} 处理后的数据
     */
    handleResponse: async function(response) {
        const contentType = response.headers.get('content-type');

        // 处理 JSON 响应
        if (contentType && contentType.includes('application/json')) {
            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.message || `HTTP ${response.status}`);
            }
            return data;
        }

        // 处理文本响应
        const text = await response.text();
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${text}`);
        }

        // 尝试解析为 JSON
        try {
            return JSON.parse(text);
        } catch {
            return text;
        }
    },
};

/* ===== DOM 操作工具 (替代 jQuery) ===== */
const dom = {
    /**
     * 查询单个元素
     * @param {string} selector - CSS 选择器
     * @param {Element} context - 上下文元素
     * @returns {Element|null}
     */
    $: function(selector, context = document) {
        return context.querySelector(selector);
    },

    /**
     * 查询所有元素
     * @param {string} selector - CSS 选择器
     * @param {Element} context - 上下文元素
     * @returns {NodeList}
     */
    $$: function(selector, context = document) {
        return context.querySelectorAll(selector);
    },

    /**
     * 添加事件监听
     * @param {string|Element} element - 元素或选择器
     * @param {string} event - 事件名称
     * @param {Function} handler - 事件处理函数
     */
    on: function(element, event, handler) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.addEventListener(event, handler);
        }
    },

    /**
     * 移除事件监听
     * @param {string|Element} element - 元素或选择器
     * @param {string} event - 事件名称
     * @param {Function} handler - 事件处理函数
     */
    off: function(element, event, handler) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.removeEventListener(event, handler);
        }
    },

    /**
     * 获取表单数据
     * @param {string|HTMLFormElement} form - 表单元素或选择器
     * @returns {object} 表单数据对象
     */
    getFormData: function(form) {
        if (typeof form === 'string') {
            form = this.$(form);
        }
        if (!form) return {};

        const formData = new FormData(form);
        const data = {};
        for (let [key, value] of formData.entries()) {
            data[key] = value;
        }
        return data;
    },

    /**
     * 显示元素
     * @param {string|Element} element - 元素或选择器
     */
    show: function(element) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.style.display = '';
        }
    },

    /**
     * 隐藏元素
     * @param {string|Element} element - 元素或选择器
     */
    hide: function(element) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.style.display = 'none';
        }
    },

    /**
     * 添加 CSS 类
     * @param {string|Element} element - 元素或选择器
     * @param {string} className - 类名
     */
    addClass: function(element, className) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.classList.add(className);
        }
    },

    /**
     * 移除 CSS 类
     * @param {string|Element} element - 元素或选择器
     * @param {string} className - 类名
     */
    removeClass: function(element, className) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.classList.remove(className);
        }
    },

    /**
     * 设置元素文本内容
     * @param {string|Element} element - 元素或选择器
     * @param {string} text - 文本内容
     */
    setText: function(element, text) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.textContent = text;
        }
    },

    /**
     * 设置元素 HTML 内容
     * @param {string|Element} element - 元素或选择器
     * @param {string} html - HTML 内容
     */
    setHTML: function(element, html) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.innerHTML = html;
        }
    },

    /**
     * 获取元素值
     * @param {string|Element} element - 元素或选择器
     * @returns {string} 元素值
     */
    getValue: function(element) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        return element ? element.value.trim() : '';
    },

    /**
     * 设置元素值
     * @param {string|Element} element - 元素或选择器
     * @param {string} value - 值
     */
    setValue: function(element, value) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.value = value;
        }
    },
};

/* ===== 消息提示工具 ===== */
const message = {
    /**
     * 显示成功消息
     * @param {string} text - 消息文本
     * @param {number} duration - 持续时间（毫秒）
     */
    success: function(text, duration = 3000) {
        this.show(text, 'success', duration);
    },

    /**
     * 显示错误消息
     * @param {string} text - 消息文本
     * @param {number} duration - 持续时间（毫秒）
     */
    error: function(text, duration = 3000) {
        this.show(text, 'error', duration);
    },

    /**
     * 显示警告消息
     * @param {string} text - 消息文本
     * @param {number} duration - 持续时间（毫秒）
     */
    warning: function(text, duration = 3000) {
        this.show(text, 'warning', duration);
    },

    /**
     * 显示信息消息
     * @param {string} text - 消息文本
     * @param {number} duration - 持续时间（毫秒）
     */
    info: function(text, duration = 3000) {
        this.show(text, 'info', duration);
    },

    /**
     * 显示消息
     * @param {string} text - 消息文本
     * @param {string} type - 消息类型
     * @param {number} duration - 持续时间（毫秒）
     */
    show: function(text, type = 'info', duration = 3000) {
        // 创建消息元素
        const messageEl = document.createElement('div');
        messageEl.className = `message message-${type}`;
        messageEl.textContent = text;
        messageEl.style.cssText = `
            position: fixed;
            top: 20px;
            left: 50%;
            transform: translateX(-50%);
            padding: 12px 24px;
            background: ${this.getBackgroundColor(type)};
            color: white;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
            z-index: 9999;
            font-size: 14px;
            animation: slideDown 0.3s ease-out;
        `;

        // 添加动画样式
        if (!document.querySelector('#message-animation-style')) {
            const style = document.createElement('style');
            style.id = 'message-animation-style';
            style.textContent = `
                @keyframes slideDown {
                    from {
                        opacity: 0;
                        transform: translateX(-50%) translateY(-20px);
                    }
                    to {
                        opacity: 1;
                        transform: translateX(-50%) translateY(0);
                    }
                }
            `;
            document.head.appendChild(style);
        }

        // 添加到页面
        document.body.appendChild(messageEl);

        // 自动移除
        setTimeout(() => {
            messageEl.style.animation = 'slideDown 0.3s ease-out reverse';
            setTimeout(() => {
                document.body.removeChild(messageEl);
            }, 300);
        }, duration);
    },

    /**
     * 获取背景颜色
     * @param {string} type - 消息类型
     * @returns {string} 颜色值
     */
    getBackgroundColor: function(type) {
        const colors = {
            success: '#16a34a',
            error: '#dc2626',
            warning: '#ea580c',
            info: '#2563eb',
        };
        return colors[type] || colors.info;
    },
};

/* ===== 本地存储工具 ===== */
const storage = {
    /**
     * 设置存储
     * @param {string} key - 键
     * @param {any} value - 值
     */
    set: function(key, value) {
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch (error) {
            console.error('存储失败:', error);
        }
    },

    /**
     * 获取存储
     * @param {string} key - 键
     * @param {any} defaultValue - 默认值
     * @returns {any} 存储的值
     */
    get: function(key, defaultValue = null) {
        try {
            const value = localStorage.getItem(key);
            return value ? JSON.parse(value) : defaultValue;
        } catch (error) {
            console.error('读取存储失败:', error);
            return defaultValue;
        }
    },

    /**
     * 移除存储
     * @param {string} key - 键
     */
    remove: function(key) {
        try {
            localStorage.removeItem(key);
        } catch (error) {
            console.error('移除存储失败:', error);
        }
    },

    /**
     * 清空存储
     */
    clear: function() {
        try {
            localStorage.clear();
        } catch (error) {
            console.error('清空存储失败:', error);
        }
    },
};

/* ===== 工具函数 ===== */
const utils = {
    /**
     * 防抖函数
     * @param {Function} func - 目标函数
     * @param {number} wait - 等待时间（毫秒）
     * @returns {Function} 防抖后的函数
     */
    debounce: function(func, wait = 300) {
        let timeout;
        return function(...args) {
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(this, args), wait);
        };
    },

    /**
     * 节流函数
     * @param {Function} func - 目标函数
     * @param {number} wait - 等待时间（毫秒）
     * @returns {Function} 节流后的函数
     */
    throttle: function(func, wait = 300) {
        let timeout;
        return function(...args) {
            if (!timeout) {
                timeout = setTimeout(() => {
                    timeout = null;
                    func.apply(this, args);
                }, wait);
            }
        };
    },

    /**
     * 格式化日期
     * @param {Date|string|number} date - 日期
     * @param {string} format - 格式
     * @returns {string} 格式化后的日期
     */
    formatDate: function(date, format = 'YYYY-MM-DD HH:mm:ss') {
        const d = new Date(date);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hour = String(d.getHours()).padStart(2, '0');
        const minute = String(d.getMinutes()).padStart(2, '0');
        const second = String(d.getSeconds()).padStart(2, '0');

        return format
            .replace('YYYY', year)
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hour)
            .replace('mm', minute)
            .replace('ss', second);
    },

    /**
     * 生成随机字符串
     * @param {number} length - 长度
     * @returns {string} 随机字符串
     */
    randomString: function(length = 8) {
        const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        let result = '';
        for (let i = 0; i < length; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return result;
    },
};
