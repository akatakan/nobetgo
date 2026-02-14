import React, { useState, useEffect } from 'react';
import { ClipboardCheck, Clock } from 'lucide-react';
import { scheduleApi } from '../services/api';
import { format } from 'date-fns';
import { tr } from 'date-fns/locale';

const ScheduleViewer: React.FC = () => {
    const [data, setData] = useState<any[]>([]);
    const [params, setParams] = useState({ month: new Date().getMonth() + 1, year: new Date().getFullYear() });

    const fetchData = async () => {
        try {
            const res = await scheduleApi.get(params.month, params.year);
            setData(res.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchData();
    }, [params]);

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Mevcut Çizelge</h2>
                <div className="flex gap-4">
                    <select
                        className="bg-gray-700 border-none rounded p-2"
                        value={params.month}
                        onChange={(e) => setParams({ ...params, month: Number(e.target.value) })}
                    >
                        {Array.from({ length: 12 }, (_, i) => (
                            <option key={i + 1} value={i + 1}>{i + 1}. Ay</option>
                        ))}
                    </select>
                    <input
                        type="number"
                        className="bg-gray-700 border-none rounded p-2 w-24"
                        value={params.year}
                        onChange={(e) => setParams({ ...params, year: Number(e.target.value) })}
                    />
                </div>
            </div>

            <div className="overflow-x-auto bg-gray-700 rounded-lg">
                <table className="w-full text-left text-sm">
                    <thead className="bg-gray-600 text-gray-300 uppercase">
                        <tr>
                            <th className="px-6 py-3">Tarih</th>
                            <th className="px-6 py-3">Personel</th>
                            <th className="px-6 py-3">Nöbet Tipi</th>
                            <th className="px-6 py-3">Durum</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-600">
                        {data.map((item: any) => (
                            <tr key={item.ID} className="hover:bg-gray-650 transition-colors">
                                <td className="px-6 py-4">
                                    {format(new Date(item.Schedule.Date), 'dd MMM yyyy, EEEE', { locale: tr })}
                                </td>
                                <td className="px-6 py-4 font-medium">
                                    {item.Schedule.Employee.FirstName} {item.Schedule.Employee.LastName}
                                </td>
                                <td className="px-6 py-4">
                                    <span
                                        className="px-2 py-1 rounded-full text-xs"
                                        style={{ backgroundColor: item.Schedule.ShiftType.Color + '33', color: item.Schedule.ShiftType.Color }}
                                    >
                                        {item.Schedule.ShiftType.Name}
                                    </span>
                                </td>
                                <td className="px-6 py-4">
                                    {item.ID ? (
                                        <span className="flex items-center gap-1 text-green-400">
                                            <ClipboardCheck className="w-4 h-4" /> Bitti
                                        </span>
                                    ) : (
                                        <span className="flex items-center gap-1 text-yellow-500">
                                            <Clock className="w-4 h-4" /> Bekliyor
                                        </span>
                                    )}
                                </td>
                            </tr>
                        ))}
                        {data.length === 0 && (
                            <tr>
                                <td colSpan={4} className="px-6 py-10 text-center text-gray-500">
                                    Bu ay için veri bulunamadı. Önce "Çizelge Planla" kısmından oluşturun.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default ScheduleViewer;
