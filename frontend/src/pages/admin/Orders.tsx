import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle
} from '@/components/ui/dialog';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { callAPI } from '@/lib/api';
import { ChevronLeft, ChevronRight, Eye, MoreHorizontal } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

interface Order {
    order_id: number;
    order_number: string;
    order_date: string;
    total_amount: number;
    order_status: string;
    customer_notes?: string;
    payment_method: string;
    payment_proof?: string;
    client_name: string;
    client_phone: string;
    ebuy_store_name?: string;
}

interface OrderItem {
    product_name: string;
    product_sku: string;
    quantity: number;
    total_price: number;
    unit_price: number;
    product_image?: string;
    size_type?: string;
}

const ORDER_STATUSES = ['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded'];

export const Orders = () => {
    const [orders, setOrders] = useState<Order[]>([]);
    const [page, setPage] = useState(1);
    const [isLoading, setIsLoading] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [orderItems, setOrderItems] = useState<OrderItem[]>([]);
    const [isItemsLoading, setIsItemsLoading] = useState(false);
    const [isDetailsOpen, setIsDetailsOpen] = useState(false);

    const fetchOrders = async () => {
        setIsLoading(true);
        try {
            // pass page query param
            //@ts-ignore
            const data = await callAPI('adminGetOrders', undefined, undefined, { params: { page } });
            setOrders(data);
        } catch (error) {
            console.error("Failed to fetch orders", error);
            toast.error("Failed to load orders");
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchOrders();
    }, [page]);

    const handleViewOrder = async (order: Order) => {
        setSelectedOrder(order);
        setIsDetailsOpen(true);
        setIsItemsLoading(true);
        try {
            const items = await callAPI('adminGetOrderItems', { orderId: order.order_id.toString() });
            setOrderItems(items);
        } catch (error) {
            toast.error("Failed to load order items");
        } finally {
            setIsItemsLoading(false);
        }
    };

    const handleStatusChange = async (orderId: number, newStatus: string) => {
        try {
            await callAPI('adminUpdateOrderStatus', { orderId: orderId.toString() }, { status: newStatus });
            toast.success(`Order status updated to ${newStatus}`);
            // Optimistic update
            setOrders(prev => prev.map(o => o.order_id === orderId ? { ...o, order_status: newStatus } : o));
        } catch (error) {
            toast.error("Failed to update status");
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-3xl font-bold tracking-tight">Orders</h2>
                <div className="flex space-x-2">
                    <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                    >
                        <ChevronLeft className="h-4 w-4" />
                    </Button>
                    <span className="flex items-center text-sm">Page {page}</span>
                    <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => setPage(p => p + 1)}
                        disabled={orders.length < 50}
                    >
                        <ChevronRight className="h-4 w-4" />
                    </Button>
                </div>
            </div>

            <div className="border rounded-md">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Order #</TableHead>
                            <TableHead>Date</TableHead>
                            <TableHead>Client</TableHead>
                            <TableHead>Store</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead>Total</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={7} className="text-center h-24">Loading...</TableCell>
                            </TableRow>
                        ) : orders.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} className="text-center h-24">No orders found.</TableCell>
                            </TableRow>
                        ) : (
                            orders.map((order) => (
                                <TableRow key={order.order_id}>
                                    <TableCell className="font-medium">{order.order_number}</TableCell>
                                    <TableCell>{new Date(order.order_date).toLocaleDateString()}</TableCell>
                                    <TableCell>
                                        <div className="flex flex-col">
                                            <span className="font-medium">{order.client_name}</span>
                                            <span className="text-xs text-muted-foreground">{order.client_phone}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell>{order.ebuy_store_name || '-'}</TableCell>
                                    <TableCell>
                                        <Badge variant={order.order_status === 'pending' ? 'secondary' : 'default'}>
                                            {order.order_status}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>MOP${order.total_amount}</TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end items-center space-x-2">
                                            <Button variant="ghost" size="icon" onClick={() => handleViewOrder(order)}>
                                                <Eye className="h-4 w-4" />
                                            </Button>
                                            
                                            <DropdownMenu>
                                                <DropdownMenuTrigger asChild>
                                                    <Button variant="ghost" size="icon">
                                                        <MoreHorizontal className="h-4 w-4" />
                                                    </Button>
                                                </DropdownMenuTrigger>
                                                <DropdownMenuContent align="end">
                                                    {ORDER_STATUSES.map(status => (
                                                        <DropdownMenuItem 
                                                            key={status}
                                                            onClick={() => handleStatusChange(order.order_id, status)}
                                                            disabled={order.order_status === status}
                                                        >
                                                            Mark as {status}
                                                        </DropdownMenuItem>
                                                    ))}
                                                </DropdownMenuContent>
                                            </DropdownMenu>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            <Dialog open={isDetailsOpen} onOpenChange={setIsDetailsOpen}>
                <DialogContent className="max-w-3xl bg-white">
                    <DialogHeader>
                        <DialogTitle>Order Details: {selectedOrder?.order_number}</DialogTitle>
                    </DialogHeader>
                    
                    <div className="grid grid-cols-2 gap-4 py-4">
                        <div>
                            <h4 className="font-semibold mb-2">Order Info</h4>
                            <p className="text-sm">Date: {selectedOrder && new Date(selectedOrder.order_date).toLocaleString()}</p>
                            <p className="text-sm">Status: {selectedOrder?.order_status}</p>
                            <p className="text-sm">Payment: {selectedOrder?.payment_method}</p>
                            <p className="text-sm">Store: {selectedOrder?.ebuy_store_name || "N/A"}</p>
                        </div>
                        <div>
                             <h4 className="font-semibold mb-2">Customer</h4>
                             <p className="text-sm">{selectedOrder?.client_name}</p>
                             <p className="text-sm">{selectedOrder?.client_phone}</p>
                             <h4 className="font-semibold mb-2 mt-2">Notes</h4>
                             <p className="text-sm italic">{selectedOrder?.customer_notes || "None"}</p>
                        </div>
                    </div>

                    <div className="border rounded-md mt-4">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[80px]">Image</TableHead>
                                    <TableHead>Product</TableHead>
                                    <TableHead>Size</TableHead>
                                    <TableHead>SKU</TableHead>
                                    <TableHead className="text-right">Qty</TableHead>
                                    <TableHead className="text-right">Price</TableHead>
                                    <TableHead className="text-right">Total</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {isItemsLoading ? (
                                    <TableRow>
                                        <TableCell colSpan={6} className="text-center">Loading items...</TableCell>
                                    </TableRow>
                                ) : (
                                    orderItems.map((item, idx) => (
                                        <TableRow key={idx}>
                                            <TableCell>
                                                {item.product_image ? (
                                                    <img 
                                                        src={item.product_image.startsWith('http') ? item.product_image : `${item.product_image}`} 
                                                        alt={item.product_name} 
                                                        className="h-32 w-32 object-contain rounded"
                                                    />
                                                ) : (
                                                    <div className="h-32 w-32 bg-gray-100 rounded flex items-center justify-center text-xs text-gray-500">No Img</div>
                                                )}
                                            </TableCell>
                                            <TableCell>{item.product_name}</TableCell>
                                            <TableCell>{item.size_type || '-'}</TableCell>
                                            <TableCell>{item.product_sku}</TableCell>
                                            <TableCell className="text-right">{item.quantity}</TableCell>
                                            <TableCell className="text-right">MOP${item.unit_price}</TableCell>
                                            <TableCell className="text-right">MOP${item.total_price}</TableCell>
                                        </TableRow>
                                    ))
                                )}
                            </TableBody>
                        </Table>
                    </div>
                </DialogContent>
            </Dialog>
        </div>
    );
};
