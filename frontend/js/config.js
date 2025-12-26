// ============================================
// API CONFIGURATION
// ============================================
const API_BASE_URL = 'https://placement-profiling-system-production-618a.up.railway.app';

// ============================================
// API CLIENT
// ============================================
const api = {
    baseUrl: API_BASE_URL,

    getHeaders() {
        const token = localStorage.getItem('token');
        return {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` })
        };
    },

    async request(method, endpoint, data = null) {
        const options = {
            method,
            headers: this.getHeaders()
        };

        if (data && method !== 'GET') {
            options.body = JSON.stringify(data);
        }

        const response = await fetch(`${this.baseUrl}${endpoint}`, options);

        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = 'index.html';
            throw new Error('Unauthorized');
        }

        const result = await response.json();
        
        if (!response.ok) {
            throw new Error(result.error || 'Request failed');
        }

        return result;
    },

    get: function(endpoint) { return this.request('GET', endpoint); },
    post: function(endpoint, data) { return this.request('POST', endpoint, data); },
    put: function(endpoint, data) { return this.request('PUT', endpoint, data); },
    patch: function(endpoint, data) { return this.request('PATCH', endpoint, data); },
    delete: function(endpoint) { return this.request('DELETE', endpoint); },

    async uploadPhoto(endpoint, file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = async () => {
                try {
                    const base64 = reader.result.split(',')[1];
                    const result = await this.post(endpoint, {
                        photo_base64: base64,
                        content_type: file.type
                    });
                    resolve(result);
                } catch (error) {
                    reject(error);
                }
            };
            reader.onerror = () => reject(new Error('Failed to read file'));
            reader.readAsDataURL(file);
        });
    }
};

// ============================================
// UTILITY FUNCTIONS
// ============================================
const utils = {
    getUser() {
        const user = localStorage.getItem('user');
        return user ? JSON.parse(user) : null;
    },

    isLoggedIn() {
        return !!localStorage.getItem('token');
    },

    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = 'index.html';
    },

    requireAuth(role = null) {
        if (!this.isLoggedIn()) {
            window.location.href = 'index.html';
            return false;
        }
        if (role) {
            const user = this.getUser();
            if (!user || user.role !== role) {
                window.location.href = 'index.html';
                return false;
            }
        }
        return true;
    },

    formatDate(dateStr) {
        if (!dateStr) return '-';
        return new Date(dateStr).toLocaleDateString('en-IN', {
            year: 'numeric', month: 'short', day: 'numeric'
        });
    },

    formatCurrency(amount) {
        if (!amount) return '-';
        return `₹${(amount).toFixed(2)} LPA`;
    },

    debounce(func, wait) {
        let timeout;
        return (...args) => {
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(this, args), wait);
        };
    },

    showToast(message, type = 'info') {
        let container = document.getElementById('toast-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'toast-container';
            container.className = 'fixed bottom-6 right-6 z-50 flex flex-col items-end gap-3';
            document.body.appendChild(container);
        }

        const toast = document.createElement('div');
        const colors = {
            success: 'bg-green-500',
            error: 'bg-red-500',
            info: 'bg-blue-500',
            warning: 'bg-yellow-500'
        };
        const icons = {
            success: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>',
            error: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>',
            info: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>',
            warning: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>'
        };

        toast.className = `${colors[type] || colors.info} text-white px-5 py-3 rounded-xl shadow-2xl transform transition-all duration-300 translate-x-[120%] flex items-center gap-3`;
        toast.innerHTML = `${icons[type] || icons.info}<span class="font-medium">${message}</span>`;
        container.appendChild(toast);

        requestAnimationFrame(() => {
            toast.classList.remove('translate-x-[120%]');
        });

        setTimeout(() => {
            toast.classList.add('translate-x-[120%]', 'opacity-0');
            setTimeout(() => toast.remove(), 300);
        }, 4000);
    },

    showLoading(show = true) {
        let overlay = document.getElementById('loading-overlay');
        if (show) {
            if (!overlay) {
                overlay = document.createElement('div');
                overlay.id = 'loading-overlay';
                overlay.className = 'fixed inset-0 bg-white/80 backdrop-blur-sm z-50 flex items-center justify-center';
                overlay.innerHTML = `
                    <div class="text-center">
                        <div class="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
                        <p class="text-gray-600 font-medium">Loading...</p>
                    </div>
                `;
                document.body.appendChild(overlay);
            }
            overlay.classList.remove('hidden');
        } else if (overlay) {
            overlay.classList.add('hidden');
        }
    }
};
