import { useState, useEffect } from 'react';
import { ShoppingCart } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer';

const CartDrawer = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [cartItemCount] = useState(3); // Example: 3 items in cart

  useEffect(() => {
    const checkIsMobile = () => {
      setIsMobile(window.innerWidth < 768); // md breakpoint
    };

    checkIsMobile();
    window.addEventListener('resize', checkIsMobile);

    return () => window.removeEventListener('resize', checkIsMobile);
  }, []);

  return (
    <Drawer direction={isMobile ? "bottom" : "right"} open={isOpen} onOpenChange={setIsOpen}>
      <DrawerTrigger asChild>
        <Button variant="ghost" size="icon" className="relative">
          <ShoppingCart className="h-4 w-4" />
          {cartItemCount > 0 && (
            <Badge
              variant="destructive"
              className="absolute -top-1.5 -right-1.5 h-5 w-5 flex items-center justify-center p-0 text-xs rounded-full"
            >
              {cartItemCount > 99 ? '99+' : cartItemCount}
            </Badge>
          )}
        </Button>
      </DrawerTrigger>
      <DrawerContent className={`max-h-[calc(100vh-2rem)] bg-white rounded-lg ${
        isMobile
          ? 'bottom-2 left-2 right-2 top-auto h-[80vh]'
          : 'right-2 top-4 bottom-4 w-[400px]'
      }`}>
        <DrawerHeader>
          <DrawerTitle>購物車</DrawerTitle>
        </DrawerHeader>
        <div className="p-4">
          <p className="text-center text-gray-500">Your cart is empty</p>
        </div>
      </DrawerContent>
    </Drawer>
  );
};

export default CartDrawer;