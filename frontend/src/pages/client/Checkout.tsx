import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { useQuery } from '@tanstack/react-query';
import type { RootState } from '../../store';
import { clearCart } from '../../store/cartSlice';
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '../../components/ui/accordion';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Separator } from '../../components/ui/separator';
import { callAPI } from '../../lib/api';
import type { Cart, CartItem, EbuyStore } from '../../lib/types';
import { toast } from 'sonner';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useUser } from '@/hooks/useUser';
import { cn } from '@/lib/utils';
import { Check, ChevronsUpDown } from "lucide-react";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';

const checkoutSchema = z.object({
  storeId: z.string().min(1, "Please select a pickup store"),
  name: z.string().min(1, "Recipient name is required"),
  phone: z.string().min(1, "Phone number is required"),
});

export default function Checkout() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { data: client } = useUser();
  const cartId = useSelector((state: RootState) => state.cart.cartId);
  
  const { data: cart, isLoading: cartLoading } = useQuery<Cart>({
    queryKey: ['cart', cartId],
    queryFn: () => callAPI('getCart'),
    enabled: !!cartId,
  });

  const { data: stores = [], isLoading: storesLoading } = useQuery<{data: EbuyStore[]}>({
    queryKey: ['ebuyStores'],
    queryFn: () => callAPI('getEbuyStores'),
    select: res => res.data
  });

  const loading = cartLoading || storesLoading;
  const [submitting, setSubmitting] = useState(false);
  const [open, setOpen] = useState(false);
  const [activeItem, setActiveItem] = useState("item-info");
  const [completedSteps, setCompletedSteps] = useState<Set<string>>(new Set(["item-1"]));

  const [proofFile, setProofFile] = useState<File | null>(null);

  const form = useForm<z.infer<typeof checkoutSchema>>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
        storeId: '',
        name: client?.username || '',
        phone: client?.phone || '',
    }
  });

  const selectedStoreId = form.watch("storeId");
  const selectedStore = stores.find((s) => String(s.store_id) === selectedStoreId);

  // Initialize form with client data if available
  useEffect(() => {
    if (client) {
      if (client.username) form.setValue('name', client.username);
      if (client.phone) form.setValue('phone', client.phone);
    }
  }, [client, form]);
  
  // Listen for open-cart-accordion event
  useEffect(() => {
    const handleOpenCart = () => {
        setActiveItem("item-1");
    };
    window.addEventListener("open-cart-accordion", handleOpenCart);
    return () => window.removeEventListener("open-cart-accordion", handleOpenCart);
  }, []);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setProofFile(e.target.files[0]);
    }
  };

  const handleStep2Next = async () => {
    const valid = await form.trigger(["storeId", "name", "phone"]);
    if (valid) {
        setCompletedSteps(prev => new Set(prev).add("item-2"));
        setActiveItem("item-3");
    }
  };

  const handleStep3Next = () => {
     setCompletedSteps(prev => new Set(prev).add("item-3"));
     setActiveItem("item-4");
  };

  const onSubmit = async (data: z.infer<typeof checkoutSchema>) => {
    if (!proofFile) {
        toast.error('Please upload payment proof');
        return;
    }

    setSubmitting(true);
    const formData = new FormData();
    formData.append('ebuy_store_id', data.storeId);
    formData.append('name', data.name);
    formData.append('phone', data.phone);
    formData.append('payment_proof', proofFile);

    try {
        await callAPI('createOrder', undefined, formData);
        dispatch(clearCart());
        toast.success('Order placed successfully!');
        navigate('/'); // Redirect to Home
    } catch (err: any) {
        console.error(err);
        toast.error(err.response?.data?.message || 'Failed to place order');
    } finally {
        setSubmitting(false);
    }
  };

  if (loading) return <div className="p-8 text-center">Loading checkout...</div>;
  if (!cart || cart.items.length === 0) return <div className="p-8 text-center">Your cart is empty</div>;

  return (
    <div className="container mx-auto max-w-2xl shadow-sm mt-2">
      <h1 className="text-lg font-semibold mb-2 text-center">訂單確認</h1>
      
      <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
      <Accordion type="single" collapsible value={activeItem} onValueChange={setActiveItem} className="w-full space-y-4">
        
        {/* Step 1: Cart Content */}
        <AccordionItem value="item-1" className="border rounded-lg px-2 mx-2 border-[rgba(0,0,0,0.3)] border-[1px]" disabled={false}>
            <AccordionTrigger className="hover:no-underline py-4">
                <span className="font-medium text-lg">1. 購物車 ({cart.items.length})</span>
            </AccordionTrigger>
            <AccordionContent>
                <div className="space-y-4 pt-2 pb-4">
                    {cart.items.map((item: CartItem) => (
                        <div key={item.cart_item_id} className="flex gap-4 items-center">
                            <div className="w-16 h-16 bg-gray-100 rounded overflow-hidden flex-shrink-0 border">
                                {item.thumbnail_url ? (
                                    <img src={item.thumbnail_url} alt={item.product_name} className="w-full h-full object-cover" />
                                ) : (
                                    <div className="w-full h-full flex items-center justify-center text-gray-400 text-xs text-center p-1">No Img</div>
                                )}
                            </div>
                            <div className="flex-1 min-w-0">
                                <h4 className="font-medium truncate">{item.product_name}</h4>
                                <p className="text-sm text-gray-500">Size: {item.size_type} | Qty: {item.quantity}</p>
                            </div>
                            <div className="font-semibold whitespace-nowrap">
                                ${item.unit_price * item.quantity}
                            </div>
                        </div>
                    ))}
                    <Separator />
                    <div className="space-y-1 pt-2">
                        <div className="flex justify-between items-center text-sm">
                            <span className="text-gray-600">Subtotal</span>
                            <span>MOP$ {cart.subtotal}</span>
                        </div>
                        <div className="flex justify-between items-center text-sm">
                            <span className="text-gray-600">Shipping</span>
                            <span>MOP$ {cart.shipping_fee}</span>
                        </div>
                        {cart.discount > 0 && (
                            <div className="flex justify-between items-center text-sm text-green-600">
                                <span>Discount</span>
                                <span>-MOP$ {cart.discount}</span>
                            </div>
                        )}
                        <div className="flex justify-between items-center font-bold text-lg pt-2 border-t mt-2">
                            <span>Total</span>
                            <span>MOP$ {cart.total}</span>
                        </div>
                    </div>
                </div>
            </AccordionContent>
        </AccordionItem>

        {/* Info Step */}
        <AccordionItem value="item-info" className="border rounded-lg px-2 mx-2 border-[rgba(0,0,0,0.3)] border-[1px]" disabled={false}>
            <AccordionTrigger className="hover:no-underline py-4">
                <span className="font-medium text-lg">重要須知 / Important Notice</span>
            </AccordionTrigger>
            <AccordionContent>
                 <ul className="space-y-4 pt-2 pb-4">
                    <li>
                        <span className="block font-semibold text-sm text-gray-900">檢查資料</span>
                        <p className='text-sm'>下單時請務必核對款式、數量及資料。下單後如需修改，請直接聯絡官方 Instagram / WeChat: <span className="font-mono bg-gray-100 px-1">@ryangel_collection</span>。</p>
                    </li>

                    <li>
                        <span className="block font-semibold text-sm text-gray-900">付款憑證</span>
                        <p className='text-sm'>所有訂單以<span className="font-semibold ">成功上傳付款截圖</span>為準。截圖須包含留言備註之手機號碼。未付款或未上傳截圖之訂單將視為無效。</p>
                    </li>

                    <li>
                        <span className="block font-semibold text-sm text-gray-900">出貨時間</span>
                        <p className='text-sm'>訂單將於付款後 <span className="font-semibold ">3 - 7 天內</span>寄出。如遇庫存不足需補貨，收貨時間將按實際情況延長。</p>
                    </li>

                    <li>
                        <span className="block font-semibold text-sm text-gray-900">到貨通知與取貨時限</span>
                        <p className="mb-2 text-sm">貨到後 eBuy 將透過短訊、官方 APP 或微信公眾號通知。請務必於指定時間內取貨：</p>
                        <ul className="ml-6 list-disc space-y-1">
                            <li className='text-sm'>eBuy 門店：<span className="font-semibold">7 天內</span></li>
                            <li className='text-sm'>24H 智能櫃：<span className="font-semibold ">48 小時內</span></li>
                        </ul>
                        <p className="mt-2 text-sm text-gray-400">註：逾期取貨 eBuy 將向閣下收取逾期費用。</p>
                    </li>

                    <li>
                        <span className="block font-semibold text-sm text-gray-900">版權聲明</span>
                        <p className='text-sm'>小店所有物品均為原創設計並受版權保護。售出貨品<span className="font-semibold">只作私人用途</span>，不得出售、轉售或作其他商業用途。</p>
                    </li>
                </ul>
                <div className="flex justify-end pt-2 pb-4">
                    <Button type="button" onClick={() => setActiveItem("item-2")}>
                        下一步: 取貨資料
                    </Button>
                </div>
            </AccordionContent>
        </AccordionItem>

        {/* Step 2: Ebuy Location & Info */}
        <AccordionItem value="item-2" className="border rounded-lg px-4" disabled={false}>
            <AccordionTrigger className="hover:no-underline py-4">
                <span className="font-semibold text-lg">2. 取貨地點及資料</span>
            </AccordionTrigger>
            <AccordionContent>
                <div className="space-y-4 pt-2 pb-4">
                    <FormField
                        control={form.control}
                        name="storeId"
                        render={({ field }) => (
                            <FormItem className="flex flex-col">
                                <FormLabel>選擇Ebuy收貨地點</FormLabel>
                                <Popover open={open} onOpenChange={setOpen}>
                                    <PopoverTrigger asChild>
                                        <FormControl>
                                            <Button
                                                variant="outline"
                                                role="combobox"
                                                aria-expanded={open}
                                                className={cn(
                                                    "w-full justify-between",
                                                    !field.value && "text-muted-foreground"
                                                )}
                                            >
                                                {field.value
                                                    ? stores.find((store) => String(store.store_id) === field.value)?.store_name
                                                    : "-- 選擇門市 --"}
                                                <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                            </Button>
                                        </FormControl>
                                    </PopoverTrigger>
                                    <PopoverContent className="w-[--radix-popover-trigger-width] p-0">
                                        <Command>
                                            <CommandInput placeholder="Search store..." />
                                            <CommandList>
                                                <CommandEmpty>No store found.</CommandEmpty>
                                                <CommandGroup>
                                                    {stores.map((store) => (
                                                        <CommandItem
                                                            value={store.store_name}
                                                            key={store.store_id}
                                                            onSelect={() => {
                                                                form.setValue("storeId", String(store.store_id));
                                                                setOpen(false);
                                                            }}
                                                        >
                                                            <Check
                                                                className={cn(
                                                                    "mr-2 h-4 w-4",
                                                                    String(store.store_id) === field.value
                                                                        ? "opacity-100"
                                                                        : "opacity-0"
                                                                )}
                                                            />
                                                            {store.store_name}
                                                        </CommandItem>
                                                    ))}
                                                </CommandGroup>
                                            </CommandList>
                                        </Command>
                                    </PopoverContent>
                                </Popover>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    {selectedStore?.latitude && selectedStore?.longitude && (
                        <div className="rounded-md border overflow-hidden h-64 mt-2">
                            <iframe
                                width="100%"
                                height="100%"
                                style={{ border: 0 }}
                                loading="lazy"
                                allowFullScreen
                                src={`https://maps.google.com/maps?q=${selectedStore.latitude},${selectedStore.longitude}+(${encodeURIComponent(selectedStore.store_name)})&z=15&output=embed`}
                            ></iframe>
                        </div>
                    )}

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>收貨人姓名</FormLabel>
                                    <FormControl>
                                        <Input placeholder="收貨人姓名" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="phone"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>收貨人電話號碼</FormLabel>
                                    <FormControl>
                                        <Input placeholder="收貨人電話號碼" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                    </div>
                    
                    <div className="flex justify-end pt-4">
                        <Button 
                            type="button" 
                            onClick={handleStep2Next}
                            disabled={!form.formState.isValid && (!selectedStoreId || !form.getValues().name || !form.getValues().phone)}
                        >
                            下一步: 付款
                        </Button>
                    </div>
                </div>
            </AccordionContent>
        </AccordionItem>

        {/* Step 3: Payment */}
        <AccordionItem value="item-3" className="border rounded-lg px-4" disabled={!completedSteps.has("item-2")}>
            <AccordionTrigger className="hover:no-underline py-4">
                <span className="font-semibold text-lg">3. 付款</span>
            </AccordionTrigger>
            <AccordionContent>
                <div className="flex flex-col items-center space-y-4 pt-2 pb-4">
                    <p className="text-center text-gray-700">請掃描二維碼支付 <strong>${cart.total}</strong>。</p>
                    <div className="w-64 h-64 bg-white flex items-center justify-center rounded-lg border-2 border-gray-200 shadow-sm">
                        {/* Placeholder for QR Code */}
                        <div className="text-center">
                            <div className="text-4xl mb-2">📱</div>
                            <span className="text-gray-500 font-medium">QR Code Placeholder</span>
                        </div>
                    </div>
                            <p>付款時務必備註「下單聯絡電話」以作身份確認。</p>
                    <div className="flex justify-end w-full pt-4">
                         <Button type="button" onClick={handleStep3Next}>
                            下一步: 上傳付款憑證
                        </Button>
                    </div>
                </div>
            </AccordionContent>
        </AccordionItem>

        {/* Step 4: Proof Upload */}
        <AccordionItem value="item-4" className="border rounded-lg px-4" disabled={!completedSteps.has("item-3")}>
            <AccordionTrigger className="hover:no-underline py-4">
                <span className="font-semibold text-lg">4. 上傳付款憑證</span>
            </AccordionTrigger>
            <AccordionContent>
                <div className="space-y-4 pt-2 pb-4">
                     <FormItem>
                        <FormLabel>上傳付款截圖</FormLabel>
                        <FormControl>
                             <Input type="file" accept="image/*" onChange={handleFileChange} className="cursor-pointer" />
                        </FormControl>
                        <p className="text-xs text-muted-foreground">Supported formats: JPG, PNG. Max size: 5MB.</p>
                     </FormItem>
                    
                    <Button 
                        type="submit"
                        className="w-full mt-6 text-lg py-6 font-bold" 
                        size="lg"
                        disabled={submitting || !proofFile}
                    >
                        {submitting ? '提交中...' : `完成訂單`}
                    </Button>
                </div>
            </AccordionContent>
        </AccordionItem>

      </Accordion>
      </form>
      </Form>
    </div>
  );
}