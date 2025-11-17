/**
 * 工具函数库 - 纯 JavaScript 实现
 * 替代 jQuery 的常用功能
 */

// HTTP 请求工具
const http = {
    /**
     * GET 请求
     */
    get: async function(url, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const fullUrl = queryString ? `${url}?${queryString}` : url;

        try {
            const response = await fetch(fullUrl, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // 包含 cookie
            });

            return await this.handleResponse(response);
        } catch (error) {
            console.error('GET request failed:', error);
            throw error;
        }
    },

    /**
     * POST 请求
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
            console.error('POST request failed:', error);
            throw error;
        }
    },

    /**
     * 处理响应
     */
    handleResponse: async function(response) {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
            return await response.json();
        }

        return await response.text();
    }
};

// DOM 操作工具
const dom = {
    /**
     * 选择元素
     */
    $: function(selector) {
        return document.querySelector(selector);
    },

    /**
     * 选择所有元素
     */
    $$: function(selector) {
        return Array.from(document.querySelectorAll(selector));
    },

    /**
     * 添加事件监听
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
     * 显示元素
     */
    show: function(element) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.style.display = 'block';
        }
    },

    /**
     * 隐藏元素
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
     * 设置 HTML 内容
     */
    html: function(element, content) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.innerHTML = content;
        }
    },

    /**
     * 设置文本内容
     */
    text: function(element, content) {
        if (typeof element === 'string') {
            element = this.$(element);
        }
        if (element) {
            element.textContent = content;
        }
    },

    /**
     * 获取表单数据
     */
    getFormData: function(form) {
        if (typeof form === 'string') {
            form = this.$(form);
        }

        const formData = new FormData(form);
        const data = {};

        for (let [key, value] of formData.entries()) {
            data[key] = value;
        }

        return data;
    }
};

// 消息提示工具
const message = {
    /**
     * 显示错误消息
     */
    error: function(msg, elementId = 'errorMessage') {
        const element = dom.$(('#' + elementId));
        if (element) {
            dom.html(element, msg);
            dom.show(element);

            // 3秒后自动隐藏
            setTimeout(() => {
                dom.hide(element);
            }, 3000);
        } else {
            alert(msg);
        }
    },

    /**
     * 显示成功消息
     */
    success: function(msg) {
        alert(msg); // 简单实现，可以改为更美观的提示
    },

    /**
     * 显示加载状态
     */
    loading: function(show = true, elementId = 'loadingIndicator') {
        const element = dom.$('#' + elementId);
        if (element) {
            if (show) {
                dom.show(element);
            } else {
                dom.hide(element);
            }
        }
    }
};

// 验证工具
const validator = {
    /**
     * 验证是否为空
     */
    isEmpty: function(value) {
        return !value || value.trim() === '';
    },

    /**
     * 验证邮箱
     */
    isEmail: function(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    },

    /**
     * 验证手机号
     */
    isPhone: function(phone) {
        const re = /^1[3-9]\d{9}$/;
        return re.test(phone);
    },

    /**
     * 验证密码强度（至少6位）
     */
    isValidPassword: function(password) {
        return password && password.length >= 6;
    }
};

// 存储工具
const storage = {
    /**
     * 设置本地存储
     */
    set: function(key, value) {
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch (e) {
            console.error('Storage set error:', e);
        }
    },

    /**
     * 获取本地存储
     */
    get: function(key) {
        try {
            const value = localStorage.getItem(key);
            return value ? JSON.parse(value) : null;
        } catch (e) {
            console.error('Storage get error:', e);
            return null;
        }
    },

    /**
     * 删除本地存储
     */
    remove: function(key) {
        localStorage.removeItem(key);
    },

    /**
     * 清空本地存储
     */
    clear: function() {
        localStorage.clear();
    }
};

// 日期格式化工具
const dateUtils = {
    /**
     * 格式化日期
     */
    format: function(date, format = 'YYYY-MM-DD HH:mm:ss') {
        const d = new Date(date);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        const seconds = String(d.getSeconds()).padStart(2, '0');

        return format
            .replace('YYYY', year)
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hours)
            .replace('mm', minutes)
            .replace('ss', seconds);
    }
};
