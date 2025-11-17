/**
 * dongfeng-pay 通用组件库
 * 纯 JavaScript 实现，无需任何框架
 * 用于快速构建管理页面
 */

/* ===== 表格组件 ===== */
const DataTable = {
    /**
     * 创建数据表格
     * @param {Object} options - 配置选项
     * @returns {HTMLElement} 表格元素
     */
    create: function(options) {
        const {
            container,
            columns,
            data = [],
            pagination = true,
            pageSize = 10,
            searchable = true,
            sortable = true,
            actions = []
        } = options;

        const wrapper = document.createElement('div');
        wrapper.className = 'datatable-wrapper';

        // 搜索栏
        if (searchable) {
            const searchBar = this.createSearchBar();
            wrapper.appendChild(searchBar);
        }

        // 表格
        const table = this.createTable(columns, data, sortable, actions);
        wrapper.appendChild(table);

        // 分页
        if (pagination) {
            const pager = this.createPagination(data.length, pageSize);
            wrapper.appendChild(pager);
        }

        if (container) {
            const containerEl = typeof container === 'string'
                ? document.querySelector(container)
                : container;
            containerEl.appendChild(wrapper);
        }

        return wrapper;
    },

    createSearchBar: function() {
        const searchBar = document.createElement('div');
        searchBar.className = 'datatable-search';
        searchBar.innerHTML = `
            <input type="text"
                   class="search-input"
                   placeholder="搜索...">
            <button class="btn btn-primary">搜索</button>
        `;
        return searchBar;
    },

    createTable: function(columns, data, sortable, actions) {
        const table = document.createElement('table');
        table.className = 'datatable';

        // 表头
        const thead = document.createElement('thead');
        const headerRow = document.createElement('tr');

        columns.forEach(col => {
            const th = document.createElement('th');
            th.textContent = col.label;
            if (sortable && col.sortable !== false) {
                th.className = 'sortable';
                th.onclick = () => this.sort(col.key);
            }
            headerRow.appendChild(th);
        });

        if (actions.length > 0) {
            const actionTh = document.createElement('th');
            actionTh.textContent = '操作';
            headerRow.appendChild(actionTh);
        }

        thead.appendChild(headerRow);
        table.appendChild(thead);

        // 表体
        const tbody = document.createElement('tbody');
        data.forEach(row => {
            const tr = document.createElement('tr');

            columns.forEach(col => {
                const td = document.createElement('td');
                const value = this.getValue(row, col.key);
                td.innerHTML = col.render
                    ? col.render(value, row)
                    : value;
                tr.appendChild(td);
            });

            if (actions.length > 0) {
                const actionTd = document.createElement('td');
                actionTd.className = 'action-cell';
                actions.forEach(action => {
                    const btn = document.createElement('button');
                    btn.className = `btn btn-sm btn-${action.type || 'default'}`;
                    btn.textContent = action.label;
                    btn.onclick = () => action.onClick(row);
                    actionTd.appendChild(btn);
                });
                tr.appendChild(actionTd);
            }

            tbody.appendChild(tr);
        });
        table.appendChild(tbody);

        return table;
    },

    createPagination: function(total, pageSize) {
        const totalPages = Math.ceil(total / pageSize);
        const pagination = document.createElement('div');
        pagination.className = 'datatable-pagination';

        pagination.innerHTML = `
            <button class="btn btn-sm" id="prevPage">上一页</button>
            <span class="page-info">第 <span id="currentPage">1</span> / ${totalPages} 页</span>
            <button class="btn btn-sm" id="nextPage">下一页</button>
            <span class="total-info">共 ${total} 条</span>
        `;

        return pagination;
    },

    getValue: function(obj, path) {
        return path.split('.').reduce((curr, key) => curr?.[key], obj) ?? '';
    },

    sort: function(key) {
        console.log('Sorting by:', key);
        // 实现排序逻辑
    }
};

