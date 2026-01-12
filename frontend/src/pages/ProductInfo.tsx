import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import Header from '../components/Header';
import { callAPI } from '../lib/api';
import { Button } from '@/components/ui/button';
import { Carousel } from 'react-responsive-carousel';
import 'react-responsive-carousel/lib/styles/carousel.min.css';
import { ArrowLeft } from 'lucide-react';
import { AspectRatio } from '@radix-ui/react-aspect-ratio';
import { cn } from '@/lib/utils';

const ProductInfo = () => {
  const { productId } = useParams<{ productId: string }>();
  const navigate = useNavigate();

  const { data: product, isLoading } = useQuery({
    queryKey: ['product', productId],
    queryFn: () => productId ? callAPI('getProduct', { productId: parseInt(productId) }) : Promise.reject('No product ID'),
    enabled: !!productId,
  });

  if (isLoading) return <div>Loading...</div>;
  if (!product) return <div>Product not found</div>;

  return (
    <div className="bg-[var(--background)] min-h-screen">
      <Header />
      <div className="my-2 mx-auto w-[95%] border-black border-b-[1px]"></div>
        <Button
          variant="ghost"
          onClick={() => navigate(-1)}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
      <div className="container mx-1 md:mx-2 lg:mx-4 p-4 pt-2">

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <AspectRatio ratio={1} className='p-4 bg-[#FFF3E8] rounded-sm border-[rgba(0,0,0,0.1)] border-[1px]'>
              <Carousel showThumbs={false}>
                {product.images.map((img: any, idx: number) => (
                  <div key={idx}>
                    <img 
                        src={img.url} 
                        alt={img.alt_text} 
                        className={cn("w-full h-full object-contain",
                            img.size_type === 'v-rect' ? 'w-[65%]' : '',
                            img.size_type === 'square' ? 'rotate-45 scale-[0.7]' : '',
                            ''
                        )} />
                  </div>
                ))}
              </Carousel>
          </AspectRatio>

          <div>
            <h2 className="text-3xl font-bold mb-2">{product.product_name}</h2>
            <p className="text-2xl font-semibold text-gray-600 mb-6">MOP$ {product.price}</p>
            <Button className="w-full">加入購物車</Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductInfo;