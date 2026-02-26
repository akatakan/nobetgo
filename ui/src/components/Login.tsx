import React, { useState } from 'react';
import { authApi } from '../services/api';
import { Lock, Mail, AlertCircle } from 'lucide-react';

interface LoginProps {
    onLoginSuccess: (token: string, role: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLoginSuccess }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const response = await authApi.login({ email, password });
            const { token, role } = response.data;

            localStorage.setItem('token', token);
            localStorage.setItem('role', role);

            onLoginSuccess(token, role);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Giriş yapılamadı. Lütfen bilgilerinizi kontrol edin.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
            <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8 border border-slate-100">
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-600 rounded-2xl mb-4 shadow-lg shadow-blue-200">
                        <Lock className="text-white w-8 h-8" />
                    </div>
                    <h1 className="text-2xl font-bold text-slate-900">NobetGo</h1>
                    <p className="text-slate-500 mt-2">Devam etmek için giriş yapın</p>
                </div>

                {error && (
                    <div className="mb-6 p-4 bg-red-50 border border-red-100 rounded-xl flex items-center gap-3 text-red-700 text-sm">
                        <AlertCircle className="w-5 h-5 flex-shrink-0" />
                        <p>{error}</p>
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                    <div>
                        <label className="block text-sm font-medium text-slate-700 mb-2">E-posta Adresi</label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Mail className="h-5 w-5 text-slate-400" />
                            </div>
                            <input
                                type="email"
                                required
                                className="block w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                placeholder="admin@nobetgo.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-slate-700 mb-2">Şifre</label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Lock className="h-5 w-5 text-slate-400" />
                            </div>
                            <input
                                type="password"
                                required
                                className="block w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                placeholder="••••••••"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold py-3 rounded-xl shadow-lg shadow-blue-200 transition-all active:scale-[0.98]"
                    >
                        {loading ? 'Giriş Yapılıyor...' : 'Giriş Yap'}
                    </button>

                    <div className="text-center">
                        <button
                            type="button"
                            onClick={() => alert('Şifre sıfırlama altyapısı hazırlandı. Yakında aktif edilecektir.')}
                            className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                        >
                            Şifremi Unuttum
                        </button>
                    </div>
                </form>

                <div className="mt-8 pt-8 border-t border-slate-100 text-center">
                    <p className="text-xs text-slate-400">
                        © 2026 NobetGo. Tüm hakları saklıdır.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Login;
