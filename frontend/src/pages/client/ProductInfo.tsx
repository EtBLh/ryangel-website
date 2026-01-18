import FaichunAvaiSizes, { FaichunSizeIcon } from '@/components/FaichunAvaiSizes';
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import FaichunSizeSq from '@public/FaiChun-size-square.svg?react';
import FaichunSizeVRect from '@public/FaiChun-size-vertical.svg?react';
import FaichunSizeFatVRect from '@public/FaiChun-size-fat.svg?react';
import FaichunSizeBigSquare from '@public/FaiChun-size-bigsquare.svg?react';
import { AspectRatio } from '@radix-ui/react-aspect-ratio';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useDispatch } from 'react-redux';
import { setCartId } from '../../store/cartSlice';
import { ArrowLeft, Minus, Plus, SlashIcon } from 'lucide-react';
import { useState, useEffect } from 'react';
import { Carousel } from 'react-responsive-carousel';
import 'react-responsive-carousel/lib/styles/carousel.min.css';
import { useNavigate, useParams } from 'react-router-dom';
import { callAPI } from '../../lib/api';

import { toast } from 'sonner';
import { dict } from '@/lib/dict';
import { ButtonGroup } from '@/components/ui/button-group';
import type { SizeType } from '@/lib/types';

const ProductInfo = () => {
    const { productId } = useParams<{ productId: string }>();
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const dispatch = useDispatch();

    const { data: product, isLoading } = useQuery({
        queryKey: ['product', productId],
        queryFn: () => productId ? callAPI('getProduct', { productId: parseInt(productId) }) : Promise.reject('No product ID'),
        enabled: !!productId,
    });
    const [selectedItem, setSelectedItem] = useState(0);
    const [selectedSize, setSelectedSize] = useState<SizeType | null>(null);
    const [quantity, setQuantity] = useState(1);

    // Set default selected size to first available size when product loads
    useEffect(() => {
        if (product?.available_sizes?.length > 0 && selectedSize === null) {
            setSelectedSize(product.available_sizes[0]);
        }
    }, [product, selectedSize]);

    const addToCartMutation = useMutation({
        mutationFn: (data: { product_id: number; size_type: SizeType | null; quantity: number }) =>
            callAPI('addToCart', undefined, data),
        onSuccess: (response) => {
            // Save cart_id to Redux store if returned
            if (response.cart_id) {
                dispatch(setCartId(response.cart_id));
            }
            // Invalidate cart queries to refresh cart data
            queryClient.invalidateQueries({ queryKey: ['cart'] });
            toast.success('已加入購物車!');
        },
        onError: (error) => {
            console.error('Failed to add to cart:', error);
            toast.error('加入購物車失敗，請再試一次。');
        },
    });

    const handleAddToCart = () => {
        if (!productId || !selectedSize) {
            toast.error('請選擇尺寸後再加入購物車。');
            return;
        }

        addToCartMutation.mutate({
            product_id: parseInt(productId),
            size_type: selectedSize,
            quantity: quantity,
        });
    };

    const increaseQuantity = () => {
        if (quantity < 5) {
            setQuantity(quantity + 1);
        }
    };

    const decreaseQuantity = () => {
        if (quantity > 1) {
            setQuantity(quantity - 1);
        }
    };

    if (isLoading) return <div>Loading...</div>;
    if (!product) return <div>Product not found</div>;


    return (
        <div className="bg-[var(--background)] min-h-screen">
            <div className="container p-4 pt-2 mx-auto w-full md:w-[800px] lg:w-[1100px]">
                <div className='flex mb-2 items-center gap-2'>
                    <Button
                        variant="ghost"
                        onClick={() => navigate(-1)}
                        className='px-1'
                    >
                        <ArrowLeft className="text-gray-700" />
                    </Button>
                    <Breadcrumb className=''>
                        <BreadcrumbList>
                            <BreadcrumbItem>
                                <BreadcrumbLink href="/">Catalog</BreadcrumbLink>
                            </BreadcrumbItem>
                            <BreadcrumbSeparator>
                                <SlashIcon />
                            </BreadcrumbSeparator>
                            <BreadcrumbItem>
                                <BreadcrumbPage>揮春</BreadcrumbPage>
                            </BreadcrumbItem>
                            <BreadcrumbSeparator>
                                <SlashIcon />
                            </BreadcrumbSeparator>
                            <BreadcrumbItem>
                                <BreadcrumbPage>{product.product_name}</BreadcrumbPage>
                            </BreadcrumbItem>
                        </BreadcrumbList>
                    </Breadcrumb>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-0 md:gap-8">
                    <div className=''>
                        <Carousel
                            showThumbs={true}
                            autoPlay={false}
                            showStatus={false}
                            showIndicators={true}
                            showArrows={true}
                            onClickThumb={(x) => console.log(x)}
                            emulateTouch
                            onChange={setSelectedItem}
                            renderThumbs={() => product.images.map((img: any, idx: number) => (
                                <div key={idx} className="w-full h-20 flex items-center justify-center bg-[#FFF3E8]">
                                    <img
                                        src={img.thumbnail_url || img.url}
                                        alt={img.alt_text}
                                        className="h-full object-contain"
                                    />
                                </div>
                            ))}
                            selectedItem={selectedItem}
                            className='first:bg-red'
                        >
                            {product.images.map((img: any, idx: number) => (
                                <AspectRatio ratio={1} key={idx} className='h-full w-full select-none p-4 bg-[#FFF3E8] rounded-sm border-[rgba(0,0,0,0.1)] border-[1px]'>
                                    <img
                                        src={img.url}
                                        alt={img.alt_text}
                                        className={cn("w-full h-full object-contain select-none",
                                            img.size_type === 'v-rect' ? 'w-[65%]' : '',
                                            img.size_type === 'square' ? 'rotate-45 scale-[0.7]' : '',
                                            ''
                                        )} />
                                </AspectRatio>
                            ))}
                        </Carousel>
                    </div>

                    <div className='relative'>
                        <span className='text-sm font-light text-gray-800 block'>
                            {
                                product.available_sizes.map((size: SizeType) => dict[size as SizeType]).join(' / ')
                            }
                            揮春
                        </span>
                        <h2 className="text-3xl font-semibold my-2">{product.product_name}</h2>
                        <div className="flex justify-start mb-6">
                            <span className='text-md text-gray-600 '>MOP$</span>
                            <span className="text-2xl font-semibold text-gray-600">{product.price}</span>
                        </div>

                        <FaichunAvaiSizes
                            sizes={product.available_sizes}
                            className='absolute top-0 right-0'
                            onClick={console.log}
                        />

                        <Accordion
                            type="single"
                            collapsible
                            className="w-full"
                            defaultValue={"policy"}
                        >
                            <AccordionItem value="size" className='border-y-[1px] border-y-gray-300'>
                                <AccordionTrigger className=''>揮春尺寸</AccordionTrigger>
                                <AccordionContent className="flex flex-col">
                                    <div>
                                        <p className='text-sm text-gray-700'>直款: 13cm x 30cm</p>
                                        <p className='text-sm text-gray-700'>方形: 23cm x 23cm</p>
                                        <p className='text-sm text-gray-700'>巨無霸方形: 28cm x 28cm</p>
                                        <p className='text-sm text-gray-700'>胖胖款: 23*30</p>
                                    </div>
                                    <div className="grid grid-cols-2 gap-2 place-items-center">
                                        <FaichunSizeSq className="max-h-[15vh]" />
                                        <FaichunSizeVRect className="max-h-[15vh]" />
                                        <FaichunSizeFatVRect className="max-h-[15vh]" />
                                        <FaichunSizeBigSquare className="max-h-[15vh]" />
                                    </div>
                                </AccordionContent>
                            </AccordionItem>
                            <AccordionItem value="policy" >
                                <AccordionTrigger className=''>出貨與付款</AccordionTrigger>
                                <AccordionContent className="flex flex-col">
                                    <div className="max-w-4xl mx-auto text-gray-700">

                                        <section className="">
                                            <div className="space-y-2">
                                                <div>
                                                    <p><span className="font-semibold text-md">1. 選購：</span>於網站選購喜歡的款式及數量，系統將自動計算優惠後金額。</p>
                                                </div>
                                                <div>
                                                    <p><span className="font-semibold text-md">2. 驗證：</span>點擊「立刻結帳」並以手機號碼或 Google 帳戶進行身份驗證。</p>
                                                </div>
                                                <div>
                                                    <p><span className="font-semibold text-md">3. 填寫資訊：</span>確認訂單內容，填寫 eBuy 寄貨資料（<span className="font-semibold">收貨人姓名、電話及收貨地點</span>）。</p>
                                                </div>
                                                <div>
                                                    <p><span className="font-semibold text-md">4. 轉帳付款：</span>截圖網頁上的付款碼，使用<span className="font-semibold ">澳門聚易用</span>平台付款。</p>
                                                    <div className="mt-2 ml-6 text-sm">
                                                        <p className="font-semibold underline">付款時務必備註「下單聯絡電話」以作身份確認。</p>
                                                        <p className="text-gray-500 mt-1">支援平台：</p>
                                                        <div className="flex flex-row flex-wrap justify-between mt-1">
                                                            <span>MPay</span>
                                                            <span>澳門中銀手機銀行</span>
                                                            <span>工銀e支付</span>
                                                            <span>豐付寶</span>
                                                            <span>支付寶（澳門）</span>
                                                            <span>廣發移動支付</span>
                                                            <span>LusoPay</span>
                                                            <span>極易付</span>
                                                        </div>
                                                    </div>
                                                </div>
                                                <div>
                                                    <p><span className="font-semibold text-md">5. 上傳憑證：</span>最後<span className="font-semibold">上傳付款截圖</span>到網站即完成訂購。您可以隨時登入網站查看訂單進度。</p>
                                                </div>
                                            </div>
                                        </section>
                                    </div>

                                </AccordionContent>
                            </AccordionItem>
                            <AccordionItem value="faq" >
                                <AccordionTrigger>注意事項</AccordionTrigger>
                                <AccordionContent className="flex flex-col">
                                    <ul className="space-y-4">
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
                                            <p className="mb-2">貨到後 eBuy 將透過短訊、官方 APP 或微信公眾號通知。請務必於指定時間內取貨：</p>
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
                                </AccordionContent>
                            </AccordionItem>
                        </Accordion>

                        {/* Size Selection */}
                        <div className="mt-4">
                            <label className="text-sm font-medium text-gray-700 mb-2 block">
                                選擇尺寸 / Select Size
                            </label>
                            <div className="flex flex-wrap gap-2">
                                {product.available_sizes.map((size: SizeType) => (
                                    <Button
                                        key={size}
                                        variant='outline'
                                        size="sm"
                                        onClick={() => {
                                            if (product.available_sizes.includes(size)) {
                                                setSelectedSize(size);
                                            }
                                        }}
                                        className={cn("flex items-center gap-2 p-3 border-[#262422] border-[1px] hover:bg-[rgba(0,0,0,0.1)]",
                                            {
                                                'bg-[#262422] hover:bg-[#262422] text-white hover:text-white border-[#3D716C]': selectedSize === size,
                                                'text-gray-800': selectedSize !== size,
                                            }
                                        )}
                                    >
                                        <FaichunSizeIcon size={size} className='my-2' />
                                        {dict[size]}
                                    </Button>
                                ))}
                            </div>
                        </div>

                        {/* Quantity Selection */}
                        <div className="mt-4">
                            <label className="text-sm font-medium text-gray-700 mb-2 block">
                                數量 / Quantity
                            </label>
                            <ButtonGroup className="flex items-center">
                                <Button
                                    variant="outline"
                                    size="icon"
                                    onClick={decreaseQuantity}
                                    disabled={quantity <= 1}
                                    className='border-gray-400'
                                >
                                    <Minus className="h-4 w-4" />
                                </Button>
                                <Input
                                    type="number"
                                    value={quantity}
                                    onChange={(e) => {
                                        const value = parseInt(e.target.value);
                                        if (value >= 1 && value <= 5) {
                                            setQuantity(value);
                                        }
                                    }}
                                    className="w-20 text-center border-gray-400"
                                    min={1}
                                    max={5}
                                />
                                <Button
                                    variant="outline"
                                    size="icon"
                                    onClick={increaseQuantity}
                                    disabled={quantity >= 5}
                                    className='border-gray-400'
                                >
                                    <Plus className="h-4 w-4" />
                                </Button>
                            </ButtonGroup>
                        </div>

                        <Button
                            className="w-full mt-4 rounded-lg"
                            onClick={handleAddToCart}
                            disabled={addToCartMutation.isPending || !selectedSize}
                        >
                            {addToCartMutation.isPending ? 'Adding...' : '加入購物車'}
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProductInfo;