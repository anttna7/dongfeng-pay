/**
 * dongfeng-pay Merchant 登录页面逻辑
 * 纯 JavaScript 实现，无需 jQuery
 */

(function() {
    'use strict';

    // ===== 页面元素引用 =====
    let elements = {};

    // ===== 初始化函数 =====
    function init() {
        // 获取 DOM 元素引用
        cacheElements();

        // 加载验证码
        loadCaptcha();

        // 绑定事件
        bindEvents();
    }

    // ===== 缓存 DOM 元素 =====
    function cacheElements() {
        elements = {
            loginForm: dom.$('#loginForm'),
            userNameInput: dom.$('#userName'),
            passwordInput: dom.$('#password'),
            captchaCodeInput: dom.$('#captchaCode'),
            captchaIdInput: dom.$('#captchaId'),
            captchaImage: dom.$('#captchaImage'),
            loginBtn: dom.$('#loginBtn'),
            btnIcon: dom.$('#btnIcon'),
            btnText: dom.$('#btnText'),
            btnLoading: dom.$('#btnLoading'),

            // 错误提示元素
            userNameError: dom.$('#userNameError'),
            passwordError: dom.$('#passwordError'),
            captchaError: dom.$('#captchaError'),
        };
    }

    // ===== 绑定事件 =====
    function bindEvents() {
        // 表单提交事件
        elements.loginForm.addEventListener('submit', handleLogin);

        // 验证码图片点击事件
        elements.captchaImage.addEventListener('click', refreshCaptcha);

        // 输入框焦点事件（清除错误提示）
        elements.userNameInput.addEventListener('focus', () => clearError('userName'));
        elements.passwordInput.addEventListener('focus', () => clearError('password'));
        elements.captchaCodeInput.addEventListener('focus', () => clearError('captcha'));

        // Enter 键提交
        [elements.userNameInput, elements.passwordInput, elements.captchaCodeInput].forEach(input => {
            input.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    handleLogin(e);
                }
            });
        });
    }

    // ===== 加载验证码 =====
    function loadCaptcha() {
        const captchaId = dom.getValue(elements.captchaIdInput);
        if (captchaId) {
            elements.captchaImage.src = `/img.do/${captchaId}.png`;
        }
    }

    // ===== 刷新验证码 =====
    async function refreshCaptcha() {
        try {
            // 调用刷新验证码接口
            const response = await http.get('/flushCaptcha.py/');

            if (response && response.captchaId) {
                // 更新验证码 ID
                dom.setValue(elements.captchaIdInput, response.captchaId);

                // 更新验证码图片
                elements.captchaImage.src = `/img.do/${response.captchaId}.png?t=${Date.now()}`;

                // 清空验证码输入
                dom.setValue(elements.captchaCodeInput, '');
            } else {
                // 如果没有返回新的 ID，直接刷新图片
                const currentId = dom.getValue(elements.captchaIdInput);
                if (currentId) {
                    elements.captchaImage.src = `/img.do/${currentId}.png?t=${Date.now()}`;
                }
            }
        } catch (error) {
            console.error('刷新验证码失败:', error);

            // 降级方案：直接刷新图片
            const currentId = dom.getValue(elements.captchaIdInput);
            if (currentId) {
                elements.captchaImage.src = `/img.do/${currentId}.png?t=${Date.now()}`;
            }
        }
    }

    // ===== 清除错误提示 =====
    function clearError(field) {
        const errorMap = {
            'userName': elements.userNameError,
            'password': elements.passwordError,
            'captcha': elements.captchaError,
        };

        const errorElement = errorMap[field];
        if (errorElement) {
            errorElement.textContent = '';
            errorElement.style.display = 'none';
        }
    }

    // ===== 显示错误提示 =====
    function showError(field, errorMessage) {
        const errorMap = {
            'userName': elements.userNameError,
            'password': elements.passwordError,
            'captcha': elements.captchaError,
        };

        const errorElement = errorMap[field];
        if (errorElement) {
            errorElement.textContent = errorMessage;
            errorElement.style.display = 'block';
        }

        // 刷新验证码
        refreshCaptcha();
    }

    // ===== 表单验证 =====
    function validateForm(formData) {
        // 清除所有错误提示
        clearError('userName');
        clearError('password');
        clearError('captcha');

        let isValid = true;

        // 验证用户名
        if (!formData.userName || formData.userName.length === 0) {
            showError('userName', '* 请输入您的账户！');
            isValid = false;
        }

        // 验证密码
        if (!formData.password || formData.password.length === 0) {
            showError('password', '* 请输入您的密码！');
            isValid = false;
        }

        // 验证验证码
        if (!formData.captchaCode || formData.captchaCode.length === 0) {
            showError('captcha', '* 请输入验证码！');
            isValid = false;
        }

        // 验证验证码 ID
        if (!formData.captchaId || formData.captchaId.length === 0) {
            showError('captcha', '* 验证码已过期，请刷新！');
            isValid = false;
        }

        return isValid;
    }

    // ===== 设置加载状态 =====
    function setLoadingState(loading) {
        if (loading) {
            elements.loginBtn.disabled = true;
            dom.hide(elements.btnIcon);
            dom.hide(elements.btnText);
            dom.show(elements.btnLoading);
        } else {
            elements.loginBtn.disabled = false;
            dom.show(elements.btnIcon);
            dom.show(elements.btnText);
            dom.hide(elements.btnLoading);
        }
    }

    // ===== 处理登录 =====
    async function handleLogin(e) {
        e.preventDefault();

        // 获取表单数据
        const userName = dom.getValue(elements.userNameInput);
        const password = dom.getValue(elements.passwordInput);
        const captchaCode = dom.getValue(elements.captchaCodeInput);
        const captchaId = dom.getValue(elements.captchaIdInput);

        const formData = {
            userName: userName,
            password: password,
            captchaCode: captchaCode,
            captchaId: captchaId,
        };

        // 验证表单
        if (!validateForm(formData)) {
            return;
        }

        // 设置加载状态
        setLoadingState(true);

        try {
            // 构建请求参数
            const params = {
                userName: userName,
                Password: password,  // 注意：后端期望的是大写 P
                captchaCode: captchaCode,
                captchaId: captchaId,
            };

            // 发送登录请求
            const response = await http.get('/login.py/', params);

            // 处理响应
            handleLoginResponse(response);
        } catch (error) {
            console.error('登录请求失败:', error);
            message.error('系统异常，请稍后再尝试！');
            refreshCaptcha();
        } finally {
            setLoadingState(false);
        }
    }

    // ===== 处理登录响应 =====
    function handleLoginResponse(data) {
        if (!data) {
            message.error('登录失败，响应数据为空');
            refreshCaptcha();
            return;
        }

        // 检查响应状态
        if (data.status === 'success' || data.code === 200 || data.success === true) {
            handleLoginSuccess();
        } else if (data.status === 'error' || data.code !== 200) {
            // 根据错误类型显示不同的错误信息
            const errorMsg = data.message || data.msg || '登录失败，请检查您的账户信息';

            // 尝试判断错误类型
            if (errorMsg.includes('账户') || errorMsg.includes('用户')) {
                showError('userName', errorMsg);
            } else if (errorMsg.includes('密码')) {
                showError('password', errorMsg);
            } else if (errorMsg.includes('验证码')) {
                showError('captcha', errorMsg);
            } else {
                message.error(errorMsg);
                refreshCaptcha();
            }
        } else {
            // 未明确状态，可能是成功
            // 检查是否有重定向信息
            if (data.redirectUrl || data.redirect) {
                handleLoginSuccess(data.redirectUrl || data.redirect);
            } else {
                // 默认成功处理
                handleLoginSuccess();
            }
        }
    }

    // ===== 登录成功处理 =====
    function handleLoginSuccess(redirectUrl) {
        message.success('登录成功，正在跳转...');

        // 延迟跳转
        setTimeout(() => {
            if (redirectUrl) {
                window.location.href = redirectUrl;
            } else {
                // 默认跳转到首页
                window.location.href = '/index/ui/';
            }
        }, 500);
    }

    // ===== 验证码验证（可选功能）=====
    async function verifyCaptcha() {
        const captchaCode = dom.getValue(elements.captchaCodeInput);
        const captchaId = dom.getValue(elements.captchaIdInput);

        if (!captchaCode || !captchaId) {
            return false;
        }

        try {
            const response = await http.get(`/verifyCaptcha.py/${captchaCode}/${captchaId}`);
            return response && response.valid === true;
        } catch (error) {
            console.error('验证码验证失败:', error);
            return false;
        }
    }

    // ===== 页面加载完成后初始化 =====
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
