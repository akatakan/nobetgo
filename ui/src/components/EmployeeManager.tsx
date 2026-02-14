import React, { useState, useEffect } from 'react';
import { Plus, User } from 'lucide-react';
import { employeeApi } from '../services/api';

const EmployeeManager: React.FC = () => {
    const [employees, setEmployees] = useState<any[]>([]);
    const [showForm, setShowForm] = useState(false);
    const [formData, setFormData] = useState({
        FirstName: '',
        LastName: '',
        Title: '',
        Department: '',
        Email: '',
        Phone: '',
        HourlyRate: 50,
    });

    const fetchEmployees = async () => {
        try {
            const res = await employeeApi.list();
            setEmployees(res.data);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchEmployees();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await employeeApi.create(formData);
            setShowForm(false);
            fetchEmployees();
        } catch (err) {
            alert('Hata oluştu');
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Personel Listesi</h2>
                <button
                    onClick={() => setShowForm(!showForm)}
                    className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded-lg"
                >
                    <Plus className="w-4 h-4" />
                    Yeni Personel
                </button>
            </div>

            {showForm && (
                <form onSubmit={handleSubmit} className="bg-gray-700 p-6 rounded-lg grid grid-cols-2 gap-4">
                    <input
                        placeholder="Ad"
                        className="bg-gray-800 border-gray-600 rounded p-2"
                        value={formData.FirstName}
                        onChange={(e) => setFormData({ ...formData, FirstName: e.target.value })}
                        required
                    />
                    <input
                        placeholder="Soyad"
                        className="bg-gray-800 border-gray-600 rounded p-2"
                        value={formData.LastName}
                        onChange={(e) => setFormData({ ...formData, LastName: e.target.value })}
                        required
                    />
                    <input
                        placeholder="Ünvan (Dr., Hemşire vb.)"
                        className="bg-gray-800 border-gray-600 rounded p-2"
                        value={formData.Title}
                        onChange={(e) => setFormData({ ...formData, Title: e.target.value })}
                    />
                    <input
                        placeholder="Bölüm"
                        className="bg-gray-800 border-gray-600 rounded p-2"
                        value={formData.Department}
                        onChange={(e) => setFormData({ ...formData, Department: e.target.value })}
                    />
                    <input
                        placeholder="Saatlik Ücret"
                        type="number"
                        className="bg-gray-800 border-gray-600 rounded p-2"
                        value={formData.HourlyRate}
                        onChange={(e) => setFormData({ ...formData, HourlyRate: Number(e.target.value) })}
                    />
                    <div className="col-span-2 flex justify-end gap-3 mt-2">
                        <button
                            type="button"
                            onClick={() => setShowForm(false)}
                            className="px-4 py-2 text-gray-400 hover:text-white"
                        >
                            İptal
                        </button>
                        <button type="submit" className="bg-green-600 hover:bg-green-700 px-6 py-2 rounded">
                            Kaydet
                        </button>
                    </div>
                </form>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {employees.map((emp) => (
                    <div key={emp.ID} className="bg-gray-700 p-4 rounded-lg flex items-center gap-4 group">
                        <div className="bg-gray-600 p-3 rounded-full">
                            <User className="w-6 h-6 text-blue-400" />
                        </div>
                        <div className="flex-1">
                            <div className="font-medium text-lg">
                                {emp.FirstName} {emp.LastName}
                            </div>
                            <div className="text-gray-400 text-sm">
                                {emp.Title} - {emp.Department}
                            </div>
                        </div>
                        <div className="text-blue-400 font-bold">${emp.HourlyRate}/h</div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default EmployeeManager;
