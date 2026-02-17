import React, { useState, useEffect } from 'react';
import { Plus, User, Pencil, Trash2, Search, Loader2, Mail, Phone, AlertTriangle, FileSpreadsheet, Check, X, ChevronLeft, ChevronRight } from 'lucide-react';
import { employeeApi, departmentApi, titleApi } from '../services/api';
import type { Employee, EmployeeFormData, Department, Title } from '../types';

const emptyForm: EmployeeFormData = {
    FirstName: '',
    LastName: '',
    TitleID: 0,
    DepartmentID: 0,
    Email: '',
    Phone: '',
    HourlyRate: 50,
    IsShiftWorker: true,
};

const ITEMS_PER_PAGE = 10;

const EmployeeManager: React.FC = () => {
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [titles, setTitles] = useState<Title[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [formData, setFormData] = useState<EmployeeFormData>({ ...emptyForm });
    const [search, setSearch] = useState('');
    const [filterDept, setFilterDept] = useState<number>(0);
    const [confirmDelete, setConfirmDelete] = useState<number | null>(null);
    const [saving, setSaving] = useState(false);
    const [importing, setImporting] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);

    const fileInputRef = React.useRef<HTMLInputElement>(null);

    const fetchEmployees = async () => {
        setLoading(true);
        try {
            const [empRes, deptRes, titleRes] = await Promise.all([
                employeeApi.list(),
                departmentApi.list(),
                titleApi.list(),
            ]);
            setEmployees(empRes.data);
            setDepartments(deptRes.data);
            setTitles(titleRes.data);
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

    const handleImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || e.target.files.length === 0) return;
        const file = e.target.files[0];
        setImporting(true);
        try {
            await employeeApi.import(file);
            alert('Personel listesi başarıyla içeri aktarıldı.');
            fetchEmployees();
        } catch (err) {
            alert('İçe aktarma hatası: ' + err);
        } finally {
            setImporting(false);
            if (fileInputRef.current) fileInputRef.current.value = '';
        }
    };

    const handleEdit = (emp: Employee) => {
        setEditingId(emp.ID);
        setFormData({
            FirstName: emp.FirstName,
            LastName: emp.LastName,
            TitleID: emp.TitleID,
            DepartmentID: emp.DepartmentID,
            Email: emp.Email,
            Phone: emp.Phone,
            HourlyRate: emp.HourlyRate,
            IsShiftWorker: emp.IsShiftWorker !== undefined ? emp.IsShiftWorker : true,
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

    const getDeptName = (deptId: number) => {
        const dept = departments.find(d => d.ID === deptId);
        return dept ? `${dept.Floor}. Kat - ${dept.Name}` : '-';
    };

    const getTitleName = (titleId: number) => {
        const title = titles.find(t => t.ID === titleId);
        return title ? title.Name : '-';
    };

    const filtered = employees.filter((emp) => {
        const q = search.toLowerCase();
        const matchSearch =
            emp.FirstName.toLowerCase().includes(q) ||
            emp.LastName.toLowerCase().includes(q) ||
            getDeptName(emp.DepartmentID).toLowerCase().includes(q) ||
            getTitleName(emp.TitleID).toLowerCase().includes(q);
        const matchDept = filterDept === 0 || emp.DepartmentID === filterDept;
        return matchSearch && matchDept;
    });

    const totalPages = Math.ceil(filtered.length / ITEMS_PER_PAGE);
    const paginatedEmployees = filtered.slice((currentPage - 1) * ITEMS_PER_PAGE, currentPage * ITEMS_PER_PAGE);

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">
                    Personel Listesi
                    <span className="text-sm font-normal text-gray-500 ml-2">({employees.length})</span>
                </h2>
                <div className="flex items-center gap-3">
                    <select
                        className="glass-input py-2 px-3 text-sm"
                        value={filterDept}
                        onChange={(e) => { setFilterDept(Number(e.target.value)); setCurrentPage(1); }}
                    >
                        <option value={0}>Tüm Bölümler</option>
                        {departments.map((d) => (
                            <option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>
                        ))}
                    </select>
                    <div className="relative">
                        <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
                        <input
                            placeholder="Ara..."
                            className="glass-input pl-9 pr-3 py-2 w-52"
                            value={search}
                            onChange={(e) => { setSearch(e.target.value); setCurrentPage(1); }}
                        />
                    </div>
                    <input
                        type="file"
                        ref={fileInputRef}
                        className="hidden"
                        accept=".xlsx,.xls"
                        onChange={handleImport}
                    />
                    <button
                        onClick={() => fileInputRef.current?.click()}
                        disabled={importing}
                        className="btn-ghost flex items-center gap-2"
                    >
                        {importing ? <Loader2 className="w-4 h-4 animate-spin" /> : <FileSpreadsheet className="w-4 h-4" />}
                        Excel İçe Aktar
                    </button>
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
                            <select
                                className="glass-input w-full"
                                value={formData.TitleID}
                                onChange={(e) => setFormData({ ...formData, TitleID: Number(e.target.value) })}
                            >
                                <option value={0}>Ünvan Seçin (opsiyonel)</option>
                                {titles.map((t) => (
                                    <option key={t.ID} value={t.ID}>{t.Name}</option>
                                ))}
                            </select>
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Bölüm *</label>
                            <select
                                className="glass-input w-full"
                                value={formData.DepartmentID}
                                onChange={(e) => setFormData({ ...formData, DepartmentID: Number(e.target.value) })}
                                required
                            >
                                <option value={0} disabled>Bölüm Seçin</option>
                                {departments.map((d) => (
                                    <option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>
                                ))}
                            </select>
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
                        <div className="space-y-1 flex items-end pb-2">
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    className="w-4 h-4 rounded border-gray-600 bg-white/5 text-blue-500 focus:ring-blue-500 focus:ring-offset-0"
                                    checked={formData.IsShiftWorker}
                                    onChange={(e) => setFormData({ ...formData, IsShiftWorker: e.target.checked })}
                                />
                                <span className="text-sm font-medium">Nöbet Tutar / Vardiyalı</span>
                            </label>
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

            {/* Employee Table */}
            {loading ? (
                <div className="space-y-2">
                    {[1, 2, 3, 4, 5].map((i) => (
                        <div key={i} className="skeleton h-12 rounded-lg" />
                    ))}
                </div>
            ) : filtered.length === 0 ? (
                <div className="text-center py-16 text-gray-500">
                    <User className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>{search || filterDept ? 'Arama sonucu bulunamadı.' : 'Henüz personel eklenmemiş.'}</p>
                </div>
            ) : (
                <div className="glass-card overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm text-left">
                            <thead className="bg-white/5 text-gray-400 uppercase text-xs font-semibold">
                                <tr>
                                    <th className="px-6 py-4">Ad Soyad</th>
                                    <th className="px-6 py-4">Bölüm / Ünvan</th>
                                    <th className="px-6 py-4">İletişim</th>
                                    <th className="px-6 py-4">Çalışma Tipi</th>
                                    <th className="px-6 py-4 text-right">Saatlik Ücret</th>
                                    <th className="px-6 py-4 text-right">İşlemler</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                {paginatedEmployees.map((emp) => (
                                    <tr key={emp.ID} className="hover:bg-white/5 transition-colors group">
                                        <td className="px-6 py-4 font-medium">
                                            <div className="flex items-center gap-3">
                                                <div className="w-9 h-9 rounded-full bg-gradient-to-br from-blue-500/20 to-purple-500/20 flex items-center justify-center text-blue-400">
                                                    <User className="w-4 h-4" />
                                                </div>
                                                {emp.FirstName} {emp.LastName}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-gray-300">
                                            <div className="flex flex-col">
                                                <span>{getDeptName(emp.DepartmentID)}</span>
                                                <span className="text-xs text-gray-500">{getTitleName(emp.TitleID)}</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-gray-400">
                                            <div className="space-y-1">
                                                {emp.Email && <div className="flex items-center gap-1.5"><Mail className="w-3 h-3" /> {emp.Email}</div>}
                                                {emp.Phone && <div className="flex items-center gap-1.5"><Phone className="w-3 h-3" /> {emp.Phone}</div>}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            {emp.IsShiftWorker ? (
                                                <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-500/10 text-green-400 border border-green-500/20">
                                                    <Check className="w-3 h-3" /> Nöbetçi
                                                </span>
                                            ) : (
                                                <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-500/10 text-gray-400 border border-gray-500/20">
                                                    <X className="w-3 h-3" /> Sabit
                                                </span>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono text-blue-400">
                                            ₺{emp.HourlyRate}
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                                <button
                                                    onClick={() => handleEdit(emp)}
                                                    className="p-2 rounded-lg hover:bg-white/10 text-gray-400 hover:text-blue-400"
                                                    title="Düzenle"
                                                >
                                                    <Pencil className="w-4 h-4" />
                                                </button>
                                                <button
                                                    onClick={() => setConfirmDelete(emp.ID)}
                                                    className="p-2 rounded-lg hover:bg-white/10 text-gray-400 hover:text-red-400"
                                                    title="Sil"
                                                >
                                                    <Trash2 className="w-4 h-4" />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>

                    {/* Pagination */}
                    {totalPages > 1 && (
                        <div className="flex items-center justify-between px-6 py-4 border-t border-white/5">
                            <div className="text-sm text-gray-500">
                                Toplam {filtered.length} kayıttan {(currentPage - 1) * ITEMS_PER_PAGE + 1} - {Math.min(currentPage * ITEMS_PER_PAGE, filtered.length)} arası gösteriliyor
                            </div>
                            <div className="flex items-center gap-2">
                                <button
                                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                                    disabled={currentPage === 1}
                                    className="p-2 rounded-lg hover:bg-white/5 disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    <ChevronLeft className="w-4 h-4" />
                                </button>
                                {Array.from({ length: totalPages }, (_, i) => i + 1).map(p => (
                                    <button
                                        key={p}
                                        onClick={() => setCurrentPage(p)}
                                        className={`w-8 h-8 rounded-lg text-sm font-medium transition-colors ${currentPage === p ? 'bg-blue-500 text-white' : 'hover:bg-white/5 text-gray-400'
                                            }`}
                                    >
                                        {p}
                                    </button>
                                ))}
                                <button
                                    onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                                    disabled={currentPage === totalPages}
                                    className="p-2 rounded-lg hover:bg-white/5 disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    <ChevronRight className="w-4 h-4" />
                                </button>
                            </div>
                        </div>
                    )}
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
