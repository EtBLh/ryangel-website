import { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { callAPI } from '@/lib/api';
import { setAdminCredentials } from '@/store/adminAuthSlice';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { toast } from 'sonner';

export const AdminLogin = () => {
    const [identifier, setIdentifier] = useState('');
    const [password, setPassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);

        try {
            const response = await callAPI('adminLogin', undefined, {
                username: identifier,
                password,
            });

            // Assuming response structure { token: "...", admin: {...} }
            // Verify backend response structure actually.
            // backend/internal/http/handlers/auth.go adminLogin usually returns token.
            
            // Let's assume standard response based on `auth.go`
            // If response is { token: "..." } then fetch me? Or standard auth response includes user.
            
            dispatch(setAdminCredentials({ 
                token: response.token, 
                admin: response.admin 
            }));
            
            toast.success('Admin login successful');
            navigate('/admin/dashboard');
        } catch (error: any) {
            console.error(error);
            toast.error(error.response?.data?.error?.message || 'Login failed');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <Card className="w-[350px]">
                <CardHeader>
                    <CardTitle className="text-center">Admin Login</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleLogin} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="identifier">Username / Email</Label>
                            <Input
                                id="identifier"
                                type="text"
                                placeholder="admin"
                                value={identifier}
                                onChange={(e) => setIdentifier(e.target.value)}
                                required
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                        </div>
                        <Button className="w-full" type="submit" disabled={isLoading}>
                            {isLoading ? 'Logging in...' : 'Login'}
                        </Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};