/* ===== 模态框组件 ===== */
const Modal = {
    /**
     * 创建模态框
     * @param {Object} options - 配置选项
     */
    create: function(options) {
        const {
            title,
            content,
            footer,
            size = 'medium',
            onClose
        } = options;

        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.innerHTML = `
            <div class="modal modal-${size}">
                <div class="modal-header">
                    <h3 class="modal-title">${title}</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    ${typeof content === 'string' ? content : ''}
                </div>
                ${footer ? `<div class="modal-footer">${footer}</div>` : ''}
            </div>
        `;

        // 关闭事件
        const closeBtn = modal.querySelector('.modal-close');
        closeBtn.onclick = () => this.close(modal, onClose);

        modal.onclick = (e) => {
            if (e.target === modal) {
                this.close(modal, onClose);
            }
        };

        if (typeof content !== 'string') {
            const body = modal.querySelector('.modal-body');
            body.appendChild(content);
        }

        document.body.appendChild(modal);

        // 动画
        setTimeout(() => modal.classList.add('active'), 10);

        return modal;
    },

    close: function(modal, callback) {
        modal.classList.remove('active');
        setTimeout(() => {
            document.body.removeChild(modal);
            if (callback) callback();
        }, 300);
    },

    confirm: function(message, callback) {
        const footer = `
            <button class="btn btn-default" id="modalCancel">取消</button>
            <button class="btn btn-primary" id="modalConfirm">确定</button>
        `;

        const modal = this.create({
            title: '确认',
            content: `<p>${message}</p>`,
            footer: footer
        });

        modal.querySelector('#modalCancel').onclick = () => {
            this.close(modal);
        };

        modal.querySelector('#modalConfirm').onclick = () => {
            this.close(modal, callback);
        };
    },

    alert: function(message, type = 'info') {
        const icons = {
            success: '✓',
            error: '✗',
            warning: '⚠',
            info: 'ℹ'
        };

        const modal = this.create({
            title: icons[type] + ' ' + (type === 'error' ? '错误' : '提示'),
            content: `<p>${message}</p>`,
            footer: '<button class="btn btn-primary" id="modalOk">确定</button>',
            size: 'small'
        });

        modal.querySelector('#modalOk').onclick = () => {
            this.close(modal);
        };
    }
};

/* ===== 表单组件 ===== */
const Form = {
    /**
     * 创建表单
     * @param {Object} options - 配置选项
     */
    create: function(options) {
        const {
            fields,
            data = {},
            onSubmit,
            labelWidth = '120px'
        } = options;

        const form = document.createElement('form');
        form.className = 'modern-form';
        form.style.setProperty('--label-width', labelWidth);

        fields.forEach(field => {
            const formGroup = this.createFormGroup(field, data[field.name]);
            form.appendChild(formGroup);
        });

        if (onSubmit) {
            form.onsubmit = (e) => {
                e.preventDefault();
                const formData = this.getFormData(form);
                onSubmit(formData);
            };
        }

        return form;
    },

    createFormGroup: function(field, value = '') {
        const group = document.createElement('div');
        group.className = 'form-group';

        const label = document.createElement('label');
        label.className = 'form-label';
        label.textContent = field.label;
        if (field.required) {
            label.innerHTML += '<span class="required">*</span>';
        }

        const inputWrapper = document.createElement('div');
        inputWrapper.className = 'form-input-wrapper';

        let input;
        switch (field.type) {
            case 'textarea':
                input = document.createElement('textarea');
                input.className = 'form-textarea';
                input.rows = field.rows || 4;
                break;
            case 'select':
                input = document.createElement('select');
                input.className = 'form-select';
                field.options?.forEach(opt => {
                    const option = document.createElement('option');
                    option.value = opt.value;
                    option.textContent = opt.label;
                    input.appendChild(option);
                });
                break;
            case 'checkbox':
                input = document.createElement('input');
                input.type = 'checkbox';
                input.className = 'form-checkbox';
                input.checked = value;
                break;
            case 'radio':
                input = document.createElement('div');
                input.className = 'form-radio-group';
                field.options?.forEach(opt => {
                    const radioLabel = document.createElement('label');
                    radioLabel.className = 'radio-label';
                    const radio = document.createElement('input');
                    radio.type = 'radio';
                    radio.name = field.name;
                    radio.value = opt.value;
                    radioLabel.appendChild(radio);
                    radioLabel.appendChild(document.createTextNode(opt.label));
                    input.appendChild(radioLabel);
                });
                break;
            default:
                input = document.createElement('input');
                input.type = field.type || 'text';
                input.className = 'form-input';
        }

        if (input.tagName !== 'DIV') {
            input.name = field.name;
            input.value = value;
            input.placeholder = field.placeholder || '';
            if (field.required) input.required = true;
            if (field.disabled) input.disabled = true;
        }

        inputWrapper.appendChild(input);

        if (field.help) {
            const help = document.createElement('small');
            help.className = 'form-help';
            help.textContent = field.help;
            inputWrapper.appendChild(help);
        }

        group.appendChild(label);
        group.appendChild(inputWrapper);

        return group;
    },

    getFormData: function(form) {
        const formData = new FormData(form);
        const data = {};
        for (let [key, value] of formData.entries()) {
            data[key] = value;
        }
        return data;
    },

    validate: function(form) {
        return form.checkValidity();
    }
};

