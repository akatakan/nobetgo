import { useState, useEffect } from 'react';
import { Bell } from 'lucide-react';
import { notificationApi } from '../services/api';
import type { Notification } from '../types';

export function NotificationBell() {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [isOpen, setIsOpen] = useState(false);
    const [loading, setLoading] = useState(true);

    // Hardcoded for now, assuming current logged in user ID is 1 or something dynamic in real app
    const currentUserId = 1;

    useEffect(() => {
        fetchNotifications();
        const interval = setInterval(fetchNotifications, 60000); // Check every minute
        return () => clearInterval(interval);
    }, []);

    const fetchNotifications = async () => {
        try {
            const response = await notificationApi.getUnread(currentUserId);
            setNotifications(response.data || []);
        } catch (err) {
            console.error('Failed to load notifications', err);
        } finally {
            setLoading(false);
        }
    };

    const markAsRead = async (id: number) => {
        try {
            await notificationApi.markAsRead(id);
            setNotifications(notifications.filter(n => n.ID !== id));
        } catch (err) {
            console.error('Failed to mark notification as read', err);
        }
    };

    const markAllAsRead = async () => {
        try {
            await notificationApi.markAllAsRead(currentUserId);
            setNotifications([]);
            setIsOpen(false);
        } catch (err) {
            console.error('Failed to mark all as read', err);
        }
    };

    return (
        <div className="relative">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="relative p-2 text-gray-400 hover:text-gray-500 focus:outline-none"
            >
                <span className="sr-only">Bildirimleri görüntüle</span>
                <Bell className="h-6 w-6" />
                {notifications.length > 0 && (
                    <span className="absolute top-0 right-0 block h-4 w-4 rounded-full bg-red-500 text-white text-xs text-center leading-4">
                        {notifications.length}
                    </span>
                )}
            </button>

            {isOpen && (
                <div className="origin-top-right absolute right-0 mt-2 w-80 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50">
                    <div className="p-4 border-b border-gray-100 flex justify-between items-center">
                        <h3 className="text-sm font-medium text-gray-900">Bildirimler</h3>
                        {notifications.length > 0 && (
                            <button
                                onClick={markAllAsRead}
                                className="text-xs text-indigo-600 hover:text-indigo-800"
                            >
                                Tümünü Oku
                            </button>
                        )}
                    </div>
                    <div className="max-h-96 overflow-y-auto">
                        {loading ? (
                            <div className="p-4 text-center text-sm text-gray-500">Yükleniyor...</div>
                        ) : notifications.length === 0 ? (
                            <div className="p-4 text-center text-sm text-gray-500">Yeni bildirim yok</div>
                        ) : (
                            <ul className="divide-y divide-gray-100">
                                {notifications.map((notification) => (
                                    <li key={notification.ID} className="p-4 hover:bg-gray-50 cursor-pointer" onClick={() => markAsRead(notification.ID)}>
                                        <p className="text-sm font-medium text-gray-900">{notification.title}</p>
                                        <p className="text-sm text-gray-500 mt-1">{notification.message}</p>
                                        <p className="text-xs text-gray-400 mt-2">
                                            {new Date(notification.CreatedAt).toLocaleString('tr-TR')}
                                        </p>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}
