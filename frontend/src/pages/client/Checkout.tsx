import {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useSelector, useDispatch} from 'react-redux';
import {useQuery} from '@tanstack/react-query';
import type {RootState} from '../../store';
import {clearCart} from '../../store/cartSlice';
import {Accordion, AccordionContent, AccordionItem, AccordionTrigger} from '../../components/ui/accordion';
import {Button} from '../../components/ui/button';
import {Input} from '../../components/ui/input';
import {callAPI} from '../../lib/api';
import type {Cart, CartItem, EbuyStore} from '../../lib/types';
import {toast} from 'sonner';
import {useForm} from 'react-hook-form';
import {zodResolver} from '@hookform/resolvers/zod';
import * as z from 'zod';
import {useUser} from '@/hooks/useUser';
import {cn} from '@/lib/utils';
import {Check, ChevronsUpDown} from "lucide-react";
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
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import {useIsMobile} from '@/hooks/use-mobile';
import {dict} from '@/lib/dict';

const checkoutSchema = z.object({
    shippingMethod: z.enum(['ebuy_store', 'direct_address']),
    storeId: z.string().optional(),
    address: z.string().optional(),
    name: z.string().min(1, "Recipient name is required"),
    phone: z.string().min(1, "Phone number is required"),
}).refine(
    (data) => {
        if (data.shippingMethod === 'ebuy_store') {
            return !!data.storeId && data.storeId.length > 0;
        }
        if (data.shippingMethod === 'direct_address') {
            return !!data.address && data.address.length > 0;
        }
        return false;
    },
    {
        message: "Please provide required shipping information",
        path: ['shippingMethod'],
    }
);

