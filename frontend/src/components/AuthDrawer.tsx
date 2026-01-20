import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Drawer,
    DrawerContent,
    DrawerHeader,
    DrawerTitle,
    DrawerDescription,
    DrawerFooter,
    DrawerClose,
} from '@/components/ui/drawer';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    InputOTP,
    InputOTPGroup,
    InputOTPSlot,
} from '@/components/ui/input-otp';
import { callAPI } from '@/lib/api';
import { useDispatch, useSelector } from 'react-redux';
import { setCredentials } from '@/store/authSlice';
import { setCartId } from '@/store/cartSlice';
import type { RootState } from '@/store';
import { toast } from 'sonner';
import { useIsMobile } from '@/hooks/use-mobile';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';

const COUNTRY_CODES = [
    { code: '+853', label: '澳門 (+853)' },
    { code: '+852', label: '香港 (+852)' },
    { code: '+886', label: '台灣 (+886)' },
    { code: '+86', label: '中國大陸 (+86)' },
];

const registerSchema = z.object({
    countryCode: z.string(),
    phone: z.string().min(1, "請輸入手機號碼"),
    username: z.string().min(1, "請輸入使用者名稱"),
    password: z.string().min(6, "密碼必須至少6個字符"),
});

const passwordLoginSchema = z.object({
    username: z.string().min(1, "請輸入使用者名稱"),
    password: z.string().min(1, "請輸入密碼"),
});

const phoneLoginSchema = z.object({
    countryCode: z.string(),
    phone: z.string().min(1, "請輸入手機號碼"),
});

const otpSchema = z.object({
    otp: z.string().length(6, "驗證碼必須由6位數字組成"),
});

// Internal Component: Social Login Buttons
const SocialLogin = ({ loading }: { loading: boolean }) => (
    <>
        <div className="relative">
            <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                    或用以下方式登入
                </span>
            </div>
        </div>
        <Button
            type="button"
            variant="outline"
            className="w-full"
            disabled={loading}
            onClick={() => window.location.href = `${import.meta.env.VITE_API_ROOT}/api/auth/google/login`}
        >
            Google
        </Button>
    </>
);

// Internal Component: Login Form
interface LoginFormProps {
    onLoginSuccess: () => void;
    onOtpSent: (phone: string) => void;
    onSwitchToRegister: () => void;
}

