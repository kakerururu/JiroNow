import type { GetStoreStatusResponseResponse } from '../api/jiroNowAPI.schemas';
import { stores } from '../stores';

interface StoreStatusDisplayProps {
    storeStatus: GetStoreStatusResponseResponse | null;
    selectedStoreId: string;
}

function StoreStatusDisplay({ storeStatus, selectedStoreId }: StoreStatusDisplayProps) {
    const selectedStore = stores.find(store => store.id === selectedStoreId);
    const storeName = selectedStore ? selectedStore.name : 'Unknown Store';

    if (!storeStatus) {
        return <p className="mt-4 text-gray-600">No store status data available.</p>;
    }

    return (
        <div className="mt-4 p-4 border rounded-md shadow-sm bg-white">
            <h2 className="text-xl font-semibold mb-2">{storeName}</h2>
            <div>
                {storeStatus.status === 'open'
                    ? <span className="text-green-600 font-bold">営業中です</span>
                    : storeStatus.status === 'closed'
                        ? <span className="text-red-600 font-bold">営業していません</span>
                        : <span>{storeStatus.status}</span>}
            </div>
        </div>
    );
}

export default StoreStatusDisplay;
