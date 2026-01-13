import { useInfiniteQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { callAPI } from '../../lib/api';
import type { Product } from '../../lib/types';
import Masonry, { ResponsiveMasonry } from "react-responsive-masonry"
import { useInView } from 'react-intersection-observer';
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
import FaichunAvaiSizes from '@/components/FaichunAvaiSizes';

import { useSearchParams } from 'react-router-dom';
import { useEffect } from 'react';

const Home = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const { ref, inView } = useInView();

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery({
    queryKey: ['products', query],
    queryFn: ({ pageParam = 1 }) => callAPI('getProducts', { q: query, page: pageParam, page_size: 20 }),
    getNextPageParam: (lastPage: any) => {
      if (lastPage.meta.page < lastPage.meta.total_pages) {
        return lastPage.meta.page + 1;
      }
      return undefined;
    },
    initialPageParam: 1,
  });

  useEffect(() => {
    if (inView && hasNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, fetchNextPage]);

  const products = data?.pages.flatMap((page: any) => page.data) || [];

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="bg-[var(--background)] min-h-screen">
      <div className="container w-full md:w-[800px] lg:w-[1100px] mx-auto p-4">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink href="/">RyAngel</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator>
              <SlashIcon />
            </BreadcrumbSeparator>
            <BreadcrumbItem>
              <BreadcrumbLink href="/">Catalog</BreadcrumbLink>
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
          gutterBreakPoints={{ 350: "8px", 750: "12px", 900: "24px" }}
        >
          <Masonry>
            {products.map((product: Product) => (
              <div
                key={product.product_id}
                className="hover:border-[black] border-[1px] cursor-pointer bg-[#FFF3E8] p-2 rounded-sm flex flex-col gap-2 items-center justify-center"
                onClick={() => navigate(`/product/${product.product_id}`)}
              >
                <img
                  src={product.images[0]?.thumbnail_url || product.images[0]?.url}
                  alt={product.images[0]?.alt_text}
                  className={cn(
                    "object-cover hover:border-black border-[1px]",
                    product.images[0].size_type === 'square' ? 'rotate-45 scale-[0.7]' : '',
                    product.images[0].size_type === 'v-rect' ? 'w-[65%]' : '',
                  )}
                />
                <div className='py-1 w-full flex flex-row justify-between items-center'>
                  <div className="flex flex-col items-start">
                    <h3 className="text-lg md:text-xl ">{product.product_name}</h3>
                    <p className="text-gray-600 text-xs md:text-sm">MOP$ {product.price}</p>
                  </div>
                  <FaichunAvaiSizes sizes={product.available_sizes} />
                </div>
              </div>
            ))}
          </Masonry>
        </ResponsiveMasonry>
        <div ref={ref} className="h-10 w-full flex justify-center p-4">
            {isFetchingNextPage && <span className="text-gray-500">Loading more...</span>}
            {!hasNextPage && products.length > 0 && <span className="text-gray-400">~~~ Fin ~~~</span>}
        </div>
      </div>
    </div>
  );
};

export default Home;