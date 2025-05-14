import { useState } from 'react';
import { getJiroNowAPI } from '../api/jiroNowAPI';
import type { GetStoreStatusResponseResponse } from '../api/jiroNowAPI.schemas';

interface UseStoreStatusResult {
    storeStatus: GetStoreStatusResponseResponse | null;
    loading: boolean;
    error: string | null;
    showStatus: boolean;
    fetchStoreStatus: (storeId: string) => Promise<void>;
}

export function useStoreStatus(): UseStoreStatusResult {
    const [storeStatus, setStoreStatus] = useState<GetStoreStatusResponseResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showStatus, setShowStatus] = useState(false);

    const fetchStoreStatus = async (storeId: string) => {
        setLoading(true);
        setError(null);
        setShowStatus(false);
        try {
            const api = getJiroNowAPI();
            const response = await api.getStoreStatus({ xAccountId: storeId });
            setStoreStatus(response.data);
            setShowStatus(true);
        } catch (err) {
            setError('Failed to fetch store status');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return {
        storeStatus,
        loading,
        error,
        showStatus,
        fetchStoreStatus,
    };
}
