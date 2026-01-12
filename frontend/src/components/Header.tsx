import { Link } from 'react-router-dom';
import { Search, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import CartDrawer from './CartDrawer';

const Header = () => {
  return (
    <header className="bg-[var(--background)] flex justify-between items-center px-8 py-2">
      <div className="flex items-start space-x-4 justify-center">
        <h1 className="text-xl font-bold">RyAngel</h1>
        <div className="flex flex-col">
          <Link to="/" className="font-sm font-light">
          Catelog
        </Link>
          <Link to="/about" className="font-sm font-light">
          About Us
        </Link>
        </div>
      </div>
      <div className="flex items-center space-x-2">
        <Button variant="ghost" size="icon">
          <Search className="h-4 w-4" />
        </Button>
        <CartDrawer />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <User className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem>登入</DropdownMenuItem>
            <DropdownMenuItem>註冊</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

    </header>
  );
};

export default Header;