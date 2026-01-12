import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import { callAPI } from '../lib/api';
import type { Product } from '../lib/types';
import Masonry, { ResponsiveMasonry } from "react-responsive-masonry"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { SlashIcon } from "lucide-react"


import { cn } from '@/lib/utils';

const Home = () => {
  const navigate = useNavigate();
  const { data: products, isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => callAPI('getProducts'),
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="bg-[var(--background)] min-h-screen">
      <Header />
      <div className="my-2 mx-auto w-[95%] border-[rgba(0,0,0,0.3)] border-b-[1px]"></div>
      <div className="container mx-1 md:mx-2 lg:mx-4 p-4">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink href="/">RyAngel</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator>
              <SlashIcon />
            </BreadcrumbSeparator>
            <BreadcrumbItem>
              <BreadcrumbLink href="/components">Catelog</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator>
              <SlashIcon />
            </BreadcrumbSeparator>
            <BreadcrumbItem>
              <BreadcrumbPage>揮春</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <h2 className='text-2xl my-2'>
          新年揮春
        </h2>
        <ResponsiveMasonry
          columnsCountBreakPoints={{ 350: 2, 750: 3, 900: 4 }}
          gutterBreakpoints={{ 350: "8px", 750: "12px", 900: "24px" }}
        >
          <Masonry>
            {products?.data?.map((product: Product) => (
              <div
                key={product.product_id}
                className="hover:border-[black] border-[1px] cursor-pointer bg-[#FFF3E8] p-2 rounded-sm flex flex-col gap-2 items-center justify-center"
                onClick={() => navigate(`/product/${product.product_id}`)}
              >
                <img
                  src={product.images[0]?.url}
                  alt={product.images[0]?.alt_text}
                  className={cn(
                    "object-cover hover:border-black border-[1px]",
                    product.images[0].size_type === 'square' ? 'rotate-45 scale-[0.7]' : '',
                    product.images[0].size_type === 'v-rect' ? 'w-[65%]' : '',
                  )}
                />
                <div className='py-1 w-full flex flex-row justify-between items-center'>
                  <div className="flex flex-col items-start">
                    <h3 className="text-xl">{product.product_name}</h3>
                    <p className="text-gray-600 text-sm">MOP$ {product.price}</p>
                  </div>
                  <div className='flex flex-col'>
                    <span className='font-muted text-sm text-right'>形狀</span>
                    <div className="flex flex-row gap-1.5 items-center">
                      {
                        product.available_sizes.map((size, idx) => (
                          size === 'v-rect' ?
                            <span key={idx} className="px-2 py-1 w-3 h-8 bg-destructive border-[#1F3D39] border-[1px]"/> :
                          size === 'square' ?
                            <span key={idx} className="px-2 py-1 mx-1 w-6 h-6 rotate-45 bg-destructive border-[#1F3D39] border-[1px]"/> :
                          size === 'fat-v-rect' ?
                            <span key={idx} className="px-2 py-1 w-6 h-8 bg-destructive border-[#1F3D39] border-[1px]"/> :
                          null
                        ))
                      }
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </Masonry>
        </ResponsiveMasonry>
      </div>
    </div>
  );
};

export default Home;