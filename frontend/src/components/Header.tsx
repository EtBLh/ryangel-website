import { useState } from 'react';
import { Search, User, ShoppingCart, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { useLocation, useSearchParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import CartDrawer from './CartDrawer';
import AuthDrawer from './AuthDrawer';
import logoWithText from '@public/logo-with-text.png';
import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useUser';
import type { RootState } from '../store';
import { useQuery } from '@tanstack/react-query';
import { callAPI } from '../lib/api';
import type { Cart } from '../lib/types';

const Header = () => {
  const [authOpen, setAuthOpen] = useState(false);
  const [authView, setAuthView] = useState<'login' | 'register'>('login');
  
  const { data: client } = useUser();
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams] = useSearchParams();
  
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [searchQuery, setSearchQuery] = useState(searchParams.get('q') || '');
  const [showNotification, setShowNotification] = useState(true);

  const cartId = useSelector((state: RootState) => state.cart.cartId);

    const { data: cart } = useQuery<Cart>({
    queryKey: ['cart', cartId],
    queryFn: () => callAPI('getCart'),
    enabled: !!cartId,
  });

  const cartItemCount = cart?.items?.reduce((sum, item) => sum + item.quantity, 0) || 0;

  const handleUserClick = () => {
    if (client) {
      navigate('/client-info');
    } else {
      setAuthView('login');
      setAuthOpen(true);
    }
  };

  const handleCartClick = () => {
      window.dispatchEvent(new CustomEvent('open-cart-accordion'));
  };

  const handleSearchSubmit = (e?: React.FormEvent) => {
    e?.preventDefault();
    navigate(`/?q=${searchQuery}`);
  };

  return (
    <>
      <header className="bg-[var(--background)] mt-0 md:mb-2 flex flex-col justify-center px-2 py-2 fixed top-0 left-0 right-0 z-10 border-b-[1px] border-gray-400 min-h-[60px]">
        {showNotification && (
          <div className="bg-[#c11e02] text-white px-4 py-1.5 flex justify-between items-center text-xs md:text-sm -mx-2 -mt-2 md:mb-2">
             <div className="flex-1 text-center font-medium">新年期間不限款式及尺寸買三送一，買滿三張免運費</div>
             <button onClick={() => setShowNotification(false)} className="text-white hover:text-gray-300">
               <X className="h-4 w-4" />
             </button>
          </div>
        )}
        
        {/* Top Row: Logo & Icons (and Desktop Search) */}
        <div className="flex justify-between items-center w-full h-[60px] md:h-auto">
            <h1 className="hidden">RyAngel</h1>
            <img 
            src={logoWithText} 
            alt="RyAngel Logo" 
            className='w-[180px] md:w-[200px] cursor-pointer object-contain md:static md:translate-x-0 md:left-auto' // Center on mobile, left on desktop
            onClick={() => navigate('/')}
            />
            
            {/* Spacer for mobile to balance the absolute logo */}
            <div className="md:hidden w-[40px]"></div> 

            <div className="flex items-center space-x-2 ml-auto">
            {/* Desktop Search */}
            <div className="hidden md:flex items-center transition-all duration-300 ease-in-out">
                {isSearchExpanded ? (
                    <form onSubmit={handleSearchSubmit} className="flex items-center gap-2 animate-in fade-in slide-in-from-right-4 duration-300">
                        <Input 
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            placeholder="搜尋店內商品..."
                            className="w-[200px] h-8"
                            autoFocus
                            onBlur={() => { if(!searchQuery) setIsSearchExpanded(false) }}
                        />
                        <Button type="button" variant="ghost" size="icon" className="h-8 w-8" onClick={() => { setSearchQuery(''); setIsSearchExpanded(false); }}>
                            <X className="h-4 w-4" />
                        </Button>
                    </form>
                ) : (
                    <Button variant="ghost" size="icon" onClick={() => setIsSearchExpanded(true)}>
                        <Search className="h-4 w-4" />
                    </Button>
                )}
            </div>

            {location.pathname === '/checkout' ? (
                <Button variant="ghost" size="icon" className="relative" onClick={handleCartClick}>
                <ShoppingCart className="h-4 w-4" />
                {cartItemCount > 0 && (
                    <Badge
                    variant="destructive"
                    className="border-[#1F3D39] border-[1px] absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs rounded-full"
                    >
                    {cartItemCount > 99 ? '99+' : cartItemCount}
                    </Badge>
                )}
                </Button>
            ) : (
                <CartDrawer />
            )}
            <Button variant="ghost" size="icon" onClick={handleUserClick}>
                <User className="h-4 w-4" />
            </Button>
            </div>
        </div>

        {/* Mobile Search - Bottom Row */}
        <div className="md:hidden w-full pb-2">
            <form onSubmit={handleSearchSubmit} className="relative">
                <Input 
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="搜尋店內商品..."
                    className="w-full pr-8"
                />
                <Button type="submit" variant="ghost" size="sm" className="absolute right-0 top-0 h-full px-3 text-gray-500">
                    <Search className="h-4 w-4" />
                </Button>
            </form>
        </div>
      </header>
      <AuthDrawer isOpen={authOpen} onOpenChange={setAuthOpen} defaultView={authView} />
      {/* Spacer with responsive height */}
      <div className={`transition-all duration-300 ${showNotification ? 'h-[150px] md:h-[100px]' : 'h-[110px] md:h-[60px]'}`}></div>
    </>
  );
};

export default Header;