export default function Checkout() {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const {data: client} = useUser();
    const cartId = useSelector((state: RootState) => state.cart.cartId);

    const form = useForm<z.infer<typeof checkoutSchema>>({
        resolver: zodResolver(checkoutSchema),
        defaultValues: {
            shippingMethod: 'ebuy_store',
            storeId: '',
            address: '',
            name: client?.username || '',
            phone: client?.phone || '',
        }
    });

    const shippingMethod = form.watch("shippingMethod");
    const selectedStoreId = form.watch("storeId");

    const {data: cart, isLoading: cartLoading} = useQuery<Cart>({
        queryKey: ['cart', cartId, shippingMethod],
        queryFn: () => callAPI('getCart', { shipping_method: shippingMethod }),
        enabled: !!cartId,
    });

    const {data: stores = [], isLoading: storesLoading} = useQuery<{ data: EbuyStore[] }, Error, EbuyStore[]>({
        queryKey: ['ebuyStores'],
        queryFn: () => callAPI('getEbuyStores'),
        select: res => res.data
    });

    const loading = cartLoading || storesLoading;
    const isMobile = useIsMobile();
    const [submitting, setSubmitting] = useState(false);
    const [open, setOpen] = useState(false);
    const [activeItem, setActiveItem] = useState("item-info");
    const [completedSteps, setCompletedSteps] = useState<Set<string>>(new Set(["item-info"]));

    const [proofFile, setProofFile] = useState<File | null>(null);
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
            setActiveItem("item-info");
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
        const fieldsToValidate = ["shippingMethod", "name", "phone"];
        if (shippingMethod === 'ebuy_store') {
            fieldsToValidate.push("storeId");
        } else if (shippingMethod === 'direct_address') {
            fieldsToValidate.push("address");
        }
        const valid = await form.trigger(fieldsToValidate as any);
        if (valid) {
            setCompletedSteps(prev => new Set(prev).add("item-2"));
            if (isMobile) {
                setActiveItem("item-cart");
            } else {
                setActiveItem("item-3");
            }
        }
    };

    const handleStepCartNext = () => {
        setCompletedSteps(prev => new Set(prev).add("item-cart"));
        setActiveItem("item-3");
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
        formData.append('shipping_method', data.shippingMethod);
        if (data.shippingMethod === 'ebuy_store' && data.storeId) {
            formData.append('ebuy_store_id', data.storeId);
        } else if (data.shippingMethod === 'direct_address' && data.address) {
            formData.append('address', data.address);
        }
        formData.append('name', data.name);
        formData.append('phone', data.phone);
        formData.append('payment_proof', proofFile);

        try {
            await callAPI('createOrder', undefined, formData);
            dispatch(clearCart());
            toast.success('成功下單!');
            navigate('/client-info'); // Redirect to Client Info
        } catch (err: any) {
            console.error(err);
            toast.error(err.response?.data?.message || 'Failed to place order');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) return <div className="p-8 text-center">Loading checkout...</div>;
    if (!cart || cart.items.length === 0) return <div className="p-8 text-center">Your cart is empty</div>;

    // Cart Summary Component (reusable for both desktop and mobile)
    const CartSummary = () => (
        <div className="space-y-4">
            {cart.items.map((item: CartItem) => (
                <div key={item.cart_item_id} className="flex gap-4 items-center">
                    <div className="w-24 h-24 border-gray-300 rounded overflow-hidden flex-shrink-0 border">
                        {item.thumbnail_url ? (
                            <img src={item.thumbnail_url} alt={item.product_name}
                                 className="w-full h-full object-contain"/>
                        ) : (
                            <div
                                className="w-full h-full flex items-center justify-center text-gray-400 text-xs text-center p-1">No
                                Img</div>
                        )}
                    </div>
                    <div className="flex-1 min-w-0">
                        <h4 className="font-medium truncate">{item.product_name}</h4>
                        <p className="text-sm text-gray-500">{item.size_type ? dict[item.size_type] || item.size_type : ''}</p>
                        <p className="text-xs text-gray-500">Qty: {item.quantity}張</p>
                    </div>
                    <div className="font-semibold whitespace-nowrap">
                        ${item.unit_price * item.quantity}
                    </div>
                </div>
            ))}
            <div className="space-y-1 pt-2">
                <div className="flex justify-between items-center text-sm">
                    <span className="text-gray-600">小計</span>
                    <span>MOP$ {cart.subtotal}</span>
                </div>
                <div className="flex justify-between items-start text-sm">
                    <span className="text-gray-600">運費</span>
                    <div className='flex flex-col text-right gap-1'>
                        <span>MOP$ {cart.shipping_fee}</span>
                        <span
                            className={cn('text-green-600', {'hidden': cart.discounted_shipping_fee >= cart.shipping_fee})}>-MOP$ {cart.shipping_fee - cart.discounted_shipping_fee}</span>
                    </div>
                </div>
                {cart.discount > 0 && (
                    <div className="flex justify-between items-center text-sm text-green-600">
                        <span>折扣</span>
                        <span>-MOP$ {cart.discount}</span>
                    </div>
                )}
                <div className="flex justify-between items-center font-bold text-lg pt-2 border-t mt-2">
                    <span>總計</span>
                    <span>MOP$ {cart.total}</span>
                </div>
            </div>
        </div>
    );

    return (
        <div className="container mx-auto mt-2 px-2">
            <h1 className="text-lg font-semibold mb-4 text-center">訂單確認</h1>

            {/* Responsive Layout */}
            <div className="grid grid-cols-1 lg:grid-cols-[400px_1fr] lg:gap-6 lg:max-w-6xl lg:mx-auto">
                {/* Desktop Left Column: Cart Summary */}
                <div className="hidden lg:block space-y-4">
                    <div className="border-black bg-[#F7EDE4] shadow-md rounded-lg border sticky top-4">
                        <h2 className="font-semibold text-md px-2 py-1 m-0 text-gray-100 bg-gray-800 rounded-t-lg">訂單摘要</h2>
                        <div className='p-2'>
                            <CartSummary/>
                        </div>
                    </div>
                </div>

                {/* Main Content (Mobile: Full width, Desktop: Right Column) */}
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)}>
                        <Accordion type="single" value={activeItem} onValueChange={setActiveItem}
                                   className="w-full space-y-1">

                            {/* Desktop: No cart accordion, it's in the sidebar */}

                            {/* Info Step */}
                            <AccordionItem
                                value="item-info"
                                className={cn("border rounded-lg mx-2 border-gray-400 border-[1px] bg-[#F7EDE4]")}
                                disabled={false}>
                                <AccordionTrigger className={cn("hover:no-underline py-2 px-4 rounded-t-lg", {
                                    'bg-gray-800 text-white': activeItem === 'item-info',
                                })}>
                                <span className="font-medium text-md flex items-center gap-2">
                                    <span
                                        className={cn('flex items-center justify-center rounded-full w-[1.2rem] h-[1.2rem]', {
                                            'bg-[#F7EDE4] text-black': activeItem === 'item-info',
                                            'border-[1px] border-gray-400 text-black': activeItem !== 'item-info'
                                        })}>
                                        i
                                    </span>
                                    重要須知 / Important Notice
                                </span>
                                </AccordionTrigger>
                                <AccordionContent className='p-2'>
                                    <ul className="space-y-4 pt-2 pb-4">
                                        <li>
                                            <span className="block font-semibold text-sm text-gray-900">檢查資料</span>
                                            <p className='text-sm'>下單時請務必核對款式、數量及資料。下單後如需修改，請直接聯絡官方
                                                Instagram / WeChat: <span
                                                    className="font-mono bg-gray-100 px-1">@ryangel_collection</span>。
                                            </p>
                                        </li>

                                        <li>
                                            <span className="block font-semibold text-sm text-gray-900">付款憑證</span>
                                            <p className='text-sm'>所有訂單以<span
                                                className="font-semibold ">成功上傳付款截圖</span>為準。截圖須包含留言備註之手機號碼。未付款或未上傳截圖之訂單將視為無效。
                                            </p>
                                        </li>

                                        <li>
                                            <span className="block font-semibold text-sm text-gray-900">出貨時間</span>
                                            <p className='text-sm'>訂單將於付款後 <span className="font-semibold ">3 - 7 天內</span>寄出。如遇庫存不足需補貨，收貨時間將按實際情況延長。
                                            </p>
                                        </li>

                                        <li>
                                            <span
                                                className="block font-semibold text-sm text-gray-900">到貨通知與取貨時限</span>
                                            <p className="mb-2 text-sm">貨到後 eBuy 將透過短訊、官方 APP
                                                或微信公眾號通知。請務必於指定時間內取貨：</p>
                                            <ul className="ml-6 list-disc space-y-1">
                                                <li className='text-sm'>eBuy 門店：<span
                                                    className="font-semibold">7 天內</span></li>
                                                <li className='text-sm'>24H 智能櫃：<span className="font-semibold ">48 小時內</span>
                                                </li>
                                            </ul>
                                            <p className="mt-2 text-sm text-gray-400">註：逾期取貨 eBuy
                                                將向閣下收取逾期費用。</p>
                                        </li>

                                        <li>
                                            <span className="block font-semibold text-sm text-gray-900">版權聲明</span>
                                            <p className='text-sm'>小店所有物品均為原創設計並受版權保護。售出貨品<span
                                                className="font-semibold">只作私人用途</span>，不得出售、轉售或作其他商業用途。
                                            </p>
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
                            <AccordionItem
                                value="item-2"
                                className={cn("border rounded-lg mx-2 border-gray-400 border-[1px] bg-[#F7EDE4]")}
                                disabled={false}>
                                <AccordionTrigger className={cn("hover:no-underline py-2 px-4 rounded-t-lg", {
                                    'bg-gray-800 text-white': activeItem === 'item-2',
                                })}>
                                <span className="font-medium text-md flex items-center gap-2">
                                    <span
                                        className={cn('flex items-center justify-center rounded-full w-[1.2rem] h-[1.2rem]', {
                                            'bg-[#F7EDE4] text-black': activeItem === 'item-2',
                                            'border-[1px] border-gray-400 text-black': activeItem !== 'item-2'
                                        })}>
                                        2
                                    </span>
                                    取貨地點及資料
                                </span>
                                </AccordionTrigger>
                                <AccordionContent className='p-2'>
                                    <div className="space-y-4 pt-2 pb-4">
                                        <FormField
                                            control={form.control}
                                            name="shippingMethod"
                                            render={({field}) => (
                                                <FormItem>
                                                    <FormLabel>選擇收貨方式</FormLabel>
                                                    <FormControl>
                                                        <RadioGroup value={field.value} onValueChange={(v) => field.onChange(v)}>
                                                            <div className="flex gap-4">
                                                                <label className={cn(
                                                                    "relative flex-1 border-[1px] rounded-lg p-4 cursor-pointer transition-all",
                                                                    field.value === 'ebuy_store' ? "border-gray-800 bg-[var(--background)]" : "border-gray-300"
                                                                )}>
                                                                    <div className="absolute top-2 right-2">
                                                                        <RadioGroupItem value="ebuy_store" id="ebuy_store" className='scale-[0.6]'/>
                                                                    </div>
                                                                    <div className="flex items-start gap-3">
                                                                        <div className="font-medium">澳門本地Ebuy取貨 ($5運費)</div>
                                                                    </div>
                                                                </label>
                                                                <label className={cn(
                                                                    "relative flex-1 border-[1px] rounded-lg p-4 cursor-pointer transition-all",
                                                                    field.value === 'direct_address' ? "border-gray-800 bg-[var(--background)]" : "border-gray-300"
                                                                )}>
                                                                    <div className="absolute top-2 right-2">
                                                                        <RadioGroupItem value="direct_address" id="direct_address" className='scale-[0.6]'/>
                                                                    </div>
                                                                    <div className="flex items-start gap-3">
                                                                        <div className="font-medium">香港送貨 ($40運費)</div>
                                                                    </div>
                                                                </label>
                                                            </div>
                                                        </RadioGroup>
                                                    </FormControl>
                                                    <FormMessage/>
                                                </FormItem>
                                            )}
                                        />
                                        {/*)}*/}
                                        {/**/}
                                        {shippingMethod === 'ebuy_store' && (
                                            <FormField
                                                control={form.control}
                                                name="storeId"
                                                render={({field}) => (
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
                                                                            : "選擇取貨地點"}
                                                                        <ChevronsUpDown
                                                                            className="ml-2 h-4 w-4 shrink-0 opacity-50"/>
                                                                    </Button>
                                                                </FormControl>
                                                            </PopoverTrigger>
                                                            <PopoverContent
                                                                className="w-[--radix-popover-trigger-width] p-0">
                                                                <Command>
                                                                    <CommandInput placeholder="搜尋ebuy取貨點"/>
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
                                                        <FormMessage/>
                                                    </FormItem>
                                                )}
                                            />
                                        )}

                                        {shippingMethod === 'ebuy_store' && selectedStore?.latitude && selectedStore?.longitude && (
                                            <div className="rounded-md border overflow-hidden h-64 mt-2">
                                                <iframe
                                                    width="100%"
                                                    height="100%"
                                                    style={{border: 0}}
                                                    loading="lazy"
                                                    allowFullScreen
                                                    src={`https://maps.google.com/maps?q=${selectedStore.latitude},${selectedStore.longitude}+(${encodeURIComponent(selectedStore.store_name)})&z=15&output=embed`}
                                                ></iframe>
                                            </div>
                                        )}

                                        {shippingMethod === 'direct_address' && (
                                            <FormField
                                                control={form.control}
                                                name="address"
                                                render={({field}) => (
                                                    <FormItem>
                                                        <FormLabel>送貨地址</FormLabel>
                                                        <FormControl>
                                                            <Input placeholder="請輸入完整送貨地址" {...field} />
                                                        </FormControl>
                                                        <FormMessage/>
                                                    </FormItem>
                                                )}
                                            />
                                        )}

                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                            <FormField
                                                control={form.control}
                                                name="name"
                                                render={({field}) => (
                                                    <FormItem>
                                                        <FormLabel>收貨人姓名</FormLabel>
                                                        <FormControl>
                                                            <Input placeholder="收貨人姓名" {...field} />
                                                        </FormControl>
                                                        <FormMessage/>
                                                    </FormItem>
                                                )}
                                            />
                                            <FormField
                                                control={form.control}
                                                name="phone"
                                                render={({field}) => (
                                                    <FormItem>
                                                        <FormLabel>收貨人電話號碼</FormLabel>
                                                        <FormControl>
                                                            <Input placeholder="收貨人電話號碼" {...field} />
                                                        </FormControl>
                                                        <FormMessage/>
                                                    </FormItem>
                                                )}
                                            />
                                        </div>

                                        <div className="flex justify-end pt-4">
                                            <Button
                                                type="button"
                                                onClick={handleStep2Next}
                                            >
                                                下一步: 付款
                                            </Button>
                                        </div>
                                    </div>
                                </AccordionContent>
                            </AccordionItem>

                            {/* Step Cart (Mobile Only) */}
                            {isMobile && (
                                <AccordionItem
                                    value="item-cart"
                                    className={cn("border rounded-lg mx-2 border-gray-400 border-[1px] bg-[#F7EDE4]")}
                                    disabled={!completedSteps.has("item-2")}>
                                    <AccordionTrigger className={cn("hover:no-underline py-2 px-4 rounded-t-lg", {
                                        'bg-gray-800 text-white': activeItem === 'item-cart',
                                    })}>
                                    <span className="font-medium text-md flex items-center gap-2">
                                        <span
                                            className={cn('flex items-center justify-center rounded-full w-[1.2rem] h-[1.2rem]', {
                                                'bg-[#F7EDE4] text-black': activeItem === 'item-cart',
                                                'border-[1px] border-gray-400 text-black': activeItem !== 'item-cart'
                                            })}>
                                            3
                                        </span>
                                        購物車確認
                                    </span>
                                    </AccordionTrigger>
                                    <AccordionContent className='p-2'>
                                        <CartSummary/>
                                        <div className="flex justify-end pt-4">
                                            <Button type="button" onClick={handleStepCartNext}>
                                                下一步: 付款
                                            </Button>
                                        </div>
                                    </AccordionContent>
                                </AccordionItem>
                            )}

                            {/* Step 3: Payment */}
                            <AccordionItem
                                value="item-3"
                                className={cn("border rounded-lg mx-2 border-gray-400 border-[1px] bg-[#F7EDE4]")}
                                disabled={isMobile ? !completedSteps.has("item-cart") : !completedSteps.has("item-2")}>
                                <AccordionTrigger className={cn("hover:no-underline py-2 px-4 rounded-t-lg", {
                                    'bg-gray-800 text-white': activeItem === 'item-3',
                                })}>
                                <span className="font-medium text-md flex items-center gap-2">
                                    <span
                                        className={cn('flex items-center justify-center rounded-full w-[1.2rem] h-[1.2rem]', {
                                            'bg-[#F7EDE4] text-black': activeItem === 'item-3',
                                            'border-[1px] border-gray-400 text-black': activeItem !== 'item-3'
                                        })}>
                                        {isMobile ? 4 : 3}
                                    </span>
                                    付款
                                </span>
                                </AccordionTrigger>
                                <AccordionContent>
                                    <div className="flex flex-col items-center space-y-4 p-2 pb-0 pt-4">
                                        <p className="text-center text-lg text-gray-700">請掃描二維碼支付 <strong>${cart.total}</strong>。
                                        </p>
                                        <div className='flex flex-row'>
                                            <div
                                                className="flex flex-col gap-2 bg-white flex items-center justify-center rounded-lg border-[1rem] border-white shadow-sm">
                                                {/* Placeholder for QR Code */}
                                                <span className='text-md'>聚易用支付碼</span>

                                                <img className='w-64 max-w-2/3' src="/media/ryangel-payment-qrcode2.png"
                                                     alt="scan this qrcode to pay"/>
                                            </div>
                                        </div>
                                        <p className='font-semibold'>付款時務必備註「下單聯絡電話」以作身份確認。</p>
                                        <div className="flex justify-end w-full pt-4">
                                            <Button type="button" onClick={handleStep3Next}>
                                                下一步: 上傳付款憑證
                                            </Button>
                                        </div>
                                    </div>
                                </AccordionContent>
                            </AccordionItem>

                            {/* Step 4: Proof Upload */}
                            <AccordionItem
                                value="item-4"
                                className={cn("border rounded-lg mx-2 border-gray-400 border-[1px] bg-[#F7EDE4]")}
                                disabled={!completedSteps.has("item-3")}>
                                <AccordionTrigger className={cn("hover:no-underline py-2 px-4 rounded-t-lg", {
                                    'bg-gray-800 text-white': activeItem === 'item-4',
                                })}>
                                <span className="font-medium text-md flex items-center gap-2">
                                    <span
                                        className={cn('flex items-center justify-center rounded-full w-[1.2rem] h-[1.2rem]', {
                                            'bg-[#F7EDE4] text-black': activeItem === 'item-4',
                                            'border-[1px] border-gray-400 text-black': activeItem !== 'item-4'
                                        })}>
                                        {isMobile ? 5 : 4}
                                    </span>
                                    上傳付款憑證
                                </span>
                                </AccordionTrigger>
                                <AccordionContent className='p-4 pb-0'>
                                    <div className="space-y-4 pt-2 pb-4">
                                        <FormItem>
                                            <FormLabel>上傳付款截圖</FormLabel>
                                            <FormControl>
                                                <Input type="file" accept="image/*" onChange={handleFileChange}
                                                       className="cursor-pointer"/>
                                            </FormControl>
                                            <p className="text-xs text-muted-foreground">支持以下圖片格式: JPG, PNG.
                                                最大檔案大小: 5MB.</p>
                                        </FormItem>
                                        <div className='h-10'/>
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
        </div>
    );
}