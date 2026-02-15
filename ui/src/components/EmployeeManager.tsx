import React, { useState, useEffect } from 'react';
import { Plus, User, Pencil, Trash2, Search, Loader2, Mail, Phone, AlertTriangle } from 'lucide-react';
import { employeeApi } from '../services/api';
import type { Employee, EmployeeFormData } from '../types';

const emptyForm: EmployeeFormData = {
    FirstName: '',
    LastName: '',
    Title: '',
    Department: '',
    Email: '',
    Phone: '',
    HourlyRate: 50,
};

const EmployeeManager: React.FC = () => {
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [formData, setFormData] = useState<EmployeeFormData>({ ...emptyForm });
    const [search, setSearch] = useState('');
    const [confirmDelete, setConfirmDelete] = useState<number | null>(null);
    const [saving, setSaving] = useState(false);

    const fetchEmployees = async () => {
        setLoading(true);
        try {
            const res = await employeeApi.list();
            setEmployees(res.data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchEmployees();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (editingId) {
                await employeeApi.update(editingId, formData);
            } else {
                await employeeApi.create(formData);
            }
            setShowForm(false);
            setEditingId(null);
            setFormData({ ...emptyForm });
            fetchEmployees();
        } catch (err) {
            alert('Hata oluştu');
        } finally {
            setSaving(false);
        }
    };

    const handleEdit = (emp: Employee) => {
        setEditingId(emp.ID);
        setFormData({
            FirstName: emp.FirstName,
            LastName: emp.LastName,
            Title: emp.Title,
            Department: emp.Department,
            Email: emp.Email,
            Phone: emp.Phone,
            HourlyRate: emp.HourlyRate,
        });
        setShowForm(true);
    };

    const handleDelete = async (id: number) => {
        try {
            await employeeApi.delete(id);
            setConfirmDelete(null);
            fetchEmployees();
        } catch (err) {
            alert('Silme hatası');
        }
    };

    const cancelForm = () => {
        setShowForm(false);
        setEditingId(null);
        setFormData({ ...emptyForm });
    };

    const filtered = employees.filter((emp) => {
        const q = search.toLowerCase();
        return (
            emp.FirstName.toLowerCase().includes(q) ||
            emp.LastName.toLowerCase().includes(q) ||
            emp.Department.toLowerCase().includes(q) ||
            emp.Title.toLowerCase().includes(q)
        );
    });

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">
                    Personel Listesi
                    <span className="text-sm font-normal text-gray-500 ml-2">({employees.length})</span>
                </h2>
                <div className="flex items-center gap-3">
                    <div className="relative">
                        <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
                        <input
                            placeholder="Ara..."
                            className="glass-input pl-9 pr-3 py-2 w-52"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                        />
                    </div>
                    <button
                        onClick={() => { cancelForm(); setShowForm(!showForm); }}
                        className="btn-primary"
                    >
                        <Plus className="w-4 h-4" />
                        Yeni Personel
                    </button>
                </div>
            </div>

            {/* Form */}
            {showForm && (
                <form onSubmit={handleSubmit} className="glass-card p-6 animate-slide-down">
                    <h3 className="text-lg font-semibold mb-4">
                        {editingId ? 'Personel Düzenle' : 'Yeni Personel Ekle'}
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Ad *</label>
                            <input
                                placeholder="Ad"
                                className="glass-input w-full"
                                value={formData.FirstName}
                                onChange={(e) => setFormData({ ...formData, FirstName: e.target.value })}
                                required
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Soyad *</label>
                            <input
                                placeholder="Soyad"
                                className="glass-input w-full"
                                value={formData.LastName}
                                onChange={(e) => setFormData({ ...formData, LastName: e.target.value })}
                                required
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Ünvan</label>
                            <input
                                placeholder="Dr., Hemşire vb."
                                className="glass-input w-full"
                                value={formData.Title}
                                onChange={(e) => setFormData({ ...formData, Title: e.target.value })}
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Bölüm</label>
                            <input
                                placeholder="Bölüm"
                                className="glass-input w-full"
                                value={formData.Department}
                                onChange={(e) => setFormData({ ...formData, Department: e.target.value })}
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">E-posta</label>
                            <input
                                type="email"
                                placeholder="ornek@email.com"
                                className="glass-input w-full"
                                value={formData.Email}
                                onChange={(e) => setFormData({ ...formData, Email: e.target.value })}
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Telefon</label>
                            <input
                                placeholder="+90 5XX XXX XX XX"
                                className="glass-input w-full"
                                value={formData.Phone}
                                onChange={(e) => setFormData({ ...formData, Phone: e.target.value })}
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Saatlik Ücret (₺)</label>
                            <input
                                type="number"
                                className="glass-input w-full"
                                value={formData.HourlyRate}
                                onChange={(e) => setFormData({ ...formData, HourlyRate: Number(e.target.value) })}
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

            {/* Employee Cards */}
            {loading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="skeleton h-28 rounded-xl" />
                    ))}
                </div>
            ) : filtered.length === 0 ? (
                <div className="text-center py-16 text-gray-500">
                    <User className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>{search ? 'Arama sonucu bulunamadı.' : 'Henüz personel eklenmemiş.'}</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {filtered.map((emp, idx) => (
                        <div
                            key={emp.ID}
                            className="glass-card p-5 group hover:border-blue-500/30 transition-all duration-300 cursor-default animate-slide-up"
                            style={{ animationDelay: `${idx * 50}ms`, animationFillMode: 'both' }}
                        >
                            <div className="flex items-start gap-4">
                                <div className="w-11 h-11 rounded-xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 flex items-center justify-center flex-shrink-0">
                                    <User className="w-5 h-5 text-blue-400" />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="font-semibold text-base truncate">
                                        {emp.FirstName} {emp.LastName}
                                    </div>
                                    <div className="text-gray-400 text-sm truncate">
                                        {emp.Title}{emp.Title && emp.Department ? ' · ' : ''}{emp.Department}
                                    </div>
                                    <div className="flex items-center gap-3 mt-2 text-xs text-gray-500">
                                        {emp.Email && (
                                            <span className="flex items-center gap-1 truncate">
                                                <Mail className="w-3 h-3" />{emp.Email}
                                            </span>
                                        )}
                                        {emp.Phone && (
                                            <span className="flex items-center gap-1">
                                                <Phone className="w-3 h-3" />{emp.Phone}
                                            </span>
                                        )}
                                    </div>
                                </div>
                                <div className="text-right flex-shrink-0">
                                    <div className="text-blue-400 font-bold text-sm">₺{emp.HourlyRate}/s</div>
                                    <div className="flex gap-1 mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <button
                                            onClick={() => handleEdit(emp)}
                                            className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-blue-400 transition-colors"
                                        >
                                            <Pencil className="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            onClick={() => setConfirmDelete(emp.ID)}
                                            className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-red-400 transition-colors"
                                        >
                                            <Trash2 className="w-3.5 h-3.5" />
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Delete Confirmation Dialog */}
            {confirmDelete !== null && (
                <div className="dialog-overlay" onClick={() => setConfirmDelete(null)}>
                    <div className="dialog-content" onClick={(e) => e.stopPropagation()}>
                        <div className="flex items-center gap-3 mb-4">
                            <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center">
                                <AlertTriangle className="w-5 h-5 text-red-400" />
                            </div>
                            <div>
                                <h4 className="font-semibold">Personeli Sil</h4>
                                <p className="text-sm text-gray-400">Bu işlem geri alınamaz.</p>
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

export default EmployeeManager;
