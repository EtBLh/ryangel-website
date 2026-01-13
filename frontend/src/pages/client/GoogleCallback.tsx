import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { setCredentials } from '../../store/authSlice';
import { setCartId } from '@/store/cartSlice';
import { toast } from 'sonner';

export default function GoogleCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const dispatch = useDispatch();

  useEffect(() => {
    const token = searchParams.get('token');
    const cartId = searchParams.get('cart_id');
    
    if (token) {
        dispatch(setCredentials({ token }));
        if (cartId) {
            dispatch(setCartId(cartId));
        }
        toast.success('Login Successful');
        navigate('/');
    } else {
        toast.error('Login Failed');
        navigate('/');
    }
  }, [searchParams, navigate, dispatch]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
    </div>
  );
}