/* ===== 分页组件 ===== */
const Pagination = {
    create: function(options) {
        const {
            total,
            pageSize = 10,
            currentPage = 1,
            onChange
        } = options;

        const totalPages = Math.ceil(total / pageSize);

        const wrapper = document.createElement('div');
        wrapper.className = 'pagination';

        // 上一页
        const prev = this.createButton('上一页', currentPage === 1, () => {
            if (currentPage > 1 && onChange) {
                onChange(currentPage - 1);
            }
        });
        wrapper.appendChild(prev);

        // 页码
        const pages = this.getPageNumbers(currentPage, totalPages);
        pages.forEach(page => {
            if (page === '...') {
                const ellipsis = document.createElement('span');
                ellipsis.className = 'pagination-ellipsis';
                ellipsis.textContent = '...';
                wrapper.appendChild(ellipsis);
            } else {
                const btn = this.createButton(
                    page.toString(),
                    false,
                    () => onChange && onChange(page)
                );
                if (page === currentPage) {
                    btn.classList.add('active');
                }
                wrapper.appendChild(btn);
            }
        });

        // 下一页
        const next = this.createButton('下一页', currentPage === totalPages, () => {
            if (currentPage < totalPages && onChange) {
                onChange(currentPage + 1);
            }
        });
        wrapper.appendChild(next);

        // 信息
        const info = document.createElement('span');
        info.className = 'pagination-info';
        info.textContent = `共 ${total} 条，第 ${currentPage}/${totalPages} 页`;
        wrapper.appendChild(info);

        return wrapper;
    },

    createButton: function(text, disabled, onClick) {
        const btn = document.createElement('button');
        btn.className = 'pagination-btn';
        btn.textContent = text;
        btn.disabled = disabled;
        btn.onclick = onClick;
        return btn;
    },

    getPageNumbers: function(current, total) {
        if (total <= 7) {
            return Array.from({ length: total }, (_, i) => i + 1);
        }

        if (current <= 3) {
            return [1, 2, 3, 4, 5, '...', total];
        }

        if (current >= total - 2) {
            return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
        }

        return [1, '...', current - 1, current, current + 1, '...', total];
    }
};

