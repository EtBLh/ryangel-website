import { Link, useLocation } from 'react-router-dom';
import { cn } from '@/lib/utils';
import { useDispatch } from 'react-redux';
import { adminLogout } from '@/store/adminAuthSlice';
import { callAPI } from '@/lib/api';

const menuItems = [
    { name: 'Dashboard', path: '/admin/dashboard' },
    { name: 'Orders', path: '/admin/orders' },
    { name: 'Products', path: '/admin/products' },
    { name: 'Discounts', path: '/admin/discounts' },
    { name: 'Users', path: '/admin/users' },
];

export const AppSidebar = () => {
    const location = useLocation();
    const dispatch = useDispatch();

    const handleLogout = async () => {
        try {
            await callAPI('adminLogout');
            dispatch(adminLogout());
             // Redirect handled by AdminLayout or useEffect
             window.location.href = "/admin/login";
        } catch (error) {
            console.error(error);
            // Force logout client side even if server fails
            dispatch(adminLogout());
            window.location.href = "/admin/login";
        }
    };

    return (
        <aside className="w-64 bg-slate-900 text-slate-100 min-h-screen flex flex-col">
            <div className="p-4 border-b border-slate-700">
                <h1 className="text-xl font-bold">Ryangel Admin</h1>
            </div>
            <nav className="flex-1 p-4 space-y-2">
                {menuItems.map((item) => (
                    <Link
                        key={item.path}
                        to={item.path}
                        className={cn(
                            "flex items-center px-4 py-2 rounded-md transition-colors",
                            location.pathname === item.path || location.pathname.startsWith(item.path + '/')
                                ? "bg-slate-700 text-white"
                                : "text-slate-400 hover:bg-slate-800 hover:text-white"
                        )}
                    >
                        {item.name}
                    </Link>
                ))}
            </nav>
            <div className="p-4 border-t border-slate-700">
                <button 
                    onClick={handleLogout}
                    className="w-full text-left px-4 py-2 text-slate-400 hover:bg-slate-800 hover:text-white rounded-md transition-colors"
                >
                    Logout
                </button>
            </div>
        </aside>
    );
};
