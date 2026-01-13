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
            <span className='font-muted text-sm text-right'>尺寸</span>
            <div className="flex flex-row gap-1.5 items-center">
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
                'px-2 py-1 w-3 h-8': size === 'v-rect',
                'px-2 py-1 mx-1 w-6 h-6 rotate-45': size === 'square',
                'px-2 py-1 w-6 h-8': size === 'fat-v-rect',
                'px-2 py-1 mx-2 w-8 h-8 rotate-45': size === 'big-square',
            }, className)}
            onClick={() => {
                onClick && onClick(size);
            }}
        />
    )
}

export default FaichunAvaiSizes;