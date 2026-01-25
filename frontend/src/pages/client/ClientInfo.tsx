import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Calendar } from "@/components/ui/calendar";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from '@/components/ui/input';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useUser } from '@/hooks/useUser';
import { callAPI } from '@/lib/api';
import { dict } from '@/lib/dict';
import type { OrderItem } from '@/lib/types';
import { cn } from "@/lib/utils";
import { logout } from '@/store/authSlice';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { format } from "date-fns";
import { CalendarIcon, Check, Pencil, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import * as z from 'zod';

const profileSchema = z.object({
  username: z.string().optional(),
  email: z.string().email("Invalid email address").optional().or(z.literal('')),
  date_of_birth: z.date().optional(),
});

const ClientInfo = () => {
  const { data: client } = useUser();
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [isEditing, setIsEditing] = useState(false);

  const form = useForm<z.infer<typeof profileSchema>>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      username: '',
      email: '',
    },
  });

  useEffect(() => {
    if (client) {
      form.reset({
        username: client.username || '',
        email: client.email || '',
        date_of_birth: client.date_of_birth ? new Date(client.date_of_birth) : undefined
      });
    }
  }, [client, form]);

  const handleUpdateProfile = async (data: z.infer<typeof profileSchema>) => {
    try {
      const payload = {
        ...data,
        date_of_birth: data.date_of_birth?.toISOString() // Backend expects ISO string
      };
      await callAPI('updateClient', undefined, payload);
      toast.success('資料已更新');
      setIsEditing(false);
      queryClient.invalidateQueries({ queryKey: ['user'] });
    } catch (error: any) {
      const msg = error.response?.data?.error?.message || '更新失敗';
      toast.error(msg);
    }
  };

  const { data: ordersData, isLoading } = useQuery({
    queryKey: ['orders', client?.phone],
    queryFn: () => callAPI('getOrders'),
    enabled: !!client,
  });

  const orders = ordersData?.orders || [];

  const handleLogout = () => {
    dispatch(logout());
    toast.success('已登出');
    navigate('/');
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'bg-yellow-500 hover:bg-yellow-600';
      case 'confirmed': return 'bg-blue-500 hover:bg-blue-600';
      case 'shipped': return 'bg-purple-500 hover:bg-purple-600';
      case 'delivered': return 'bg-green-500 hover:bg-green-600';
      case 'cancelled': return 'bg-red-500 hover:bg-red-600';
      case 'refunded': return 'bg-gray-500 hover:bg-gray-600';
      default: return 'bg-gray-500';
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('zh-MO', { style: 'currency', currency: 'MOP' }).format(price);
  };

  if (!client) {
    return (
      <div className="container mx-auto p-4 text-center">
        <p className="mb-4">請先登入</p>
        <Button onClick={() => navigate('/')}>回首頁</Button>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4 flex flex-col md:flex-row gap-4">
      <div className="flex flex-col w-full md:w-1/3">
        <h1 className="text-xl font-bold my-2">會員中心</h1>

        <div className="bg-card border-gray-400 border-[1px] p-3 rounded-md border shadow-sm space-y-1 w-full">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold">基本資料</h2>
            {!isEditing ? (
              <Button variant="ghost" size="sm" onClick={() => setIsEditing(true)}>
                <Pencil className="h-4 w-4" />
                編輯
              </Button>
            ) : (
              <div className="flex gap-2">
                <Button variant="ghost" size="sm" onClick={() => { setIsEditing(false); form.reset(); }}>
                  <X className="h-4 w-4" />
                </Button>
                <Button variant="default" size="sm" onClick={form.handleSubmit(handleUpdateProfile)}>
                  <Check className="h-4 w-4" />
                </Button>
              </div>
            )}
          </div>

          <Form {...form}>
            <div className="grid grid-cols-1 gap-1">
              <div>
                <label className="text-sm font-medium text-gray-500">電話號碼</label>
                <p className="text-md font-mono text-gray-500 bg-gray-100 p-2 rounded mb-2">{client.phone} <span className="text-xs text-gray-400 ml-2"></span></p>
              </div>

              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem className='space-y-0'>
                    <FormLabel className="text-sm font-medium text-gray-500">用戶名稱</FormLabel>
                    {isEditing ? (
                      <FormControl>
                        <Input placeholder="輸入用戶名稱" {...field} />
                      </FormControl>
                    ) : (
                      <p className="text-md pb-2">{client.username || '-'}</p>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem className='space-y-0'>
                    <FormLabel className="text-sm font-medium text-gray-500">Email</FormLabel>
                    {isEditing ? (
                      <FormControl>
                        <Input placeholder="example@email.com" {...field} />
                      </FormControl>
                    ) : (
                      <p className="text-md pb-2">{client.email || '-'}</p>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="date_of_birth"
                render={({ field }) => (
                  <FormItem className="flex flex-col">
                    <FormLabel className="text-sm font-medium text-gray-500">Birthday</FormLabel>
                    {isEditing ? (
                      <Popover>
                        <PopoverTrigger asChild>
                          <FormControl>
                            <Button
                              variant={"outline"}
                              className={cn(
                                "text-left font-normal",
                                !field.value && "text-muted-foreground"
                              )}
                            >
                              {field.value ? (
                                format(field.value, "PPP")
                              ) : (
                                <span>選擇日期</span>
                              )}
                              <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                            </Button>
                          </FormControl>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                          <Calendar
                            mode="single"
                            selected={field.value}
                            onSelect={field.onChange}
                            disabled={(date) =>
                              date > new Date() || date < new Date("1900-01-01")
                            }
                            initialFocus
                          />
                        </PopoverContent>
                      </Popover>
                    ) : (
                      <p className="text-md mb-2">{client.date_of_birth ? format(new Date(client.date_of_birth), 'PPP') : '-'}</p>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </Form>
        </div>
        <Button variant="destructive" onClick={handleLogout} className="my-2 mr-0 ml-auto block border border-gray-400 border-[1px]">
          登出
        </Button>
      </div>

      <div className="flex flex-col flex-1">

        <div className="space-y-4">
          <h2 className="text-lg font-semibold">歷史訂單</h2>
          {isLoading ? (
            <div className="p-4 text-center text-gray-500">載入中...</div>
          ) : orders.length === 0 ? (
            <div className="p-8 rounded-lg text-center text-gray-500 border-gray-400 border-[1px] bg-[#FFF4EB]">
              尚無訂單記錄
            </div>
          ) : (
            <div className="space-y-4 w-full">
              {orders.map((orderData: any) => {
                // Backward compatibility: check if it's new structure {order, items} or old structure Order object
                const order = orderData.order || orderData;
                const items = orderData.items || [];

                return (
                  <div key={order.order_id} className="p-4 rounded-lg border-gray-400 border-[1px] bg-[#FFF4EB] w-full">
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="flex items-center space-x-2">
                          <p className="font-semibold text-md md:text-lg">#{order.order_number}</p>
                        </div>
                        <p className="text-sm text-gray-500 mt-1">
                          {new Date(order.order_date).toLocaleString('zh-TW')}
                        </p>
                        {order.ebuy_store_name && (
                          <p className="text-sm font-medium text-blue-800 mt-1">
                            取貨地址: {order.ebuy_store_name}
                          </p>
                        )}
                      </div>
                      <div className="text-right">
                        <p className="font-bold text-lg text-primary">
                          {formatPrice(order.total_amount)}
                        </p>
                        <p className="text-xs text-gray-500 capitalize">
                          {/* {order.payment_method.replace('_', ' ')} */}
                          支付方式: 聚易用
                        </p>

                        <Badge className={getStatusColor(order.order_status)}>
                          {dict[order.order_status as keyof typeof dict]}
                        </Badge>
                      </div>
                    </div>

                    <div className="mt-4 border-t pt-4">
                      <h4 className="text-sm font-semibold mb-2">訂單內容</h4>
                      <div className="space-y-2">
                        {items && items.map((item: OrderItem) => (
                          <div key={item.order_item_id} className="flex justify-between text-sm items-center bg-white border-gray-200 border-[1px] p-2 rounded">
                            <div className="flex items-center space-x-3">
                              {item.product_image && (
                                <img
                                  src={item.product_image.startsWith('http') ? item.product_image : item.product_image.replace(/(\.[^.]+)$/, '-sm$1')}
                                  alt={item.product_name}
                                  className="w-24 h-24 object-contain rounded"
                                />
                              )}
                              <div>
                                <p className="font-medium">{item.product_name}</p>
                                {item.size_type && <p className="text-xs text-gray-500">尺寸: {dict[item.size_type as keyof typeof dict]}</p>}
                              </div>
                            </div>
                            <div className="text-right">
                              <span className="text-gray-600">x{item.quantity}</span>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div className="text-sm text-gray-600 mt-4 border-t pt-2 space-y-1">
                      <div className="flex justify-between">
                        <span>商品小計:</span>
                        <span>{formatPrice(order.subtotal_amount)}</span>
                      </div>
                      {order.discount_amount > 0 && (
                        <div className="flex justify-between text-green-600">
                          <span>折扣:</span>
                          <span>-{formatPrice(order.discount_amount)}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span>運費:</span>
                        <span>{formatPrice(order.shipping_amount)}</span>
                      </div>
                    </div>

                    {order.payment_proof && (
                      <div className="mt-3 pt-2 border-t border-gray-200">
                        <p className="text-xs font-semibold text-gray-500 mb-2">付款憑證:</p>
                        <div className="relative w-24 h-24 group">
                          <img
                            src={`${import.meta.env.VITE_API_ROOT}${order.payment_proof}`}
                            alt="Payment Proof"
                            className="w-full h-full object-cover rounded-md border shadow-sm cursor-pointer transition-transform hover:scale-105"
                            onClick={(e) => {
                              e.stopPropagation();
                              window.open(`${import.meta.env.VITE_API_ROOT}${order.payment_proof}`, '_blank');
                            }}
                          />
                        </div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ClientInfo;
