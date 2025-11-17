/**
 * 登录页面逻辑
 */

// 刷新验证码
function refreshCaptcha() {
    const captchaImage = dom.$('#captchaImage');
    if (captchaImage) {
        captchaImage.src = '/getVerifyImg?' + new Date().getTime();
    }
}

// 登录表单提交
dom.on(document, 'DOMContentLoaded', function() {
    const loginForm = dom.$('#loginForm');

    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();

            // 获取表单数据
            const formData = dom.getFormData(loginForm);

            // 验证
            if (validator.isEmpty(formData.username)) {
                message.error('请输入用户名');
                return;
            }

            if (validator.isEmpty(formData.password)) {
                message.error('请输入密码');
                return;
            }

            if (validator.isEmpty(formData.captcha)) {
                message.error('请输入验证码');
                return;
            }

            try {
                // 发送登录请求
                const response = await http.post('/login', formData);

                if (response.code === '0000' || response.success) {
                    // 登录成功
                    message.success('登录成功');

                    // 保存用户信息
                    if (response.data) {
                        storage.set('userInfo', response.data);
                    }

                    // 跳转到主页
                    setTimeout(() => {
                        window.location.href = '/index.html';
                    }, 500);
                } else {
                    // 登录失败
                    message.error(response.message || '登录失败，请重试');
                    refreshCaptcha(); // 刷新验证码
                }
            } catch (error) {
                console.error('Login error:', error);
                message.error('网络错误，请稍后重试');
                refreshCaptcha();
            }
        });
    }

    // 回车键登录
    document.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            const submitBtn = loginForm.querySelector('button[type="submit"]');
            if (submitBtn) {
                submitBtn.click();
            }
        }
    });
});
