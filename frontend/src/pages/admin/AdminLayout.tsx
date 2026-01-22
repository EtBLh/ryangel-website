import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import type { RootState } from '../../store';
import { AppSidebar } from '@/components/admin/AppSidebar';

export const AdminLayout = () => {
    const { token } = useSelector((state: RootState) => state.adminAuth);
    const navigate = useNavigate();
    const location = useLocation();

    useEffect(() => {
        if (!token && location.pathname !== '/admin/login') {
            navigate('/admin/login');
        } else if (token && location.pathname === '/admin/login') {
            navigate('/admin');
        } else if (token && location.pathname === '/admin') {
            // Default to dashboard
            navigate('/admin/dashboard');
        }
    }, [token, location.pathname, navigate]);

    if (!token && location.pathname !== '/admin/login') {
        return null;
    }

    if (!token || location.pathname === '/admin/login') {
        // Login page doesn't need sidebar
        return <Outlet />;
    }

    return (
        <div className="flex w-full min-h-screen">
            <AppSidebar />
            <main className="flex-1 w-full bg-slate-50 overflow-auto">
                <div className="p-4 flex items-center border-b bg-white shadow-sm">
                    {/* Placeholder for hamburger menu if needed */}
                    <h1 className="text-lg font-semibold">Admin Panel</h1>
                </div>
                <div className="p-6">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