const LoginForm = ({ onLoginSuccess, onOtpSent, onSwitchToRegister }: LoginFormProps) => {
    const [loginMethod, setLoginMethod] = useState<'phone' | 'password'>('phone');
    const [loading, setLoading] = useState(false);
    const dispatch = useDispatch();
    const cartId = useSelector((state: RootState) => state.cart.cartId);

    const passwordLoginForm = useForm<z.infer<typeof passwordLoginSchema>>({
        resolver: zodResolver(passwordLoginSchema),
        defaultValues: { username: '', password: '' }
    });

    const phoneLoginForm = useForm<z.infer<typeof phoneLoginSchema>>({
        resolver: zodResolver(phoneLoginSchema),
        defaultValues: { countryCode: '+853', phone: '' }
    });

    const onPasswordLogin = async (data: z.infer<typeof passwordLoginSchema>) => {
        setLoading(true);
        try {
            const apiData = {
                username: data.username,
                password: data.password,
                cart_id: cartId ? String(cartId) : undefined
            };
            const response = await callAPI('clientLogin', {}, apiData);
            dispatch(setCredentials({ token: response.token }));
            if (response.cart_id) {
                dispatch(setCartId(response.cart_id));
            }
            toast.success('登入成功');
            onLoginSuccess();
        } catch (error: any) {
            console.error(error);
            const msg = error.response?.data?.error?.message;
            if (error.response?.status === 404) {
                toast.error('找不到此使用者');
            } else if (error.response?.status === 401) {
                toast.error('密碼錯誤');
            } else if (error.response?.status === 403) {
                toast.error('帳號已被停用');
            } else {
                toast.error(msg || '發生錯誤');
            }
        } finally {
            setLoading(false);
        }
    };

    const onPhoneLogin = async (data: z.infer<typeof phoneLoginSchema>) => {
        setLoading(true);
        const fullPhone = `${data.countryCode}${data.phone.replace(/^0+/, '')}`;
        try {
            await callAPI('clientLogin', {}, { phone: fullPhone });
            toast.success('驗證碼已發送');
            onOtpSent(fullPhone);
        } catch (error: any) {
            console.error(error);
            const msg = error.response?.data?.error?.message;
            if (error.response?.status === 404) {
                toast.error('找不到此電話號碼，請先註冊');
                onSwitchToRegister();
            } else if (error.response?.status === 403) {
                toast.error('帳號已被停用');
            } else {
                toast.error(msg || '發生錯誤');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="space-y-4">
            <Tabs defaultValue="phone" value={loginMethod} onValueChange={(val) => setLoginMethod(val as 'phone' | 'password')} className="w-full">
                <TabsList className="grid w-full grid-cols-2 mb-4">
                    <TabsTrigger value="phone">手機驗證</TabsTrigger>
                    <TabsTrigger value="password">密碼登入</TabsTrigger>
                </TabsList>

                <TabsContent value="phone">
                    <Form {...phoneLoginForm}>
                        <form onSubmit={phoneLoginForm.handleSubmit(onPhoneLogin)} className="space-y-4">
                            <FormField
                                control={phoneLoginForm.control}
                                name="phone"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>手機號碼</FormLabel>
                                        <div className="flex space-x-2">
                                            <FormField
                                                control={phoneLoginForm.control}
                                                name="countryCode"
                                                render={({ field: ccField }) => (
                                                    <FormItem className="space-y-0">
                                                        <Select onValueChange={ccField.onChange} defaultValue={ccField.value} disabled={loading}>
                                                            <FormControl>
                                                                <SelectTrigger className="w-[140px]">
                                                                    <SelectValue placeholder="Code" />
                                                                </SelectTrigger>
                                                            </FormControl>
                                                            <SelectContent>
                                                                {COUNTRY_CODES.map((c) => (
                                                                    <SelectItem key={c.code} value={c.code}>
                                                                        {c.label}
                                                                    </SelectItem>
                                                                ))}
                                                            </SelectContent>
                                                        </Select>
                                                    </FormItem>
                                                )}
                                            />
                                            <FormControl>
                                                <Input placeholder="電話號碼" type="tel" {...field} className="flex-1" />
                                            </FormControl>
                                        </div>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <Button type="submit" className="w-full" disabled={loading}>
                                {loading ? '發送中...' : '取得驗證碼'}
                            </Button>
                        </form>
                    </Form>
                </TabsContent>

                <TabsContent value="password">
                    <Form {...passwordLoginForm}>
                        <form onSubmit={passwordLoginForm.handleSubmit(onPasswordLogin)} className="space-y-4 w-full">
                            <FormField
                                control={passwordLoginForm.control}
                                name="username"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>使用者名稱</FormLabel>
                                        <FormControl>
                                            <Input placeholder="使用者名稱" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={passwordLoginForm.control}
                                name="password"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>密碼</FormLabel>
                                        <FormControl>
                                            <Input type="password" placeholder="請輸入密碼" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <Button type="submit" className="w-full" disabled={loading}>
                                {loading ? '登入中...' : '登入'}
                            </Button>
                        </form>
                    </Form>
                </TabsContent>
            </Tabs>

            <div className="mt-4 text-center text-sm">
                <button type="button" onClick={onSwitchToRegister} className="underline">
                    還沒有帳號？立即註冊
                </button>
            </div>

            <SocialLogin loading={loading} />
        </div>
    );
};


// Internal Component: Register Form
interface RegisterFormProps {
    onOtpSent: (phone: string) => void;
    onSwitchToLogin: () => void;
}

const RegisterForm = ({ onOtpSent, onSwitchToLogin }: RegisterFormProps) => {
    const [loading, setLoading] = useState(false);
    
    const registerForm = useForm<z.infer<typeof registerSchema>>({
        resolver: zodResolver(registerSchema),
        defaultValues: { countryCode: '+853', phone: '', username: '', password: '' }
    });

    const onRegister = async (data: z.infer<typeof registerSchema>) => {
        setLoading(true);
        const fullPhone = `${data.countryCode}${data.phone.replace(/^0+/, '')}`;
        try {
            await callAPI('clientRegister', {}, { phone: fullPhone, username: data.username, password: data.password });
            toast.success('驗證碼已發送，請查收');
            onOtpSent(fullPhone);
        } catch (error: any) {
            console.error(error);
            const msg = error.response?.data?.error?.message;
            if (error.response?.status === 409) {
                toast.error('此電話號碼或使用者名稱已註冊');
            } else {
                toast.error(msg || '發生錯誤');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="space-y-4">
            <Form {...registerForm}>
                <form onSubmit={registerForm.handleSubmit(onRegister)} className="space-y-4">
                    <FormField
                        control={registerForm.control}
                        name="phone"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>手機號碼</FormLabel>
                                <div className="flex space-x-2">
                                    <FormField
                                        control={registerForm.control}
                                        name="countryCode"
                                        render={({ field: ccField }) => (
                                            <FormItem className="space-y-0">
                                                <Select onValueChange={ccField.onChange} defaultValue={ccField.value} disabled={loading}>
                                                    <FormControl>
                                                        <SelectTrigger className="w-[140px]">
                                                            <SelectValue placeholder="Code" />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        {COUNTRY_CODES.map((c) => (
                                                            <SelectItem key={c.code} value={c.code}>
                                                                {c.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            </FormItem>
                                        )}
                                    />
                                    <FormControl>
                                        <Input placeholder="電話號碼" type="tel" {...field} className="flex-1" />
                                    </FormControl>
                                </div>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={registerForm.control}
                        name="username"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>使用者名稱</FormLabel>
                                <FormControl>
                                    <Input placeholder="使用者名稱" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={registerForm.control}
                        name="password"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>密碼</FormLabel>
                                <FormControl>
                                    <Input type="password" placeholder="請輸入密碼" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <Button type="submit" className="w-full" disabled={loading}>
                        {loading ? '註冊中...' : '註冊帳號'}
                    </Button>
                </form>
            </Form>

            <div className="mt-4 text-center text-sm">
                <button type="button" onClick={onSwitchToLogin} className="underline">
                    已有帳號？登入
                </button>
            </div>

            <SocialLogin loading={loading} />
        </div>
    );
};

// Internal Component: OTP Form
interface OtpFormProps {
    phone: string;
    onSuccess: () => void;
    onBack: () => void;
}

const OtpForm = ({ phone, onSuccess, onBack }: OtpFormProps) => {
    const [loading, setLoading] = useState(false);
    const dispatch = useDispatch();
    const cartId = useSelector((state: RootState) => state.cart.cartId);

    const otpForm = useForm<z.infer<typeof otpSchema>>({
        resolver: zodResolver(otpSchema),
        defaultValues: { otp: '' }
    });

    const onVerifyOtp = async (data: z.infer<typeof otpSchema>) => {
        setLoading(true);
        try {
            const apiData = {
                phone: phone,
                otp: data.otp,
                cart_id: cartId ? String(cartId) : undefined
            };
            const response = await callAPI('verifyOTP', {}, apiData);
            dispatch(setCredentials({ token: response.token }));
            if (response.cart_id) {
                dispatch(setCartId(response.cart_id));
            }
            toast.success('登入成功');
            onSuccess();
        } catch (error: any) {
            console.error(error);
            toast.error('驗證碼錯誤或過期');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Form {...otpForm}>
            <form onSubmit={otpForm.handleSubmit(onVerifyOtp)} className="space-y-4">
                <div className="space-y-2 flex flex-col items-center">
                    <FormField
                        control={otpForm.control}
                        name="otp"
                        render={({ field }) => (
                            <FormItem className="w-full flex flex-col items-center">
                                <FormLabel className="w-full text-left">驗證碼</FormLabel>
                                <FormControl>
                                    <InputOTP
                                        maxLength={6}
                                        value={field.value}
                                        onChange={field.onChange}
                                        disabled={loading}
                                    >
                                        <InputOTPGroup>
                                            <InputOTPSlot index={0} />
                                            <InputOTPSlot index={1} />
                                            <InputOTPSlot index={2} />
                                            <InputOTPSlot index={3} />
                                            <InputOTPSlot index={4} />
                                            <InputOTPSlot index={5} />
                                        </InputOTPGroup>
                                    </InputOTP>
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                </div>
                <Button type="submit" className="w-full" disabled={loading}>
                    {loading ? '驗證中...' : '登入'}
                </Button>
                <Button
                    type="button"
                    variant="ghost"
                    className="w-full"
                    onClick={onBack}
                    disabled={loading}
                >
                    返回
                </Button>
            </form>
        </Form>
    );
};

interface AuthDrawerProps {
    isOpen: boolean;
    onOpenChange: (open: boolean) => void;
    defaultView?: 'login' | 'register';
}

export default function AuthDrawer({ isOpen, onOpenChange, defaultView = 'login' }: AuthDrawerProps) {
    const [view, setView] = useState<'login' | 'register'>(defaultView);
    const [step, setStep] = useState<'phone' | 'otp'>('phone');
    const [tempPhone4Otp, setTempPhone4Otp] = useState('');
    const isMobile = useIsMobile();

    useEffect(() => {
        if (isOpen) {
            setView(defaultView);
            setStep('phone');
        }
    }, [isOpen, defaultView]);

    const handleSwitchToRegister = () => {
        setView('register');
        setStep('phone');
    };

    const handleSwitchToLogin = () => {
        setView('login');
        setStep('phone');
    };

    const handleOtpStep = (phone: string) => {
        setTempPhone4Otp(phone);
        setStep('otp');
    };

    const AuthContent = () => {
        if (step === 'otp') {
            return (
                <OtpForm 
                    phone={tempPhone4Otp} 
                    onSuccess={() => onOpenChange(false)} 
                    onBack={() => setStep('phone')} 
                />
            );
        }

        if (view === 'login') {
            return (
                <LoginForm 
                    onLoginSuccess={() => onOpenChange(false)} 
                    onOtpSent={handleOtpStep} 
                    onSwitchToRegister={handleSwitchToRegister} 
                />
            );
        }

        return (
            <RegisterForm 
                onOtpSent={handleOtpStep} 
                onSwitchToLogin={handleSwitchToLogin} 
            />
        );
    };

    const title = step === 'otp' ? '輸入驗證碼' : (view === 'login' ? '會員登入' : '註冊帳號');
    const description = step === 'otp'
        ? '請輸入簡訊驗證碼'
        : (view === 'login' ? '請輸入您的手機號碼以登入' : '請輸入手機號碼以建立帳號');

    if (isMobile) {
        return (
            <Drawer open={isOpen} direction="bottom" onOpenChange={onOpenChange}>
                <DrawerContent className="h-[85vh] bottom-0 mt-0 px-2 pb-1 bg-white rounded-t-xl w-full">
                    <div className="relative flex flex-col w-full h-full">
                        <DrawerHeader className='mt-4 text-left'>
                            <DrawerTitle>{title}</DrawerTitle>
                            <DrawerDescription>{description}</DrawerDescription>
                        </DrawerHeader>
                        <div className="p-4 overflow-y-auto">
                            <AuthContent />
                        </div>
                        <DrawerFooter>
                            <DrawerClose asChild>
                                <Button variant="outline">取消</Button>
                            </DrawerClose>
                        </DrawerFooter>
                    </div>
                </DrawerContent>
            </Drawer>
        );
    }

    return (
        <Dialog open={isOpen} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[400px] bg-white">
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    <DialogDescription>{description}</DialogDescription>
                </DialogHeader>
                <div className="pb-0">
                    <AuthContent />
                </div>
            </DialogContent>
        </Dialog>
    );
}

