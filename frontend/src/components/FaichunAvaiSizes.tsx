import type { SizeType } from "@/lib/types";
import { cn } from "@/lib/utils";

interface FaichunAvaiSizesProps {
    sizes: SizeType[];
    className?: string;
    onClick?: (size: SizeType) => void;
}

const FaichunAvaiSizes = ({ sizes, className, onClick }: FaichunAvaiSizesProps) => {
    return (
        <div className={cn('flex flex-col justify-start items-end gap-1', className)}>
            <span className='font-muted text-xs md:text-sm text-right'>尺寸</span>
            <div className="flex flex-row gap-1 items-center">
                {
                    sizes.map((size, idx) => <FaichunSizeIcon key={idx} size={size} onClick={onClick} />)
                }
            </div>
        </div>
    )
}

interface FaichunSizeIconProps {
    size: SizeType;
    className?: string;
    onClick?: (size: SizeType) => void;
}

export const FaichunSizeIcon = ({ size, className, onClick }: FaichunSizeIconProps) => {
    return (
        <span
            className={cn("bg-destructive border-[#1F3D39] border-[1px] hover:bg-destructive/60", {
                'px-1 py-0.5 w-4 h-7': size === 'v-rect',
                'px-1 py-0.5 mx-1 w-5 h-5 rotate-45': size === 'square',
                'px-1 py-0.5 w-6 h-7': size === 'fat-v-rect',
                'px-1 py-0.5 mx-2 w-6 h-6 rotate-45': size === 'big-square',
            }, className)}
            onClick={() => {
                onClick && onClick(size);
            }}
        />
    )
}

export default FaichunAvaiSizes;