/* ===== 卡片组件 ===== */
const Card = {
    create: function(options) {
        const {
            title,
            content,
            footer,
            className = ''
        } = options;

        const card = document.createElement('div');
        card.className = `card ${className}`;

        if (title) {
            const header = document.createElement('div');
            header.className = 'card-header';
            header.innerHTML = `<h3 class="card-title">${title}</h3>`;
            card.appendChild(header);
        }

        const body = document.createElement('div');
        body.className = 'card-body';
        if (typeof content === 'string') {
            body.innerHTML = content;
        } else {
            body.appendChild(content);
        }
        card.appendChild(body);

        if (footer) {
            const footerEl = document.createElement('div');
            footerEl.className = 'card-footer';
            footerEl.innerHTML = footer;
            card.appendChild(footerEl);
        }

        return card;
    }
};

/* ===== 加载指示器 ===== */
const Loading = {
    show: function(container) {
        const loading = document.createElement('div');
        loading.className = 'loading-overlay';
        loading.innerHTML = `
            <div class="loading-spinner"></div>
            <div class="loading-text">加载中...</div>
        `;

        if (container) {
            const containerEl = typeof container === 'string'
                ? document.querySelector(container)
                : container;
            containerEl.style.position = 'relative';
            containerEl.appendChild(loading);
        } else {
            document.body.appendChild(loading);
        }

        return loading;
    },

    hide: function(loading) {
        if (loading && loading.parentNode) {
            loading.parentNode.removeChild(loading);
        }
    }
};

/* ===== 通知组件 ===== */
const Notification = {
    show: function(message, type = 'info', duration = 3000) {
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;

        const icons = {
            success: '✓',
            error: '✗',
            warning: '⚠',
            info: 'ℹ'
        };

        notification.innerHTML = `
            <span class="notification-icon">${icons[type]}</span>
            <span class="notification-message">${message}</span>
            <button class="notification-close">&times;</button>
        `;

        document.body.appendChild(notification);

        setTimeout(() => notification.classList.add('show'), 10);

        const close = () => {
            notification.classList.remove('show');
            setTimeout(() => document.body.removeChild(notification), 300);
        };

        notification.querySelector('.notification-close').onclick = close;

        if (duration > 0) {
            setTimeout(close, duration);
        }

        return notification;
    }
};

/* ===== Tab 组件 ===== */
const Tabs = {
    create: function(options) {
        const {
            tabs,
            defaultActive = 0,
            onChange
        } = options;

        const wrapper = document.createElement('div');
        wrapper.className = 'tabs';

        // Tab 头部
        const header = document.createElement('div');
        header.className = 'tabs-header';

        tabs.forEach((tab, index) => {
            const tabBtn = document.createElement('button');
            tabBtn.className = 'tab-btn';
            if (index === defaultActive) tabBtn.classList.add('active');
            tabBtn.textContent = tab.label;
            tabBtn.onclick = () => {
                this.switchTab(wrapper, index);
                if (onChange) onChange(index);
            };
            header.appendChild(tabBtn);
        });

        wrapper.appendChild(header);

        // Tab 内容
        const content = document.createElement('div');
        content.className = 'tabs-content';

        tabs.forEach((tab, index) => {
            const pane = document.createElement('div');
            pane.className = 'tab-pane';
            if (index === defaultActive) pane.classList.add('active');

            if (typeof tab.content === 'string') {
                pane.innerHTML = tab.content;
            } else {
                pane.appendChild(tab.content);
            }

            content.appendChild(pane);
        });

        wrapper.appendChild(content);

        return wrapper;
    },

    switchTab: function(wrapper, index) {
        const buttons = wrapper.querySelectorAll('.tab-btn');
        const panes = wrapper.querySelectorAll('.tab-pane');

        buttons.forEach((btn, i) => {
            btn.classList.toggle('active', i === index);
        });

        panes.forEach((pane, i) => {
            pane.classList.toggle('active', i === index);
        });
    }
};

// 导出组件
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        DataTable,
        Modal,
        Form,
        Pagination,
        Card,
        Loading,
        Notification,
        Tabs
    };
}
