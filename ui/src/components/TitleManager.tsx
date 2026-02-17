import React, { useState, useEffect } from 'react';
import { Plus, Award, Pencil, Trash2, Loader2, AlertTriangle } from 'lucide-react';
import { titleApi } from '../services/api';
import type { Title, TitleFormData } from '../types';

const emptyForm: TitleFormData = {
    Name: '',
};

const TitleManager: React.FC = () => {
    const [titles, setTitles] = useState<Title[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [formData, setFormData] = useState<TitleFormData>({ ...emptyForm });
    const [confirmDelete, setConfirmDelete] = useState<number | null>(null);
    const [saving, setSaving] = useState(false);

    const fetchTitles = async () => {
        setLoading(true);
        try {
            const res = await titleApi.list();
            setTitles(res.data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTitles();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (editingId) {
                await titleApi.update(editingId, formData);
            } else {
                await titleApi.create(formData);
            }
            setShowForm(false);
            setEditingId(null);
            setFormData({ ...emptyForm });
            fetchTitles();
        } catch (err) {
            alert('Hata oluştu');
        } finally {
            setSaving(false);
        }
    };

    const handleEdit = (title: Title) => {
        setEditingId(title.ID);
        setFormData({ Name: title.Name });
        setShowForm(true);
    };

    const handleDelete = async (id: number) => {
        try {
            await titleApi.delete(id);
            setConfirmDelete(null);
            fetchTitles();
        } catch (err) {
            alert('Silme hatası: Bu ünvana bağlı personel olabilir.');
        }
    };

    const cancelForm = () => {
        setShowForm(false);
        setEditingId(null);
        setFormData({ ...emptyForm });
    };

    const TITLE_COLORS = [
        'from-violet-500/20 to-purple-500/20',
        'from-blue-500/20 to-indigo-500/20',
        'from-cyan-500/20 to-blue-500/20',
        'from-emerald-500/20 to-cyan-500/20',
        'from-amber-500/20 to-yellow-500/20',
        'from-rose-500/20 to-pink-500/20',
    ];

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">
                    Ünvan Listesi
                    <span className="text-sm font-normal text-gray-500 ml-2">({titles.length})</span>
                </h2>
                <button
                    onClick={() => { cancelForm(); setShowForm(!showForm); }}
                    className="btn-primary"
                >
                    <Plus className="w-4 h-4" />
                    Yeni Ünvan
                </button>
            </div>

            {/* Form */}
            {showForm && (
                <form onSubmit={handleSubmit} className="glass-card p-6 animate-slide-down">
                    <h3 className="text-lg font-semibold mb-4">
                        {editingId ? 'Ünvan Düzenle' : 'Yeni Ünvan Ekle'}
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Ünvan Adı *</label>
                            <input
                                placeholder="Dr., Hemşire, Uzman vb."
                                className="glass-input w-full"
                                value={formData.Name}
                                onChange={(e) => setFormData({ ...formData, Name: e.target.value })}
                                required
                            />
                        </div>
                    </div>
                    <div className="flex justify-end gap-3 mt-6">
                        <button type="button" onClick={cancelForm} className="btn-ghost">
                            İptal
                        </button>
                        <button type="submit" disabled={saving} className="btn-success">
                            {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                            {editingId ? 'Güncelle' : 'Kaydet'}
                        </button>
                    </div>
                </form>
            )}

            {/* Title Cards */}
            {loading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="skeleton h-20 rounded-xl" />
                    ))}
                </div>
            ) : titles.length === 0 ? (
                <div className="text-center py-16 text-gray-500">
                    <Award className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Henüz ünvan eklenmemiş.</p>
                    <p className="text-sm mt-1">Personel eklerken ünvan seçimi yapabilmek için önce ünvan tanımlayın.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {titles.map((title, idx) => (
                        <div
                            key={title.ID}
                            className="glass-card p-5 group hover:border-violet-500/30 transition-all duration-300 cursor-default animate-slide-up"
                            style={{ animationDelay: `${idx * 50}ms`, animationFillMode: 'both' }}
                        >
                            <div className="flex items-center gap-4">
                                <div className={`w-11 h-11 rounded-xl bg-gradient-to-br ${TITLE_COLORS[idx % TITLE_COLORS.length]} flex items-center justify-center flex-shrink-0`}>
                                    <Award className="w-5 h-5 text-violet-400" />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="font-semibold text-base truncate">
                                        {title.Name}
                                    </div>
                                </div>
                                <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
                                    <button
                                        onClick={() => handleEdit(title)}
                                        className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-blue-400 transition-colors"
                                    >
                                        <Pencil className="w-3.5 h-3.5" />
                                    </button>
                                    <button
                                        onClick={() => setConfirmDelete(title.ID)}
                                        className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-red-400 transition-colors"
                                    >
                                        <Trash2 className="w-3.5 h-3.5" />
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Delete Confirmation */}
            {confirmDelete !== null && (
                <div className="dialog-overlay" onClick={() => setConfirmDelete(null)}>
                    <div className="dialog-content" onClick={(e) => e.stopPropagation()}>
                        <div className="flex items-center gap-3 mb-4">
                            <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center">
                                <AlertTriangle className="w-5 h-5 text-red-400" />
                            </div>
                            <div>
                                <h4 className="font-semibold">Ünvanı Sil</h4>
                                <p className="text-sm text-gray-400">Bu işlem geri alınamaz. Bağlı personel varsa silinemez.</p>
                            </div>
                        </div>
                        <div className="flex justify-end gap-3">
                            <button onClick={() => setConfirmDelete(null)} className="btn-ghost">İptal</button>
                            <button onClick={() => handleDelete(confirmDelete)} className="btn-danger">Sil</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default TitleManager;
