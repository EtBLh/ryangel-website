import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ShoppingCart, Minus, Plus, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSelector } from 'react-redux';
import type { RootState } from '../store';
import { callAPI } from '../lib/api';
import type { CartItem } from '../lib/types';
import AuthDrawer from './AuthDrawer';
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer';
import { cn } from '@/lib/utils';
import { useIsMobile } from '@/hooks/use-mobile';
import { useUser } from '@/hooks/useUser';

const CartDrawer = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [isAuthDrawerOpen, setIsAuthDrawerOpen] = useState(false);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const cartId = useSelector((state: RootState) => state.cart.cartId);
  const { data: client } = useUser();
  const isMobile = useIsMobile();

  const { data: cart, isLoading } = useQuery({
    queryKey: ['cart', cartId],
    queryFn: () => callAPI('getCart'),
    enabled: !!cartId, // Always fetch when cartId exists, not just when drawer is open
  });

  const cartItemCount = cart?.items?.reduce((total: number, item: CartItem) => total + item.quantity, 0) || 0;

  const updateQuantityMutation = useMutation({
    mutationFn: ({ cartItemId, quantity }: { cartItemId: number; quantity: number }) =>
      callAPI('updateCartItem', { cartItemId }, { quantity }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
    },
  });

  const removeItemMutation = useMutation({
    mutationFn: (cartItemId: number) =>
      callAPI('removeCartItem', { cartItemId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
    },
  });

  const handleQuantityChange = (cartItemId: number, newQuantity: number) => {
    if (!cartItemId) {
      console.error('Cart item ID is midssing');
      return;
    }
    if (newQuantity < 1) {
      removeItemMutation.mutate(cartItemId);
    } else {
      updateQuantityMutation.mutate({ cartItemId, quantity: newQuantity });
    }
  };

  const handleCheckout = () => {
    if (client) {
      setIsOpen(false);
      navigate('/checkout');
    } else {
      setIsAuthDrawerOpen(true);
    }
  };

  return (
    <>
    <Drawer direction={isMobile ? "bottom" : "right"} open={isOpen} onOpenChange={setIsOpen}>
      <DrawerTrigger asChild>
        <Button variant="ghost" size="icon" className="relative">
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
      </DrawerTrigger>
      <DrawerContent className={cn(`bg-white rounded-lg`, {
          'bottom-2 left-2 right-2 top-auto h-[90vh]': isMobile,
          'right-2 top-2 bottom-2 w-[400px] h-[calc(100vh-1rem)]': !isMobile,
      })}>
        <DrawerHeader>
          <DrawerTitle className='items-center justify-center'>
            <ShoppingCart className="h-[1.2rem] w-[1.2rem] inline-block mr-2" />購物車
          </DrawerTitle>
        </DrawerHeader>
        <div className="flex-1 overflow-y-auto p-4">
          {isLoading ? (
            <p className="text-center text-gray-500">Loading...</p>
          ) : !cartId || !cart?.items?.length? (
            <p className="text-center text-gray-500">幫手買兩張揮春啦( • ̀ω•́ )</p>
          ) : (
            <div className="space-y-4">
              {cart.items.map((item: CartItem, index: number) => {
                return (
                <div key={index} className="flex items-center gap-3 p-3 border rounded-lg">
                  {item.thumbnail_url && (
                    <div className="w-16 h-16 flex-shrink-0 bg-[#FFF3E8] rounded overflow-hidden flex items-center justify-center">
                      <img 
                        src={item.thumbnail_url} 
                        alt={item.product_name} 
                        className="max-w-full max-h-full object-contain"
                      />
                    </div>
                  )}
                  <div className="flex-1">
                    <h4 className="font-medium text-sm">{item.product_name}</h4>
                    {item.size_type && (
                      <p className="text-xs text-gray-500">Size: {item.size_type}</p>
                    )}
                    <p className="text-sm font-semibold">MOP$ {item.unit_price}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8"
                      disabled={updateQuantityMutation.isPending || removeItemMutation.isPending}
                      onClick={() => handleQuantityChange(item.cart_item_id, item.quantity - 1)}
                    >
                      <Minus className="h-3 w-3" />
                    </Button>
                    <span className="w-8 text-center">{item.quantity}</span>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8"
                      disabled={updateQuantityMutation.isPending || removeItemMutation.isPending}
                      onClick={() => handleQuantityChange(item.cart_item_id, item.quantity + 1)}
                    >
                      <Plus className="h-3 w-3" />
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 text-red-500"
                      disabled={updateQuantityMutation.isPending || removeItemMutation.isPending}
                      onClick={() => removeItemMutation.mutate(item.cart_item_id)}
                    >
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
                );
              })}
            </div>
          )}
        </div>
        {cart?.items?.length > 0 && (
          <div className="border-t p-4 space-y-2">
            <div className="flex justify-between text-sm">
              <span>小計:</span>
              <div className="flex flex-col items-end">
                {cart.discount > 0 ? (
                  <>
                    <span className="line-through text-gray-400 text-xs">MOP$ {cart.subtotal.toFixed(2)}</span>
                    <span className="text-green-600 font-medium">MOP$ {cart.discounted_subtotal.toFixed(2)}</span>
                  </>
                ) : (
                  <span>MOP$ {cart.subtotal.toFixed(2)}</span>
                )}
              </div>
            </div>
            <div className="flex justify-between text-sm">
              <span>運費:</span>
              <div className="flex flex-col items-end">
                 {cart.shipping_fee !== cart.discounted_shipping_fee ? (
                  <>
                     <span className="line-through text-gray-400 text-xs">MOP$ {cart.shipping_fee.toFixed(2)}</span>
                     <span className="text-gray-800 font-medium">{`MOP$ ${cart.discounted_shipping_fee.toFixed(2)}`}</span>
                  </>
                 ) : (
                    <span>MOP$ {cart.shipping_fee.toFixed(2)}</span>
                 )}
              </div>
            </div>
            <div className="flex justify-between font-semibold text-lg border-t pt-2">
              <span>Total:</span>
              <span>MOP$ {cart.total.toFixed(2)}</span>
            </div>
            <Button className="w-full" size="lg" onClick={handleCheckout}>
              立刻結帳
            </Button>
          </div>
        )}
      </DrawerContent>
    </Drawer>
    <AuthDrawer isOpen={isAuthDrawerOpen} onOpenChange={setIsAuthDrawerOpen} />
    </>
  );
};

export default CartDrawer;