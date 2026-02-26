import React, { useState } from 'react';
import { authApi } from '../services/api';
import { Lock, Check, AlertCircle, Loader2 } from 'lucide-react';

interface ChangePasswordModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const ChangePasswordModal: React.FC<ChangePasswordModalProps> = ({ isOpen, onClose }) => {
    const [oldPassword, setOldPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);

    if (!isOpen) return null;

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (newPassword !== confirmPassword) {
            setError('Yeni şifreler eşleşmiyor.');
            return;
        }

        if (newPassword.length < 6) {
            setError('Yeni şifre en az 6 karakter olmalıdır.');
            return;
        }

        setLoading(true);
        try {
            await authApi.changePassword({ oldPassword, newPassword });
            setSuccess(true);
            setTimeout(() => {
                onClose();
                setOldPassword('');
                setNewPassword('');
                setConfirmPassword('');
                setSuccess(false);
            }, 2000);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Şifre değiştirilemedi.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
            <div className="bg-[#0f172a] border border-white/10 rounded-2xl w-full max-w-md shadow-2xl animate-in zoom-in-95 duration-200">
                <div className="p-6">
                    <div className="flex items-center justify-between mb-6">
                        <div className="flex items-center gap-3">
                            <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
                                <Lock className="w-5 h-5 text-blue-400" />
                            </div>
                            <h2 className="text-xl font-bold text-white">Şifre Değiştir</h2>
                        </div>
                        <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
                            <span className="text-2xl">&times;</span>
                        </button>
                    </div>

                    {success ? (
                        <div className="py-8 text-center animate-in fade-in duration-300">
                            <div className="w-16 h-16 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
                                <Check className="w-8 h-8 text-emerald-400" />
                            </div>
                            <p className="text-emerald-400 font-medium">Şifreniz başarıyla değiştirildi!</p>
                        </div>
                    ) : (
                        <form onSubmit={handleSubmit} className="space-y-4">
                            {error && (
                                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl flex items-center gap-3 text-red-400 text-sm">
                                    <AlertCircle className="w-4 h-4" />
                                    {error}
                                </div>
                            )}

                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-1.5">Mevcut Şifre</label>
                                <input
                                    type="password"
                                    required
                                    value={oldPassword}
                                    onChange={(e) => setOldPassword(e.target.value)}
                                    className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                    placeholder="••••••••"
                                />
                            </div>

                            <div className="h-px bg-white/5 my-2"></div>

                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-1.5">Yeni Şifre</label>
                                <input
                                    type="password"
                                    required
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                    placeholder="••••••••"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-1.5">Yeni Şifre (Tekrar)</label>
                                <input
                                    type="password"
                                    required
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                    placeholder="••••••••"
                                />
                            </div>

                            <div className="flex gap-3 pt-2">
                                <button
                                    type="button"
                                    onClick={onClose}
                                    className="flex-1 px-4 py-2.5 rounded-xl border border-white/10 text-gray-300 hover:bg-white/5 transition-all"
                                >
                                    İptal
                                </button>
                                <button
                                    type="submit"
                                    disabled={loading}
                                    className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-bold px-4 py-2.5 rounded-xl shadow-lg shadow-blue-500/20 transition-all flex items-center justify-center gap-2"
                                >
                                    {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Güncelle'}
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ChangePasswordModal;
