/**
 * dongfeng-pay 管理后台布局逻辑
 * 纯 JavaScript ES6+ 实现，无框架依赖
 */

(function() {
    'use strict';

    // ===== DOM 元素 =====
    const elements = {
        // 用户菜单
        userMenuBtn: document.getElementById('userMenuBtn'),
        userMenu: document.getElementById('userMenu'),
        userDropdown: document.querySelector('.user-dropdown'),

        // 侧边栏
        sidebar: document.getElementById('sidebar'),
        menuParents: document.querySelectorAll('.menu-item-parent'),
        menuItems: document.querySelectorAll('.menu-item'),
        submenuItems: document.querySelectorAll('.submenu-item'),

        // 主内容区
        mainContent: document.getElementById('mainContent'),

        // 按钮
        changePasswordBtn: document.getElementById('changePasswordBtn'),
        logoutBtn: document.getElementById('logoutBtn'),

        // 密码模态框
        passwordModal: document.getElementById('passwordModal'),
        closePasswordModal: document.getElementById('closePasswordModal'),
        cancelPasswordBtn: document.getElementById('cancelPasswordBtn'),
        savePasswordBtn: document.getElementById('savePasswordBtn'),
        passwordForm: document.getElementById('passwordForm'),
        oldPassword: document.getElementById('oldPassword'),
        newPassword: document.getElementById('newPassword'),
        confirmPassword: document.getElementById('confirmPassword'),
        oldPasswordError: document.getElementById('oldPasswordError'),
        newPasswordError: document.getElementById('newPasswordError'),
        confirmPasswordError: document.getElementById('confirmPasswordError')
    };

    // ===== 状态管理 =====
    const state = {
        currentPage: null,
        activeMenu: null
    };

    // ===== 工具函数 =====
    const utils = {
        // 清空错误提示
        clearErrors() {
            elements.oldPasswordError.textContent = '';
            elements.newPasswordError.textContent = '';
            elements.confirmPasswordError.textContent = '';
        },

        // 设置错误提示
        setError(element, message) {
            element.textContent = message;
        },

        // 密码验证
        validatePassword(password) {
            if (password.length < 8) {
                return '密码长度不能小于8！';
            }
            if (password.length > 16) {
                return '密码长度不能大于16！';
            }
            if (!/([a-zA-Z]+[0-9]+|[0-9]+[a-zA-Z])/.test(password)) {
                return '新密码中必须包含数字和字母！';
            }
            return '';
        }
    };

    // ===== 用户下拉菜单 =====
    const userMenu = {
        init() {
            // 点击用户按钮切换菜单
            elements.userMenuBtn?.addEventListener('click', (e) => {
                e.stopPropagation();
                this.toggle();
            });

            // 点击页面其他地方关闭菜单
            document.addEventListener('click', () => {
                this.close();
            });

            // 阻止菜单内部点击事件冒泡
            elements.userMenu?.addEventListener('click', (e) => {
                e.stopPropagation();
            });
        },

        toggle() {
            elements.userDropdown?.classList.toggle('active');
        },

        close() {
            elements.userDropdown?.classList.remove('active');
        }
    };

    // ===== 侧边栏菜单 =====
    const sidebar = {
        init() {
            // 初始化：隐藏所有子菜单
            document.querySelectorAll('.submenu').forEach(submenu => {
                submenu.classList.remove('active');
            });

            // 父菜单点击事件
            elements.menuParents.forEach(parent => {
                parent.addEventListener('click', (e) => {
                    e.preventDefault();
                    this.toggleSubmenu(parent);
                });
            });

            // 单页菜单点击事件
            document.querySelectorAll('.menu-item-single').forEach(item => {
                item.addEventListener('click', (e) => {
                    e.preventDefault();
                    const page = item.getAttribute('data-page');
                    if (page) {
                        this.setActiveMenu(item);
                        pageLoader.loadPage(page);
                    }
                });
            });

            // 子菜单点击事件
            elements.submenuItems.forEach(item => {
                item.addEventListener('click', (e) => {
                    e.preventDefault();
                    const page = item.getAttribute('data-page');
                    if (page) {
                        this.setActiveSubmenu(item);
                        pageLoader.loadPage(page);
                    }
                });
            });
        },

        toggleSubmenu(parent) {
            const menuId = parent.getAttribute('data-menu');
            const submenu = document.getElementById(`menu-${menuId}`);

            if (!submenu) return;

            // 切换当前菜单
            const isActive = submenu.classList.contains('active');

            // 关闭所有子菜单
            document.querySelectorAll('.submenu').forEach(menu => {
                menu.classList.remove('active');
            });
            document.querySelectorAll('.menu-item-parent').forEach(item => {
                item.classList.remove('active');
            });

            // 如果之前是关闭的，现在打开
            if (!isActive) {
                submenu.classList.add('active');
                parent.classList.add('active');
            }
        },

        setActiveMenu(menuItem) {
            // 移除所有菜单的激活状态
            document.querySelectorAll('.menu-item').forEach(item => {
                item.classList.remove('active');
            });
            document.querySelectorAll('.submenu-item').forEach(item => {
                item.classList.remove('active');
            });

            // 设置当前菜单为激活
            menuItem.classList.add('active');
        },

        setActiveSubmenu(submenuItem) {
            // 移除所有子菜单的激活状态
            document.querySelectorAll('.submenu-item').forEach(item => {
                item.classList.remove('active');
            });

            // 设置当前子菜单为激活
            submenuItem.classList.add('active');
        }
    };

    // ===== 页面加载器 =====
    const pageLoader = {
        async loadPage(url) {
            if (state.currentPage === url) return;

            try {
                // 显示加载状态
                this.showLoading();

                // 使用 Fetch API 加载页面
                const response = await fetch(`/${url}`);

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const html = await response.text();

                // 更新内容
                elements.mainContent.innerHTML = html;
                elements.mainContent.classList.add('fade-in');

                // 更新状态
                state.currentPage = url;

                // 执行页面中的脚本（如果有）
                this.executeScripts();

            } catch (error) {
                console.error('页面加载失败:', error);
                this.showError('页面加载失败，请稍后重试。');
            }
        },

        showLoading() {
            elements.mainContent.innerHTML = `
                <div class="loading-page">
                    <div class="loading-spinner"></div>
                    <div class="loading-text">加载中...</div>
                </div>
            `;
        },

        showError(message) {
            elements.mainContent.innerHTML = `
                <div class="welcome-screen">
                    <h2 style="color: var(--danger);">加载失败</h2>
                    <p>${message}</p>
                </div>
            `;
        },

        executeScripts() {
            // 执行新加载内容中的脚本
            const scripts = elements.mainContent.querySelectorAll('script');
            scripts.forEach(oldScript => {
                const newScript = document.createElement('script');
                Array.from(oldScript.attributes).forEach(attr => {
                    newScript.setAttribute(attr.name, attr.value);
                });
                newScript.appendChild(document.createTextNode(oldScript.innerHTML));
                oldScript.parentNode.replaceChild(newScript, oldScript);
            });
        }
    };

    // ===== 密码更改 =====
    const passwordChange = {
        init() {
            // 打开密码模态框
            elements.changePasswordBtn?.addEventListener('click', (e) => {
                e.preventDefault();
                this.openModal();
                userMenu.close();
            });

            // 关闭模态框
            elements.closePasswordModal?.addEventListener('click', () => {
                this.closeModal();
            });

            elements.cancelPasswordBtn?.addEventListener('click', () => {
                this.closeModal();
            });

            // 点击模态框外部关闭
            elements.passwordModal?.addEventListener('click', (e) => {
                if (e.target === elements.passwordModal) {
                    this.closeModal();
                }
            });

            // 保存密码
            elements.savePasswordBtn?.addEventListener('click', () => {
                this.savePassword();
            });

            // 表单提交
            elements.passwordForm?.addEventListener('submit', (e) => {
                e.preventDefault();
                this.savePassword();
            });
        },

        openModal() {
            elements.passwordModal?.classList.add('active');
            setTimeout(() => {
                elements.oldPassword?.focus();
            }, 100);
        },

        closeModal() {
            elements.passwordModal?.classList.remove('active');
            elements.passwordForm?.reset();
            utils.clearErrors();
        },

        async savePassword() {
            // 清除之前的错误
            utils.clearErrors();

            // 获取表单数据
            const oldPassword = elements.oldPassword?.value || '';
            const newPassword = elements.newPassword?.value || '';
            const confirmPassword = elements.confirmPassword?.value || '';

            // 验证
            if (!oldPassword) {
                utils.setError(elements.oldPasswordError, '旧密码不能为空!');
                elements.oldPassword?.focus();
                return;
            }

            if (!newPassword) {
                utils.setError(elements.newPasswordError, '新密码不能为空！');
                elements.newPassword?.focus();
                return;
            }

            if (!confirmPassword) {
                utils.setError(elements.confirmPasswordError, '请再次输入新密码!');
                elements.confirmPassword?.focus();
                return;
            }

            if (newPassword !== confirmPassword) {
                utils.setError(elements.confirmPasswordError, '新密码两次输入不一致!');
                elements.confirmPassword?.focus();
                return;
            }

            if (oldPassword === newPassword) {
                utils.setError(elements.newPasswordError, '新密码不能和旧密码一样!');
                elements.newPassword?.focus();
                return;
            }

            const passwordError = utils.validatePassword(newPassword);
            if (passwordError) {
                utils.setError(elements.newPasswordError, passwordError);
                elements.newPassword?.focus();
                return;
            }

            // 提交到服务器
            try {
                elements.savePasswordBtn.disabled = true;
                elements.savePasswordBtn.textContent = '保存中...';

                const response = await fetch('/update/password', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    },
                    body: new URLSearchParams({
                        oldPassword,
                        newPassword,
                        twicePassword: confirmPassword
                    })
                });

                const result = await response.json();

                if (result.Code === 200) {
                    // 成功
                    this.closeModal();
                    alert('密码修改成功!');
                    setTimeout(() => {
                        window.location.href = '/login.html';
                    }, 1000);
                } else if (result.Code === 404) {
                    // 登录过期
                    window.location.href = '/login.html';
                } else {
                    // 失败
                    const errorElement = document.getElementById(result.Key + 'Error');
                    if (errorElement) {
                        utils.setError(errorElement, result.Msg);
                    } else {
                        alert(result.Msg || '密码修改失败！');
                    }
                }
            } catch (error) {
                console.error('密码修改失败:', error);
                alert('系统异常，请稍后再试!');
            } finally {
                elements.savePasswordBtn.disabled = false;
                elements.savePasswordBtn.textContent = '保存';
            }
        }
    };

    // ===== 退出登录 =====
    const logout = {
        init() {
            elements.logoutBtn?.addEventListener('click', async (e) => {
                e.preventDefault();

                if (!confirm('确定要退出登录吗？')) {
                    return;
                }

                try {
                    const response = await fetch('/logout', {
                        method: 'GET'
                    });

                    const result = await response.json();

                    if (result.Code === 200) {
                        window.location.href = '/login.html';
                    } else {
                        alert('系统异常，退出失败!');
                    }
                } catch (error) {
                    console.error('退出登录失败:', error);
                    alert('系统异常!');
                }

                userMenu.close();
            });
        }
    };

    // ===== 初始化 =====
    function init() {
        userMenu.init();
        sidebar.init();
        passwordChange.init();
        logout.init();

        // 默认加载控制面板
        setTimeout(() => {
            const mainMenuItem = document.querySelector('[data-page="main.html"]');
            if (mainMenuItem) {
                mainMenuItem.click();
            }
        }, 100);

        console.log('Admin layout initialized');
    }

    // 页面加载完成后初始化
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // 暴露部分 API 供外部使用
    window.AdminLayout = {
        loadPage: pageLoader.loadPage.bind(pageLoader),
        closeUserMenu: userMenu.close.bind(userMenu)
    };

})();
