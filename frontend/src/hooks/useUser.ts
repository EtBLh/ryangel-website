import { useQuery } from '@tanstack/react-query';
import { callAPI } from '@/lib/api';
import { useSelector } from 'react-redux';
import type { RootState } from '@/store';
import type { Client } from '@/lib/types';

export const useUser = () => {
    const token = useSelector((state: RootState) => state.auth.token);

    return useQuery<{ client: Client }>({
        queryKey: ['user', token],
        queryFn: () => callAPI('clientMe'),
        enabled: !!token,
        staleTime: 1000 * 60 * 5, // 5 minutes
        retry: false,
        select: (data) => data.client
    });
};
