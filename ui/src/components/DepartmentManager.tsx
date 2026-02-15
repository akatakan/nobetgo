import React, { useState, useEffect } from 'react';
import { Plus, Building2, Pencil, Trash2, Loader2, AlertTriangle } from 'lucide-react';
import { departmentApi } from '../services/api';
import type { Department, DepartmentFormData } from '../types';

const emptyForm: DepartmentFormData = {
    Name: '',
    Floor: 1,
    Description: '',
};

const DepartmentManager: React.FC = () => {
    const [departments, setDepartments] = useState<Department[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [formData, setFormData] = useState<DepartmentFormData>({ ...emptyForm });
    const [confirmDelete, setConfirmDelete] = useState<number | null>(null);
    const [saving, setSaving] = useState(false);

    const fetchDepartments = async () => {
        setLoading(true);
        try {
            const res = await departmentApi.list();
            setDepartments(res.data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchDepartments();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (editingId) {
                await departmentApi.update(editingId, formData);
            } else {
                await departmentApi.create(formData);
            }
            setShowForm(false);
            setEditingId(null);
            setFormData({ ...emptyForm });
            fetchDepartments();
        } catch (err) {
            alert('Hata oluştu');
        } finally {
            setSaving(false);
        }
    };

    const handleEdit = (dept: Department) => {
        setEditingId(dept.ID);
        setFormData({
            Name: dept.Name,
            Floor: dept.Floor,
            Description: dept.Description,
        });
        setShowForm(true);
    };

    const handleDelete = async (id: number) => {
        try {
            await departmentApi.delete(id);
            setConfirmDelete(null);
            fetchDepartments();
        } catch (err) {
            alert('Silme hatası: Bu bölüme bağlı personel olabilir.');
        }
    };

    const cancelForm = () => {
        setShowForm(false);
        setEditingId(null);
        setFormData({ ...emptyForm });
    };

    const FLOOR_COLORS = [
        'from-blue-500/20 to-cyan-500/20',
        'from-purple-500/20 to-pink-500/20',
        'from-amber-500/20 to-orange-500/20',
        'from-emerald-500/20 to-teal-500/20',
        'from-rose-500/20 to-red-500/20',
    ];

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">
                    Bölüm Listesi
                    <span className="text-sm font-normal text-gray-500 ml-2">({departments.length})</span>
                </h2>
                <button
                    onClick={() => { cancelForm(); setShowForm(!showForm); }}
                    className="btn-primary"
                >
                    <Plus className="w-4 h-4" />
                    Yeni Bölüm
                </button>
            </div>

            {/* Form */}
            {showForm && (
                <form onSubmit={handleSubmit} className="glass-card p-6 animate-slide-down">
                    <h3 className="text-lg font-semibold mb-4">
                        {editingId ? 'Bölüm Düzenle' : 'Yeni Bölüm Ekle'}
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Bölüm Adı *</label>
                            <input
                                placeholder="Dahiliye, Cerrahi vb."
                                className="glass-input w-full"
                                value={formData.Name}
                                onChange={(e) => setFormData({ ...formData, Name: e.target.value })}
                                required
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Kat *</label>
                            <input
                                type="number"
                                min="1"
                                className="glass-input w-full"
                                value={formData.Floor}
                                onChange={(e) => setFormData({ ...formData, Floor: Number(e.target.value) })}
                                required
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Açıklama</label>
                            <input
                                placeholder="İsteğe bağlı açıklama"
                                className="glass-input w-full"
                                value={formData.Description}
                                onChange={(e) => setFormData({ ...formData, Description: e.target.value })}
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

            {/* Department Cards */}
            {loading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="skeleton h-28 rounded-xl" />
                    ))}
                </div>
            ) : departments.length === 0 ? (
                <div className="text-center py-16 text-gray-500">
                    <Building2 className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Henüz bölüm eklenmemiş.</p>
                    <p className="text-sm mt-1">Nöbet oluşturmak için önce bölüm tanımlayın.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {departments.map((dept, idx) => (
                        <div
                            key={dept.ID}
                            className="glass-card p-5 group hover:border-blue-500/30 transition-all duration-300 cursor-default animate-slide-up"
                            style={{ animationDelay: `${idx * 50}ms`, animationFillMode: 'both' }}
                        >
                            <div className="flex items-start gap-4">
                                <div className={`w-11 h-11 rounded-xl bg-gradient-to-br ${FLOOR_COLORS[(dept.Floor - 1) % FLOOR_COLORS.length]} flex items-center justify-center flex-shrink-0`}>
                                    <Building2 className="w-5 h-5 text-blue-400" />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="font-semibold text-base truncate">
                                        {dept.Name}
                                    </div>
                                    <div className="text-gray-400 text-sm">
                                        {dept.Floor}. Kat
                                    </div>
                                    {dept.Description && (
                                        <div className="text-gray-500 text-xs mt-1 truncate">
                                            {dept.Description}
                                        </div>
                                    )}
                                </div>
                                <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
                                    <button
                                        onClick={() => handleEdit(dept)}
                                        className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-blue-400 transition-colors"
                                    >
                                        <Pencil className="w-3.5 h-3.5" />
                                    </button>
                                    <button
                                        onClick={() => setConfirmDelete(dept.ID)}
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
                                <h4 className="font-semibold">Bölümü Sil</h4>
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

export default DepartmentManager;
