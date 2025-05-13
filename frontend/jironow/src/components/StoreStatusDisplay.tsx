import type { GetStoreStatusResponseResponse } from '../api/jiroNowAPI.schemas';
import './StoreStatusDisplay.css';
import { stores } from '../stores';

interface StoreStatusDisplayProps {
    storeStatus: GetStoreStatusResponseResponse | null;
    loading: boolean;
    error: string | null;
    showStatus: boolean;
    selectedStoreId: string;
}

function StoreStatusDisplay({ storeStatus, loading, error, showStatus, selectedStoreId }: StoreStatusDisplayProps) {
    const selectedStore = stores.find(store => store.id === selectedStoreId);
    const storeName = selectedStore ? selectedStore.name : 'Unknown Store';

    if (loading) {
        return <div className="store-status-container loading">Loading...</div>;
    }

    if (error) {
        return <div className="store-status-container error">Error: {error}</div>;
    }

    if (showStatus) {
        if (storeStatus) {
            return (
                <div className="store-status-container">
                    <h2>{storeName}</h2>
                    <div>
                        {storeStatus.status === 'open'
                            ? '営業中です'
                            : storeStatus.status === 'closed'
                                ? '営業していません'
                                : storeStatus.status}
                    </div>
                </div>
            );
        } else {
            return <p className="store-status-container">No store status data available.</p>;
        }
    }

    return null;
}

export default StoreStatusDisplay;
