/**
 * dongfeng-pay Boss 登录页面逻辑
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
        refreshCaptcha();

        // 绑定事件
        bindEvents();

        // 尝试加载记住的密码
        loadRememberedCredentials();
    }

    // ===== 缓存 DOM 元素 =====
    function cacheElements() {
        elements = {
            loginForm: dom.$('#loginForm'),
            userIDInput: dom.$('#userID'),
            passwdInput: dom.$('#passwd'),
            verifyCodeInput: dom.$('#verifyCode'),
            captchaImage: dom.$('#captchaImage'),
            rememberMeCheckbox: dom.$('#rememberMe'),
            loginBtn: dom.$('#loginBtn'),
            btnText: dom.$('#btnText'),
            btnLoading: dom.$('#btnLoading'),

            // 错误提示元素
            userIDError: dom.$('#userIDError'),
            passwdError: dom.$('#passwdError'),
            codeError: dom.$('#codeError'),
        };
    }

    // ===== 绑定事件 =====
    function bindEvents() {
        // 表单提交事件
        elements.loginForm.addEventListener('submit', handleLogin);

        // 验证码图片点击事件
        elements.captchaImage.addEventListener('click', refreshCaptcha);

        // 输入框焦点事件（清除错误提示）
        elements.userIDInput.addEventListener('focus', () => clearError('userID'));
        elements.passwdInput.addEventListener('focus', () => clearError('passwd'));
        elements.verifyCodeInput.addEventListener('focus', () => clearError('code'));

        // Enter 键提交
        [elements.userIDInput, elements.passwdInput, elements.verifyCodeInput].forEach(input => {
            input.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    handleLogin(e);
                }
            });
        });
    }

    // ===== 刷新验证码 =====
    function refreshCaptcha() {
        const timestamp = new Date().getTime();
        elements.captchaImage.src = `/getVerifyImg?rand=${timestamp}`;
    }

    // ===== 清除错误提示 =====
    function clearError(field) {
        const errorMap = {
            'userID': elements.userIDError,
            'passwd': elements.passwdError,
            'code': elements.codeError,
        };

        const errorElement = errorMap[field];
        if (errorElement) {
            errorElement.textContent = '';
            errorElement.style.display = 'none';
        }
    }

    // ===== 显示错误提示 =====
    function showError(field, message) {
        const errorMap = {
            'userID': elements.userIDError,
            'passwd': elements.passwdError,
            'code': elements.codeError,
        };

        const errorElement = errorMap[field];
        if (errorElement) {
            errorElement.textContent = message;
            errorElement.style.display = 'block';
        }

        // 刷新验证码
        refreshCaptcha();
    }

    // ===== 表单验证 =====
    function validateForm(formData) {
        // 清除所有错误提示
        clearError('userID');
        clearError('passwd');
        clearError('code');

        let isValid = true;

        // 验证手机号
        if (!formData.userID || formData.userID.length === 0) {
            showError('userID', '* 登录手机号不能为空！');
            isValid = false;
        }

        // 验证密码
        if (!formData.passwd || formData.passwd.length === 0) {
            showError('passwd', '* 密码不能为空！');
            isValid = false;
        }

        // 验证验证码
        if (!formData.Code || formData.Code.length < 4) {
            showError('code', '* 验证码不正确！');
            isValid = false;
        }

        return isValid;
    }

    // ===== 设置加载状态 =====
    function setLoadingState(loading) {
        if (loading) {
            elements.loginBtn.disabled = true;
            dom.hide(elements.btnText);
            dom.show(elements.btnLoading);
        } else {
            elements.loginBtn.disabled = false;
            dom.show(elements.btnText);
            dom.hide(elements.btnLoading);
        }
    }

    // ===== 处理登录 =====
    async function handleLogin(e) {
        e.preventDefault();

        // 获取表单数据
        const userID = dom.getValue(elements.userIDInput);
        const passwd = dom.getValue(elements.passwdInput);
        const Code = dom.getValue(elements.verifyCodeInput);

        const formData = {
            userID: userID,
            passwd: passwd,
            Code: Code,
        };

        // 验证表单
        if (!validateForm(formData)) {
            return;
        }

        // 设置加载状态
        setLoadingState(true);

        try {
            // 发送登录请求
            const response = await http.postForm('/login', formData);

            // 处理响应
            handleLoginResponse(response, {
                userID: userID,
                passwd: passwd,
            });
        } catch (error) {
            console.error('登录请求失败:', error);
            message.error('系统异常，请稍后再尝试！');
            refreshCaptcha();
        } finally {
            setLoadingState(false);
        }
    }

    // ===== 处理登录响应 =====
    function handleLoginResponse(data, credentials) {
        if (!data) {
            message.error('登录失败，响应数据为空');
            refreshCaptcha();
            return;
        }

        // 根据返回的 Key 字段判断错误类型
        if (data.Key) {
            switch (data.Key) {
                case 'userID':
                    showError('userID', data.Msg || '* 用户名错误！');
                    break;
                case 'passWD':
                    showError('passwd', data.Msg || '* 密码错误！');
                    break;
                case 'code':
                    showError('code', data.Msg || '* 验证码错误！');
                    break;
                case 'unactive':
                case 'del':
                    message.error(data.Msg || '账户状态异常');
                    refreshCaptcha();
                    break;
                default:
                    // Key 为空字符串表示登录成功
                    if (data.Key.length === 0) {
                        handleLoginSuccess(credentials);
                    } else {
                        message.error(data.Msg || '登录失败');
                        refreshCaptcha();
                    }
            }
        } else {
            // 没有 Key 字段，可能是成功
            handleLoginSuccess(credentials);
        }
    }

    // ===== 登录成功处理 =====
    function handleLoginSuccess(credentials) {
        message.success('登录成功，正在跳转...');

        // 如果勾选了"记住密码"，保存凭据
        if (elements.rememberMeCheckbox.checked) {
            saveCredentials(credentials);
        } else {
            clearSavedCredentials();
        }

        // 延迟跳转到首页
        setTimeout(() => {
            window.location.href = '/index.html';
        }, 500);
    }

    // ===== 保存凭据 =====
    function saveCredentials(credentials) {
        storage.set('boss_remember', {
            userID: credentials.userID,
            passwd: credentials.passwd,
            timestamp: Date.now(),
        });
    }

    // ===== 清除保存的凭据 =====
    function clearSavedCredentials() {
        storage.remove('boss_remember');
    }

    // ===== 加载记住的凭据 =====
    function loadRememberedCredentials() {
        const remembered = storage.get('boss_remember');

        if (!remembered) {
            return;
        }

        // 检查是否过期（7天）
        const MAX_AGE = 7 * 24 * 60 * 60 * 1000; // 7天
        const age = Date.now() - remembered.timestamp;

        if (age > MAX_AGE) {
            clearSavedCredentials();
            return;
        }

        // 填充表单
        dom.setValue(elements.userIDInput, remembered.userID);
        dom.setValue(elements.passwdInput, remembered.passwd);
        elements.rememberMeCheckbox.checked = true;
    }

    // ===== 页面加载完成后初始化 =====
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
